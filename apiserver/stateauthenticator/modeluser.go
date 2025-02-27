// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package stateauthenticator

import (
	"time"

	"github.com/juju/errors"
	"github.com/juju/names/v4"

	"github.com/DavinZhang/juju/apiserver/common"
	"github.com/DavinZhang/juju/core/permission"
	"github.com/DavinZhang/juju/state"
	stateerrors "github.com/DavinZhang/juju/state/errors"
)

// modelUserEntityFinder implements state.EntityFinder by returning
// an Entity value for model users, ensuring that the user exists in
// the state's current model, while also supporting external users.
type modelUserEntityFinder struct {
	st *state.State
}

// FindEntity implements state.EntityFinder.
func (f modelUserEntityFinder) FindEntity(tag names.Tag) (state.Entity, error) {
	utag, ok := tag.(names.UserTag)
	if !ok {
		return f.st.FindEntity(tag)
	}

	model, err := f.st.Model()
	if err != nil {
		return nil, errors.Trace(err)
	}
	modelUser, err := f.st.UserAccess(utag, model.ModelTag())
	if err != nil && !errors.IsNotFound(err) {
		return nil, errors.Trace(err)
	}

	// No model user found, so see if the user has been granted
	// access to the controller.
	if permission.IsEmptyUserAccess(modelUser) {
		controllerUser, err := state.ControllerAccess(f.st, utag)
		if err != nil && !errors.IsNotFound(err) {
			return nil, errors.Trace(err)
		}
		// TODO(perrito666) remove the following section about everyone group
		// when groups are implemented, this accounts only for the lack of a local
		// ControllerUser when logging in from an external user that has not been granted
		// permissions on the controller but there are permissions for the special
		// everyone group.
		if permission.IsEmptyUserAccess(controllerUser) && !utag.IsLocal() {
			everyoneTag := names.NewUserTag(common.EveryoneTagName)
			controllerUser, err = f.st.UserAccess(everyoneTag, f.st.ControllerTag())
			if err != nil && !errors.IsNotFound(err) {
				return nil, errors.Annotatef(err, "obtaining ControllerUser for everyone group")
			}
		}
		if permission.IsEmptyUserAccess(controllerUser) {
			return nil, errors.NotFoundf("model or controller user")
		}
	}

	u := &modelUserEntity{
		st:        f.st,
		modelUser: modelUser,
		tag:       utag,
	}
	if utag.IsLocal() {
		user, err := f.st.User(utag)
		if err != nil {
			return nil, errors.Trace(err)
		}
		u.user = user
	}
	return u, nil
}

// modelUserEntity encapsulates an model user
// and, if the user is local, the local state user
// as well. This enables us to implement FindEntity
// in such a way that the authentication mechanisms
// can work without knowing these details.
type modelUserEntity struct {
	st *state.State

	modelUser permission.UserAccess
	user      *state.User
	tag       names.Tag
}

// Refresh implements state.Authenticator.Refresh.
func (u *modelUserEntity) Refresh() error {
	if u.user == nil {
		return nil
	}
	return u.user.Refresh()
}

// SetPassword implements state.Authenticator.SetPassword
// by setting the password on the local user.
func (u *modelUserEntity) SetPassword(pass string) error {
	if u.user == nil {
		return errors.New("cannot set password on external user")
	}
	return u.user.SetPassword(pass)
}

// PasswordValid implements state.Authenticator.PasswordValid.
func (u *modelUserEntity) PasswordValid(pass string) bool {
	if u.user == nil {
		return false
	}
	return u.user.PasswordValid(pass)
}

// Tag implements state.Entity.Tag.
func (u *modelUserEntity) Tag() names.Tag {
	return u.tag
}

// LastLogin implements loginEntity.LastLogin.
func (u *modelUserEntity) LastLogin() (time.Time, error) {
	// The last connection for the model takes precedence over
	// the local user last login time.
	var err error
	var t time.Time

	model, err := u.st.Model()
	if err != nil {
		return t, errors.Trace(err)
	}

	if !permission.IsEmptyUserAccess(u.modelUser) {
		t, err = model.LastModelConnection(u.modelUser.UserTag)
	} else {
		err = stateerrors.NewNeverConnectedError("controller user")
	}
	if stateerrors.IsNeverConnectedError(err) || permission.IsEmptyUserAccess(u.modelUser) {
		if u.user != nil {
			// There's a global user, so use that login time instead.
			return u.user.LastLogin()
		}
		// Since we're implementing LastLogin, we need
		// to implement LastLogin error semantics too.
		err = stateerrors.NewNeverLoggedInError(err.Error())
	}
	return t, errors.Trace(err)
}

// UpdateLastLogin implements loginEntity.UpdateLastLogin.
func (u *modelUserEntity) UpdateLastLogin() error {
	updateLastLogin := func() error {
		// If user is nil, don't attempt to perform the update and exit early.
		if u.user == nil {
			return nil
		}

		if err := u.user.UpdateLastLogin(); err != nil {
			return errors.Trace(err)
		}
		return nil
	}

	if !permission.IsEmptyUserAccess(u.modelUser) {
		if u.modelUser.Object.Kind() != names.ModelTagKind {
			return errors.NotValidf("%s as model user", u.modelUser.Object.Kind())
		}

		model, err := u.st.Model()
		if err != nil {
			return errors.Trace(err)
		}

		if err := model.UpdateLastModelConnection(u.modelUser.UserTag); err != nil {
			// Attempt to update the users last login data, if the update
			// fails, then just report it as a log message and return the
			// original error message.
			if err := updateLastLogin(); err != nil {
				logger.Warningf("Unable to update last login with %s", err.Error())
			}
			return errors.Trace(err)
		}
	}

	return updateLastLogin()
}
