// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package testing

import (
	"github.com/juju/utils/v2"

	"github.com/DavinZhang/juju/environs"
	"github.com/DavinZhang/juju/provider/common"
)

// impatientAttempt is an extremely short polling time suitable for tests.
// It polls at least once, never delays, and times out very quickly.
//
// TODO(katco): 2016-08-09: lp:1611427
var impatientAttempt = utils.AttemptStrategy{}

// savedAttemptStrategy holds the state needed to restore an AttemptStrategy's
// original setting.
//
// TODO(katco): 2016-08-09: lp:1611427
type savedAttemptStrategy struct {
	address  *utils.AttemptStrategy
	original utils.AttemptStrategy
}

// saveAttemptStrategies captures the information required to restore the
// given AttemptStrategy objects.
//
// TODO(katco): 2016-08-09: lp:1611427
func saveAttemptStrategies(strategies []*utils.AttemptStrategy) []savedAttemptStrategy {
	savedStrategies := make([]savedAttemptStrategy, len(strategies))
	for index, strategy := range strategies {
		savedStrategies[index] = savedAttemptStrategy{
			address:  strategy,
			original: *strategy,
		}
	}
	return savedStrategies
}

// restore returns a saved strategy to its original configuration.
func (saved savedAttemptStrategy) restore() {
	*saved.address = saved.original
}

// restoreAttemptStrategies restores multiple saved AttemptStrategies.
func restoreAttemptStrategies(strategies []savedAttemptStrategy) {
	for _, saved := range strategies {
		saved.restore()
	}
}

// internalPatchAttemptStrategies sets the given AttemptStrategy objects
// to the impatientAttempt configuration, and returns a function that restores
// the original configurations.
//
// TODO(katco): 2016-08-09: lp:1611427
func internalPatchAttemptStrategies(strategies []*utils.AttemptStrategy) func() {
	snapshot := saveAttemptStrategies(strategies)
	for _, strategy := range strategies {
		*strategy = impatientAttempt
	}
	return func() { restoreAttemptStrategies(snapshot) }
}

// TODO: Everything up to this point is generic.  Move it to utils?

// PatchAttemptStrategies patches environs' global polling strategy, plus any
// otther AttemptStrategy objects whose addresses you pass, to very short
// polling and timeout times so that tests can run fast.
// It returns a cleanup function that restores the original settings.  You must
// call this afterwards.
//
// TODO(katco): 2016-08-09: lp:1611427
func PatchAttemptStrategies(strategies ...*utils.AttemptStrategy) func() {
	// The one irregularity here is that LongAttempt goes on the list of
	// strategies that need patching.  To keep testing simple, we treat
	// the given attempts and LongAttempt as a single slice from here on.
	combinedStrategies := append(
		strategies,
		&common.LongAttempt,
		&common.ShortAttempt,
		&environs.AddressesRefreshAttempt,
	)
	return internalPatchAttemptStrategies(combinedStrategies)
}
