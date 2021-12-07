// Copyright 2021 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package secrets

import (
	"testing"

	gc "gopkg.in/check.v1"

	apiservererrors "github.com/DavinZhang/juju/apiserver/errors"
	"github.com/DavinZhang/juju/apiserver/facade"
	"github.com/DavinZhang/juju/secrets"
	coretesting "github.com/DavinZhang/juju/testing"
)

func TestPackage(t *testing.T) {
	gc.TestingT(t)
}

func NewTestAPI(
	service secrets.SecretsService,
	authorizer facade.Authorizer,
) (*SecretsAPI, error) {
	if !authorizer.AuthClient() {
		return nil, apiservererrors.ErrPerm
	}

	return &SecretsAPI{
		authorizer:     authorizer,
		controllerUUID: coretesting.ControllerTag.Id(),
		modelUUID:      coretesting.ModelTag.Id(),
		secretsService: service,
	}, nil
}
