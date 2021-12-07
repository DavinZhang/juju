// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// Package commands provides functionality for registering all the romulus commands.
package commands

import (
	"github.com/juju/cmd/v3"

	"github.com/DavinZhang/juju/cmd/juju/romulus/agree"
	"github.com/DavinZhang/juju/cmd/juju/romulus/budget"
	"github.com/DavinZhang/juju/cmd/juju/romulus/createwallet"
	"github.com/DavinZhang/juju/cmd/juju/romulus/listagreements"
	"github.com/DavinZhang/juju/cmd/juju/romulus/listplans"
	"github.com/DavinZhang/juju/cmd/juju/romulus/listwallets"
	"github.com/DavinZhang/juju/cmd/juju/romulus/setplan"
	"github.com/DavinZhang/juju/cmd/juju/romulus/setwallet"
	"github.com/DavinZhang/juju/cmd/juju/romulus/showwallet"
	"github.com/DavinZhang/juju/cmd/juju/romulus/sla"
)

type commandRegister interface {
	Register(cmd.Command)
}

// RegisterAll registers all romulus commands with the
// provided command registry.
func RegisterAll(r commandRegister) {
	r.Register(agree.NewAgreeCommand())
	r.Register(listagreements.NewListAgreementsCommand())
	r.Register(budget.NewBudgetCommand())
	r.Register(createwallet.NewCreateWalletCommand())
	r.Register(listplans.NewListPlansCommand())
	r.Register(setwallet.NewSetWalletCommand())
	r.Register(setplan.NewSetPlanCommand())
	r.Register(showwallet.NewShowWalletCommand())
	r.Register(sla.NewSLACommand())
	r.Register(listwallets.NewListWalletsCommand())
}
