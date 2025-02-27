// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package applicationoffers

import (
	"github.com/juju/charm/v9"
	"github.com/juju/errors"
	"github.com/juju/names/v4"

	commoncrossmodel "github.com/DavinZhang/juju/apiserver/common/crossmodel"
	"github.com/DavinZhang/juju/core/crossmodel"
	"github.com/DavinZhang/juju/core/network"
	"github.com/DavinZhang/juju/core/permission"
	"github.com/DavinZhang/juju/environs/context"
	"github.com/DavinZhang/juju/state"
)

// StatePool provides the subset of a state pool.
type StatePool interface {
	// Get returns a State for a given model from the pool.
	Get(modelUUID string) (Backend, func(), error)

	// Get returns a Model from the pool.
	GetModel(modelUUID string) (Model, func(), error)
}

var GetStatePool = func(sp *state.StatePool) StatePool {
	return &statePoolShim{sp}

}

type statePoolShim struct {
	*state.StatePool
}

func (pool statePoolShim) Get(modelUUID string) (Backend, func(), error) {
	st, err := pool.StatePool.Get(modelUUID)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}
	return &stateShim{
		st:      st.State,
		Backend: commoncrossmodel.GetBackend(st.State),
	}, func() { st.Release() }, nil
}

func (pool statePoolShim) GetModel(modelUUID string) (Model, func(), error) {
	m, ph, err := pool.StatePool.GetModel(modelUUID)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}
	return &modelShim{m}, func() { ph.Release() }, nil
}

// Backend provides selected methods off the state.State struct.
type Backend interface {
	commoncrossmodel.Backend
	Charm(*charm.URL) (commoncrossmodel.Charm, error)
	ApplicationOffer(name string) (*crossmodel.ApplicationOffer, error)
	Model() (Model, error)
	OfferConnections(string) ([]OfferConnection, error)
	SpaceByName(string) (Space, error)
	User(names.UserTag) (User, error)

	CreateOfferAccess(offer names.ApplicationOfferTag, user names.UserTag, access permission.Access) error
	UpdateOfferAccess(offer names.ApplicationOfferTag, user names.UserTag, access permission.Access) error
	RemoveOfferAccess(offer names.ApplicationOfferTag, user names.UserTag) error
	GetOfferUsers(offerUUID string) (map[string]permission.Access, error)

	// GetModelCallContext gets everything that is needed to make cloud calls on behalf of the state current model.
	GetModelCallContext() context.ProviderCallContext

	AllSpaceInfos() (network.SpaceInfos, error)
}

var GetStateAccess = func(st *state.State) Backend {
	return &stateShim{
		st:      st,
		Backend: commoncrossmodel.GetBackend(st),
	}
}

type stateShim struct {
	commoncrossmodel.Backend
	st *state.State
}

func (s stateShim) UserPermission(subject names.UserTag, target names.Tag) (permission.Access, error) {
	return s.st.UserPermission(subject, target)
}

func (s stateShim) CreateOfferAccess(offer names.ApplicationOfferTag, user names.UserTag, access permission.Access) error {
	return s.st.CreateOfferAccess(offer, user, access)
}

func (s stateShim) UpdateOfferAccess(offer names.ApplicationOfferTag, user names.UserTag, access permission.Access) error {
	return s.st.UpdateOfferAccess(offer, user, access)
}

func (s stateShim) RemoveOfferAccess(offer names.ApplicationOfferTag, user names.UserTag) error {
	return s.st.RemoveOfferAccess(offer, user)
}

func (s stateShim) GetOfferUsers(offerUUID string) (map[string]permission.Access, error) {
	return s.st.GetOfferUsers(offerUUID)
}

func (s *stateShim) SpaceByName(name string) (Space, error) {
	sp, err := s.st.SpaceByName(name)
	return &spaceShim{sp}, err
}

func (s *stateShim) Model() (Model, error) {
	m, err := s.st.Model()
	return &modelShim{m}, err
}

func (s *stateShim) GetModelCallContext() context.ProviderCallContext {
	return context.CallContext(s.st)
}

func (s *stateShim) AllSpaceInfos() (network.SpaceInfos, error) {
	return s.st.AllSpaceInfos()
}

type stateCharmShim struct {
	*state.Charm
}

func (s stateShim) Charm(curl *charm.URL) (commoncrossmodel.Charm, error) {
	ch, err := s.st.Charm(curl)
	if err != nil {
		return nil, err
	}
	return stateCharmShim{ch}, nil
}

func (s *stateShim) ApplicationOffer(name string) (*crossmodel.ApplicationOffer, error) {
	offers := state.NewApplicationOffers(s.st)
	return offers.ApplicationOffer(name)
}

var GetApplicationOffers = func(backend interface{}) crossmodel.ApplicationOffers {
	switch st := backend.(type) {
	case *state.State:
		return state.NewApplicationOffers(st)
	case *stateShim:
		return state.NewApplicationOffers(st.st)
	}
	return nil
}

type Subnet interface {
	CIDR() string
	VLANTag() int
	ProviderId() network.Id
	ProviderNetworkId() network.Id
	AvailabilityZones() []string
}

type subnetShim struct {
	*state.Subnet
}

type Space interface {
	Name() string
	Subnets() ([]Subnet, error)
	ProviderId() network.Id
}

type spaceShim struct {
	*state.Space
}

func (s *spaceShim) Subnets() ([]Subnet, error) {
	subnets, err := s.Space.Subnets()
	if err != nil {
		return nil, errors.Trace(err)
	}
	result := make([]Subnet, len(subnets))
	for i, subnet := range subnets {
		result[i] = &subnetShim{subnet}
	}
	return result, nil
}

type Model interface {
	UUID() string
	ModelTag() names.ModelTag
	Name() string
	Type() state.ModelType
	Owner() names.UserTag
}

type modelShim struct {
	*state.Model
}

func (s *stateShim) OfferConnections(offerUUID string) ([]OfferConnection, error) {
	conns, err := s.st.OfferConnections(offerUUID)
	if err != nil {
		return nil, err
	}
	result := make([]OfferConnection, len(conns))
	for i, oc := range conns {
		result[i] = offerConnectionShim{oc}
	}
	return result, nil
}

type OfferConnection interface {
	SourceModelUUID() string
	UserName() string
	RelationKey() string
	RelationId() int
}

type offerConnectionShim struct {
	*state.OfferConnection
}

func (s *stateShim) User(tag names.UserTag) (User, error) {
	return s.st.User(tag)
}

type User interface {
	DisplayName() string
}
