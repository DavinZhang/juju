// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common

import (
	stdcontext "context"

	"github.com/DavinZhang/juju/api/base"
	"github.com/DavinZhang/juju/api/credentialvalidator"
	"github.com/DavinZhang/juju/environs/context"
)

// CredentialAPI exposes functionality of the credential validator API facade to a worker.
type CredentialAPI interface {
	InvalidateModelCredential(reason string) error
}

// NewCredentialInvalidatorFacade creates an API facade capable of invalidating credential.
func NewCredentialInvalidatorFacade(apiCaller base.APICaller) (CredentialAPI, error) {
	return credentialvalidator.NewFacade(apiCaller), nil
}

// NewCloudCallContextFunc creates a function returning a cloud call context to be used by workers.
func NewCloudCallContextFunc(c CredentialAPI) CloudCallContextFunc {
	return func(ctx stdcontext.Context) context.ProviderCallContext {
		return &context.CloudCallContext{
			Context:                  ctx,
			InvalidateCredentialFunc: c.InvalidateModelCredential,
		}
	}
}

// CloudCallContextFunc is a function returning a ProviderCallContext.
type CloudCallContextFunc func(ctx stdcontext.Context) context.ProviderCallContext
