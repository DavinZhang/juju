// Copyright 2012-2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package testing

import (
	"time"

	"github.com/juju/collections/set"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	"gopkg.in/tomb.v2"

	"github.com/DavinZhang/juju/core/watcher"
	"github.com/DavinZhang/juju/testing"
)

type Stopper interface {
	Stop() error
}

func AssertStop(c *gc.C, stopper Stopper) {
	c.Assert(stopper.Stop(), gc.IsNil)
}

type KillWaiter interface {
	Kill()
	Wait() error
}

func AssertKillAndWait(c *gc.C, killWaiter KillWaiter) {
	killWaiter.Kill()
	c.Assert(killWaiter.Wait(), jc.ErrorIsNil)
}

// AssertCanStopWhenSending ensures even when there are changes
// pending to be delivered by the watcher it can still stop
// cleanly. This is necessary to check for deadlocks in case the
// watcher's inner loop is blocked trying to send and its tomb is
// already dying.
func AssertCanStopWhenSending(c *gc.C, stopper Stopper) {
	// Leave some time for the event to be delivered and the watcher
	// to block on sending it.
	<-time.After(testing.ShortWait)
	stopped := make(chan bool)
	// Stop() blocks, so we need to call it in a separate goroutine.
	go func() {
		c.Check(stopper.Stop(), gc.IsNil)
		stopped <- true
	}()
	select {
	case <-time.After(testing.LongWait):
		// NOTE: If this test fails here it means we have a deadlock
		// in the client-side watcher implementation.
		c.Fatalf("watcher did not stop as expected")
	case <-stopped:
	}
}

type NotifyWatcher interface {
	Stop() error
	Changes() <-chan struct{}
}

// NotifyWatcherC embeds a gocheck.C and adds methods to help verify
// the behaviour of any watcher that uses a <-chan struct{}.
type NotifyWatcherC struct {
	*gc.C
	State   SyncStarter
	Watcher NotifyWatcher
}

// SyncStarter is an interface that watcher checkers will use to ensure
// that changes to the watched object have been synchronized. This is
// primarily implemented by state.State.
type SyncStarter interface {
	StartSync()
}

// NewNotifyWatcherC returns a NotifyWatcherC that checks for aggressive
// event coalescence.
func NewNotifyWatcherC(c *gc.C, st SyncStarter, w NotifyWatcher) NotifyWatcherC {
	return NotifyWatcherC{
		C:       c,
		State:   st,
		Watcher: w,
	}
}

func (c NotifyWatcherC) AssertNoChange() {
	c.State.StartSync()
	select {
	case _, ok := <-c.Watcher.Changes():
		if ok {
			c.Fatalf("watcher sent unexpected change")
		} else {
			c.Fatalf("watcher closed Changes channel")
		}
	case <-time.After(testing.ShortWait):
	}
}

func (c NotifyWatcherC) AssertOneChange() {
	// Wait a very small amount of time, so that if there is already an event
	// queued to be processed, we see it, before the StartSync flushes new
	// events into the queue.
	shortTimeout := time.After(1 * time.Millisecond)
	longTimeout := time.After(testing.LongWait)
loop:
	for {
		select {
		case _, ok := <-c.Watcher.Changes():
			c.C.Logf("got change")
			c.Assert(ok, jc.IsTrue)
			break loop
		case <-shortTimeout:
			c.C.Logf("StartSync()")
			c.State.StartSync()
			shortTimeout = nil
		case <-longTimeout:
			c.Fatalf("watcher did not send change")
			break loop
		}
	}
	c.AssertNoChange()
}

func (c NotifyWatcherC) AssertClosed() {
	select {
	case _, ok := <-c.Watcher.Changes():
		c.Assert(ok, jc.IsFalse)
	default:
		c.Fatalf("watcher not closed")
	}
}

// StringsWatcherC embeds a gocheck.C and adds methods to help verify
// the behaviour of any watcher that uses a <-chan []string.
type StringsWatcherC struct {
	*gc.C
	State   SyncStarter
	Watcher StringsWatcher
}

// NewStringsWatcherC returns a StringsWatcherC that checks for aggressive
// event coalescence.
func NewStringsWatcherC(c *gc.C, st SyncStarter, w StringsWatcher) StringsWatcherC {
	return StringsWatcherC{
		C:       c,
		State:   st,
		Watcher: w,
	}
}

type StringsWatcher interface {
	Stop() error
	Changes() <-chan []string
}

func (c StringsWatcherC) AssertNoChange() {
	c.State.StartSync()
	select {
	case actual, ok := <-c.Watcher.Changes():
		c.Fatalf("watcher sent unexpected change: (%v, %v)", actual, ok)
	case <-time.After(testing.ShortWait):
	}
}

func (c StringsWatcherC) AssertChanges() {
	c.State.StartSync()
	select {
	case <-c.Watcher.Changes():
	case <-time.After(testing.LongWait):
		c.Fatalf("watcher did not send change")
	}
}

func (c StringsWatcherC) AssertChange(expect ...string) {
	// We should assert for either a single or multiple changes,
	// based on the number of `expect` changes.
	c.assertChange(len(expect) == 1, expect...)
}

func (c StringsWatcherC) AssertChangeInSingleEvent(expect ...string) {
	c.assertChange(true, expect...)
}

// AssertChangeMaybeIncluding verifies that there is a change that may
// contain zero to all of the passed in strings, and no other changes.
func (c StringsWatcherC) AssertChangeMaybeIncluding(expect ...string) {
	maxCount := len(expect)
	actual := c.collectChanges(true, maxCount)

	if maxCount == 0 {
		c.Assert(actual, gc.HasLen, 0)
	} else {
		actualCount := len(actual)
		c.Assert(actualCount <= maxCount, jc.IsTrue, gc.Commentf("expected at most %d, got %d", maxCount, actualCount))
		unexpected := set.NewStrings(actual...).Difference(set.NewStrings(expect...))
		c.Assert(unexpected.Values(), gc.HasLen, 0)
	}
}

// assertChange asserts the given list of changes was reported by
// the watcher, but does not assume there are no following changes.
func (c StringsWatcherC) assertChange(single bool, expect ...string) {
	actual := c.collectChanges(single, len(expect))
	if len(expect) == 0 {
		c.Assert(actual, gc.HasLen, 0)
	} else {
		c.Assert(actual, jc.SameContents, expect)
	}
}

// collectChanges gets up to the max number of changes within the
// testing.LongWait period.
func (c StringsWatcherC) collectChanges(single bool, max int) []string {
	c.State.StartSync()
	timeout := time.After(testing.LongWait)
	var actual []string
	gotOneChange := false
loop:
	for {
		select {
		case changes, ok := <-c.Watcher.Changes():
			c.Assert(ok, jc.IsTrue)
			gotOneChange = true
			actual = append(actual, changes...)
			if single || len(actual) >= max {
				break loop
			}
		case <-timeout:
			if !gotOneChange {
				c.Fatalf("watcher did not send change")
			}
			// If we triggered a timeout, stop looking for more changes
			break loop
		}
	}
	return actual
}

func (c StringsWatcherC) AssertClosed() {
	select {
	case _, ok := <-c.Watcher.Changes():
		c.Assert(ok, jc.IsFalse)
	default:
		c.Fatalf("watcher not closed")
	}
}

// RelationUnitsWatcherC embeds a gocheck.C and adds methods to help
// verify the behaviour of any watcher that uses a <-chan
// params.RelationUnitsChange.
type RelationUnitsWatcherC struct {
	*gc.C
	State   SyncStarter
	Watcher RelationUnitsWatcher
	// settingsVersions keeps track of the settings version of each
	// changed unit since the last received changes to ensure version
	// always increases.
	settingsVersions    map[string]int64
	appSettingsVersions map[string]int64
}

// NewRelationUnitsWatcherC returns a RelationUnitsWatcherC that
// checks for aggressive event coalescence.
func NewRelationUnitsWatcherC(c *gc.C, st SyncStarter, w RelationUnitsWatcher) RelationUnitsWatcherC {
	return RelationUnitsWatcherC{
		C:                   c,
		State:               st,
		Watcher:             w,
		settingsVersions:    make(map[string]int64),
		appSettingsVersions: make(map[string]int64),
	}
}

type RelationUnitsWatcher interface {
	Stop() error
	Changes() watcher.RelationUnitsChannel
}

func (c RelationUnitsWatcherC) AssertNoChange() {
	c.State.StartSync()
	select {
	case actual, ok := <-c.Watcher.Changes():
		c.Fatalf("watcher sent unexpected change: (%v, %v)", actual, ok)
	case <-time.After(testing.ShortWait):
	}
}

// AssertChange asserts the given changes was reported by the watcher,
// but does not assume there are no following changes.
func (c RelationUnitsWatcherC) AssertChange(changed []string, appChanged []string, departed []string) {
	// Get all items in changed in a map for easy lookup.
	changedNames := set.NewStrings(changed...)
	appChangedNames := set.NewStrings(appChanged...)
	c.State.StartSync()
	timeout := time.After(testing.LongWait)
	select {
	case actual, ok := <-c.Watcher.Changes():
		c.Logf("Watcher.Changes() => %# v", actual)
		c.Assert(ok, jc.IsTrue)
		c.Check(actual.Changed, gc.HasLen, len(changed))
		c.Check(actual.AppChanged, gc.HasLen, len(appChanged))
		// Because the versions can change, we only need to make sure
		// the keys match, not the contents (UnitSettings == txnRevno).
		for k, settings := range actual.Changed {
			c.Check(changedNames.Contains(k), jc.IsTrue)
			oldVer, ok := c.settingsVersions[k]
			if !ok {
				// TODO(jam): 2019-10-22 shouldn't we update this *every* time we see it?
				// This is the first time we see this unit, so
				// save the settings version for later.
				c.settingsVersions[k] = settings.Version
			} else {
				// Already seen; make sure the version increased.
				c.Assert(settings.Version, jc.GreaterThan, oldVer,
					gc.Commentf("expected unit settings to increase got %d had %d",
						settings.Version, oldVer))
			}
		}
		for k, version := range actual.AppChanged {
			c.Check(appChangedNames.Contains(k), jc.IsTrue)
			oldVer, ok := c.appSettingsVersions[k]
			if ok {
				// Make sure if we've seen this setting before, it has been updated
				c.Assert(version, jc.GreaterThan, oldVer,
					gc.Commentf("expected app settings to increase got %d had %d",
						version, oldVer))
			}
			c.appSettingsVersions[k] = version
		}
		c.Check(actual.Departed, jc.SameContents, departed)
	case <-timeout:
		c.Fatalf("watcher did not send change")
	}
}

func (c RelationUnitsWatcherC) AssertClosed() {
	select {
	case _, ok := <-c.Watcher.Changes():
		c.Assert(ok, jc.IsFalse)
	default:
		c.Fatalf("watcher not closed")
	}
}

// SecretsRotationWatcherC embeds a gocheck.C and adds methods to help
// verify the behaviour of any watcher that uses a
// <-chan []SecretRotationChange
type SecretsRotationWatcherC struct {
	*gc.C
	State   SyncStarter
	Watcher SecretsRotationWatcher
}

// NewSecretsRotationWatcherC returns a SecretsRotationWatcherC that
// checks for aggressive event coalescence.
func NewSecretsRotationWatcherC(c *gc.C, st SyncStarter, w SecretsRotationWatcher) SecretsRotationWatcherC {
	return SecretsRotationWatcherC{
		C:       c,
		State:   st,
		Watcher: w,
	}
}

type SecretsRotationWatcher interface {
	Stop() error
	Changes() watcher.SecretRotationChannel
}

func (c SecretsRotationWatcherC) AssertNoChange() {
	c.State.StartSync()
	select {
	case actual, ok := <-c.Watcher.Changes():
		c.Fatalf("watcher sent unexpected change: (%v, %v)", actual, ok)
	case <-time.After(testing.ShortWait):
	}
}

// AssertChange asserts the given changes was reported by the watcher,
// but does not assume there are no following changes.
func (c SecretsRotationWatcherC) AssertChange(expect ...watcher.SecretRotationChange) {
	c.State.StartSync()
	var received []watcher.SecretRotationChange
	timeout := time.After(testing.LongWait)
	for a := testing.LongAttempt.Start(); a.Next(); {
		select {
		case actual, ok := <-c.Watcher.Changes():
			c.Logf("Watcher.Changes() => %# v", actual)
			c.Assert(ok, jc.IsTrue)
			received = append(received, actual...)
			if len(received) >= len(expect) {
				c.Assert(received, jc.DeepEquals, expect)
				return
			}
		case <-timeout:
			c.Fatalf("watcher did not send change")
		}
	}
}

func (c SecretsRotationWatcherC) AssertClosed() {
	select {
	case _, ok := <-c.Watcher.Changes():
		c.Assert(ok, jc.IsFalse)
	default:
		c.Fatalf("watcher not closed")
	}
}

// MockNotifyWatcher implements state.NotifyWatcher.
type MockNotifyWatcher struct {
	tomb tomb.Tomb
	ch   <-chan struct{}
}

func NewMockNotifyWatcher(ch <-chan struct{}) *MockNotifyWatcher {
	w := &MockNotifyWatcher{ch: ch}
	w.tomb.Go(func() error {
		<-w.tomb.Dying()
		return tomb.ErrDying
	})
	return w
}

func (w *MockNotifyWatcher) Changes() <-chan struct{} {
	return w.ch
}

func (w *MockNotifyWatcher) Stop() error {
	w.Kill()
	return w.Wait()
}

func (w *MockNotifyWatcher) Kill() {
	w.tomb.Kill(nil)
}

func (w *MockNotifyWatcher) Err() error {
	return w.tomb.Err()
}

func (w *MockNotifyWatcher) Wait() error {
	return w.tomb.Wait()
}

// MockStringsWatcher implements state.StringsWatcher.
type MockStringsWatcher struct {
	tomb tomb.Tomb
	ch   <-chan []string
}

func NewMockStringsWatcher(ch <-chan []string) *MockStringsWatcher {
	w := &MockStringsWatcher{ch: ch}
	w.tomb.Go(func() error {
		<-w.tomb.Dying()
		return tomb.ErrDying
	})
	return w
}

func (w *MockStringsWatcher) Changes() <-chan []string {
	return w.ch
}

func (w *MockStringsWatcher) Stop() error {
	w.Kill()
	return w.Wait()
}

func (w *MockStringsWatcher) Kill() {
	w.tomb.Kill(nil)
}

func (w *MockStringsWatcher) Err() error {
	return w.tomb.Err()
}

func (w *MockStringsWatcher) Wait() error {
	return w.tomb.Wait()
}
