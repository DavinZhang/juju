// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package raftlease

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	"github.com/juju/collections/set"
	"github.com/juju/errors"
	"gopkg.in/yaml.v2"

	"github.com/DavinZhang/juju/core/globalclock"
	"github.com/DavinZhang/juju/core/lease"
)

const (
	// CommandVersion is the current version of the command format. If
	// this changes then we need to be sure that reading and applying
	// commands for previous versions still works.
	CommandVersion = 1

	// SnapshotVersion is the current version of the snapshot
	// format. Similarly, changes to the snapshot representation need
	// to be backward-compatible.
	SnapshotVersion = 1

	// OperationClaim denotes claiming a new lease.
	OperationClaim = "claim"

	// OperationExtend denotes extending an already-held lease.
	OperationExtend = "extend"

	// OperationRevoke denotes revoking an existing lease.
	OperationRevoke = "revoke"

	// OperationSetTime denotes updating stored global time (which
	// will also remove any expired leases).
	OperationSetTime = "setTime"

	// OperationPin pins a lease, preventing it from expiring
	// until it is unpinned.
	OperationPin = "pin"

	// OperationUnpin unpins a lease, restoring normal
	// lease expiry behaviour.
	OperationUnpin = "unpin"
)

// FSMResponse defines what will be available on the return value from
// FSM apply calls.
type FSMResponse interface {
	// Error is a lease error (rather than anything to do with the
	// raft machinery).
	Error() error

	// Notify tells the target what changes occurred because of the
	// applied command.
	// On encountering an error, notify will continue until all resulting
	// claims and expiries are exhausted, before returning the error messages.
	Notify(NotifyTarget) error
}

// groupKey stores the namespace and model uuid that identifies all
// the leases for a particular model and lease type.
type groupKey struct {
	namespace string
	modelUUID string
}

// groupKeyFor builds a group key for the given lease key.
func groupKeyFor(key lease.Key) groupKey {
	return groupKey{
		namespace: key.Namespace,
		modelUUID: key.ModelUUID,
	}
}

// NewFSM returns a new FSM to store lease information.
func NewFSM() *FSM {
	return &FSM{
		groups: make(map[groupKey]map[lease.Key]*entry),
		pinned: make(map[lease.Key]set.Strings),
	}
}

// FSM stores the state of leases in the system.
type FSM struct {
	mu         sync.RWMutex
	globalTime time.Time
	groups     map[groupKey]map[lease.Key]*entry

	// Pinned leases are denoted by having a non-empty collection of tags
	// representing the applications requiring pinned behaviour,
	// against their key.
	// This allows different Juju concerns to pin leases, but remove only
	// their own pins. It is done to avoid restoring normal expiration
	// to a lease pinned by another concern operating under the
	// assumption that the lease-holder will not change.
	pinned map[lease.Key]set.Strings
}

func (f *FSM) getGroup(key lease.Key) (map[lease.Key]*entry, bool) {
	entries, found := f.groups[groupKeyFor(key)]
	return entries, found
}

func (f *FSM) ensureGroup(key lease.Key) map[lease.Key]*entry {
	result, found := f.getGroup(key)
	if found {
		return result
	}
	result = make(map[lease.Key]*entry)
	f.groups[groupKeyFor(key)] = result
	return result
}

func (f *FSM) claim(key lease.Key, holder string, duration time.Duration) *response {
	entries := f.ensureGroup(key)
	if entry, found := entries[key]; found {
		// If the claim is for a lease held by someone else,
		// indicate it is already held, so they should not retry.
		if entry.holder != holder {
			return alreadyHeldResponse()
		}

		// If a claim (instead of an extension) is being made by the lease
		// holder, this may be due to a HA situation where the local Raft node
		// is not in sync with the leader. Let them retry.
		logger.Tracef(
			"lease for %q in model %q is already held by claimant %q; allow retry",
			key.Namespace, key.ModelUUID, holder,
		)
		return invalidResponse()
	}
	entries[key] = &entry{
		holder:   holder,
		start:    f.globalTime,
		duration: duration,
	}
	return &response{claimed: key, claimer: holder}
}

func (f *FSM) extend(key lease.Key, holder string, duration time.Duration) *response {
	entries, groupFound := f.getGroup(key)
	if !groupFound {
		return invalidResponse()
	}
	entry, found := entries[key]
	if !found {
		return invalidResponse()
	}
	if entry.holder != holder {
		logger.Tracef(
			"unable to extend lease for %q in model %q; requested for %q, but held by %q",
			key.Namespace, key.ModelUUID, entry.holder, holder,
		)
		return invalidResponse()
	}
	expiry := f.globalTime.Add(duration)
	if !expiry.After(entry.start.Add(entry.duration)) {
		// No extension needed - the lease already expires after the
		// new time.
		return &response{}
	}
	// entry is a pointer back into the f.entries map, so this update
	// isn't lost.
	entry.start = f.globalTime
	entry.duration = duration
	return &response{}
}

func (f *FSM) revoke(key lease.Key, holder string) *response {
	entries, groupFound := f.getGroup(key)
	if !groupFound {
		return invalidResponse()
	}
	entry, found := entries[key]
	if !found {
		return invalidResponse()
	}
	if entry.holder != holder {
		logger.Tracef(
			"unable to revoke lease for %q in model %q; requested for %q, but held by %q",
			key.Namespace, key.ModelUUID, entry.holder, holder,
		)
		return invalidResponse()
	}
	delete(entries, key)
	if len(entries) == 0 {
		delete(f.groups, groupKeyFor(key))
	}
	return &response{expired: []lease.Key{key}}
}

func (f *FSM) pin(key lease.Key, entity string) *response {
	if f.pinned[key] == nil {
		f.pinned[key] = set.NewStrings()
	}
	f.pinned[key].Add(entity)
	return &response{}
}

func (f *FSM) unpin(key lease.Key, entity string) *response {
	if f.pinned[key] != nil {
		f.pinned[key].Remove(entity)
	}
	return &response{}
}

func (f *FSM) setTime(oldTime, newTime time.Time) *response {
	if f.globalTime != oldTime {
		return &response{err: globalclock.ErrOutOfSyncUpdate}
	}
	f.globalTime = newTime
	return &response{expired: f.removeExpired(newTime)}
}

// removeExpired deletes leases that have expired and
// returns a collection of the deleted lease keys.
// Pinned leases are not deleted.
func (f *FSM) removeExpired(newTime time.Time) []lease.Key {
	var expired []lease.Key
	for gKey, entries := range f.groups {
		for key, entry := range entries {
			expiry := entry.start.Add(entry.duration)
			if expiry.Before(newTime) && !f.isPinned(key) {
				delete(entries, key)
				expired = append(expired, key)
			}
		}
		if len(entries) == 0 {
			delete(f.groups, gKey)
		}
	}
	return expired
}

// GlobalTime returns the FSM's internal time.
func (f *FSM) GlobalTime() time.Time {
	return f.globalTime
}

// Leases gets information about all the leases in the system,
// optionally filtered by the input lease keys.
func (f *FSM) Leases(getLocalTime func() time.Time, keys ...lease.Key) map[lease.Key]lease.Info {
	if len(keys) > 0 {
		return f.filteredLeases(getLocalTime, keys)
	}
	return f.allLeases(getLocalTime)
}

// filteredLeases is an optimisation for anticipated usage.
// There will usually be a single key for filtering, so iterating over the
// filter list and retrieving from entries will be fastest by far.
func (f *FSM) filteredLeases(getLocalTime func() time.Time, keys []lease.Key) map[lease.Key]lease.Info {
	results := make(map[lease.Key]lease.Info)
	localTime := getLocalTime()

	f.mu.RLock()
	defer f.mu.RUnlock()

	for _, key := range keys {
		entries, found := f.getGroup(key)
		if !found {
			continue
		}
		if entry, ok := entries[key]; ok {
			results[key] = f.infoFromEntry(localTime, key, entry)
		}
	}
	return results
}

func (f *FSM) allLeases(getLocalTime func() time.Time) map[lease.Key]lease.Info {
	results := make(map[lease.Key]lease.Info)
	localTime := getLocalTime()

	f.mu.RLock()
	defer f.mu.RUnlock()

	for _, entries := range f.groups {
		for key, entry := range entries {
			results[key] = f.infoFromEntry(localTime, key, entry)
		}
	}
	return results
}

func (f *FSM) infoFromEntry(localTime time.Time, key lease.Key, entry *entry) lease.Info {
	globalExpiry := entry.start.Add(entry.duration)

	// Pinned leases are always represented as having an expiry in the future.
	// This prevents the lease manager from waking up thinking it has some
	// expiry events to handle.
	remaining := globalExpiry.Sub(f.globalTime)
	if f.isPinned(key) {
		remaining = 30 * time.Second
	}
	localExpiry := localTime.Add(remaining)

	return lease.Info{
		Holder: entry.holder,
		Expiry: localExpiry,
	}
}

// LeaseGroup returns all leases matching the namespace and model -
// when there are many models this is more efficient than getting all
// the leases and filtering by model.
func (f *FSM) LeaseGroup(getLocalTime func() time.Time, namespace, modelUUID string) map[lease.Key]lease.Info {
	f.mu.RLock()
	defer f.mu.RUnlock()

	gKey := groupKey{namespace: namespace, modelUUID: modelUUID}
	entries, found := f.groups[gKey]
	if !found {
		return nil
	}
	localTime := getLocalTime()
	results := make(map[lease.Key]lease.Info, len(entries))
	for key, entry := range entries {
		results[key] = f.infoFromEntry(localTime, key, entry)
	}
	return results
}

// Pinned returns all the currently known lease pins and
// applications requiring the pinned behaviour.
func (f *FSM) Pinned() map[lease.Key][]string {
	pinned := make(map[lease.Key][]string)

	f.mu.RLock()
	defer f.mu.RUnlock()

	for key, entities := range f.pinned {
		if !entities.IsEmpty() {
			pinned[key] = entities.SortedValues()
		}
	}
	return pinned
}

func (f *FSM) isPinned(key lease.Key) bool {
	return !f.pinned[key].IsEmpty()
}

// entry holds the details of a lease.
type entry struct {
	// holder identifies the current holder of the lease.
	holder string

	// start is the global time at which the lease started.
	start time.Time

	// duration is the duration for which the lease is valid,
	// from the start time.
	duration time.Duration
}

var _ FSMResponse = (*response)(nil)

// response stores what happened as a result of applying a command.
type response struct {
	err     error
	claimer string
	claimed lease.Key
	expired []lease.Key
}

// Error is part of FSMResponse.
func (r *response) Error() error {
	return r.err
}

// Notify is part of FSMResponse.
func (r *response) Notify(target NotifyTarget) error {
	// This response is either for a claim (in which case claimer will be set)
	// or a set-time (so it will have zero or more expirations).
	var errs []error
	if r.claimer != "" {
		if err := target.Claimed(r.claimed, r.claimer); err != nil {
			errs = append(errs, errors.Annotatef(err, "claim lease"))
		}
	}
	for _, expiredKey := range r.expired {
		if err := target.Expired(expiredKey); err != nil {
			errs = append(errs, errors.Annotatef(err, "expire lease"))
		}
	}
	if errs == nil {
		return nil
	}
	return multiErr(errs)
}

type multiErr []error

func (e multiErr) Error() string {
	// TODO(stickupkid): It maybe more prudent to use a pooled bytes buffer for
	// this.
	var s []string
	for _, err := range e {
		if err == nil {
			continue
		}
		s = append(s, err.Error())
	}
	return strings.Join(s, "\n")
}

// Apply log is invoked once a log entry is committed.
// It returns a value which will be made available in the
// ApplyFuture returned by Raft.Apply method if that
// method was called on the same Raft node as the FSM.
// Apply is part of raft.FSM.
func (f *FSM) Apply(log *raft.Log) interface{} {
	command, err := unmarshalCommand(log)
	if err != nil {
		return &response{err: errors.Trace(err)}
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	return f.apply(command)
}

// apply extacts out of the command operation invocations, so that the batching
// FSM can make use of the same logic.
// The caller is expected hold the lock before calling apply.
func (f *FSM) apply(command Command) *response {
	switch command.Operation {
	case OperationClaim:
		return f.claim(command.LeaseKey(), command.Holder, command.Duration)
	case OperationExtend:
		return f.extend(command.LeaseKey(), command.Holder, command.Duration)
	case OperationRevoke:
		return f.revoke(command.LeaseKey(), command.Holder)
	case OperationPin:
		return f.pin(command.LeaseKey(), command.PinEntity)
	case OperationUnpin:
		return f.unpin(command.LeaseKey(), command.PinEntity)
	case OperationSetTime:
		return f.setTime(command.OldTime, command.NewTime)
	default:
		return &response{err: errors.NotValidf("operation %q", command.Operation)}
	}
}

func unmarshalCommand(log *raft.Log) (Command, error) {
	var command Command
	err := yaml.Unmarshal(log.Data, &command)
	if err != nil {
		return command, errors.Trace(err)
	}
	if err := command.Validate(); err != nil {
		return command, errors.Trace(err)
	}
	return command, nil
}

// Snapshot is part of raft.FSM.
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	entries := make(map[SnapshotKey]SnapshotEntry)

	f.mu.RLock()
	defer f.mu.RUnlock()

	for _, group := range f.groups {
		for key, entry := range group {
			entries[SnapshotKey{
				Namespace: key.Namespace,
				ModelUUID: key.ModelUUID,
				Lease:     key.Lease,
			}] = SnapshotEntry{
				Holder:   entry.holder,
				Start:    entry.start,
				Duration: entry.duration,
			}
		}
	}

	pinned := make(map[SnapshotKey][]string)
	for key, entities := range f.pinned {
		if entities.IsEmpty() {
			continue
		}
		pinned[SnapshotKey{
			Namespace: key.Namespace,
			ModelUUID: key.ModelUUID,
			Lease:     key.Lease,
		}] = entities.SortedValues()
	}

	return &Snapshot{
		Version:    SnapshotVersion,
		Entries:    entries,
		Pinned:     pinned,
		GlobalTime: f.globalTime,
	}, nil
}

// Restore is part of raft.FSM.
func (f *FSM) Restore(reader io.ReadCloser) error {
	defer func() { _ = reader.Close() }()

	var snapshot Snapshot
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&snapshot); err != nil {
		return errors.Trace(err)
	}
	if snapshot.Version != SnapshotVersion {
		return errors.NotValidf("snapshot version %d", snapshot.Version)
	}
	if snapshot.Entries == nil {
		return errors.NotValidf("nil entries")
	}

	newGroups := make(map[groupKey]map[lease.Key]*entry, len(snapshot.Entries))
	for key, ssEntry := range snapshot.Entries {
		gKey := groupKey{
			namespace: key.Namespace,
			modelUUID: key.ModelUUID,
		}
		newEntries, found := newGroups[gKey]
		if !found {
			newEntries = make(map[lease.Key]*entry)
			newGroups[gKey] = newEntries
		}

		newEntries[lease.Key{
			Namespace: key.Namespace,
			ModelUUID: key.ModelUUID,
			Lease:     key.Lease,
		}] = &entry{
			holder:   ssEntry.Holder,
			start:    ssEntry.Start,
			duration: ssEntry.Duration,
		}
	}

	newPinned := make(map[lease.Key]set.Strings, len(snapshot.Pinned))
	for key, entities := range snapshot.Pinned {
		newPinned[lease.Key{
			Namespace: key.Namespace,
			ModelUUID: key.ModelUUID,
			Lease:     key.Lease,
		}] = set.NewStrings(entities...)
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.globalTime = snapshot.GlobalTime
	f.groups = newGroups
	f.pinned = newPinned

	return nil
}

// BatchFSM creates a FSM that allows for batching operations. Raft takes
// care of applying the batches in chunked sizes, which allows for restoring
// and snapshotting a the library level. Those should be transparent to the
// the FSM.
type BatchFSM struct {
	*FSM
}

// NewBatchFSM creates a BatchFSM from an existing FSM. By lifting the FSM
// into a BatchFSM allows for better performance when applying logs, it does
// this by only stealing a lock only when required.
func NewBatchFSM(fsm *FSM) *BatchFSM {
	return &BatchFSM{
		FSM: fsm,
	}
}

// ApplyBatch is invoked once a batch of log entries has been committed and
// are ready to be applied to the FSM. ApplyBatch will take in an array of
// log entries. These log entries will be in the order they were committed,
// will not have gaps, and could be of a few log types. Clients should check
// the log type prior to attempting to decode the data attached. Presently
// the LogCommand and LogConfiguration types will be sent.
//
// The returned slice must be the same length as the input and each response
// should correlate to the log at the same index of the input. The returned
// values will be made available in the ApplyFuture returned by Raft.Apply
// method if that method was called on the same Raft node as the FSM.
// ApplyBatch is part of raft.BatchingFSM.
func (f *BatchFSM) ApplyBatch(logs []*raft.Log) interface{} {
	// Unmarshal all the logs up front, we can validate them without
	// stealing the lock. Additionally we can ensure that we get the correct
	// type of log.
	commands := make([]Command, len(logs))
	responses := make([]interface{}, len(logs))

	var numErrs int
	for i, log := range logs {
		command, err := unmarshalCommand(log)
		if err != nil {
			responses[i] = &response{err: errors.Trace(err)}
			numErrs++
			continue
		}

		commands[i] = command
	}

	// If the number of errors matches the number of responses, then we can just
	// return, as there is nothing to process.
	if numErrs == len(logs) {
		return responses
	}

	// Ensure we steal the lock so we serialize all log applications. That way
	// we never allow another Apply or ApplyBatch to intervene.
	f.mu.Lock()

	for i, resp := range responses {
		// We already have a response for the batch item. This happens if we
		// are unable to unmarshal the command from the raft.Log.
		if resp != nil {
			continue
		}

		responses[i] = f.apply(commands[i])
	}

	f.mu.Unlock()

	return responses
}

// Snapshot defines the format of the FSM snapshot.
type Snapshot struct {
	Version    int                           `yaml:"version"`
	Entries    map[SnapshotKey]SnapshotEntry `yaml:"entries"`
	Pinned     map[SnapshotKey][]string      `yaml:"pinned"`
	GlobalTime time.Time                     `yaml:"global-time"`
}

// Persist is part of raft.FSMSnapshot.
func (s *Snapshot) Persist(sink raft.SnapshotSink) (err error) {
	defer func() {
		if err != nil {
			_ = sink.Cancel()
		}
	}()

	encoder := yaml.NewEncoder(sink)
	if err := encoder.Encode(s); err != nil {
		return errors.Trace(err)
	}
	if err := encoder.Close(); err != nil {
		return errors.Trace(err)
	}
	return sink.Close()
}

// Release is part of raft.FSMSnapshot.
func (s *Snapshot) Release() {}

// SnapshotKey defines the format of a lease key in a snapshot.
type SnapshotKey struct {
	Namespace string `yaml:"namespace"`
	ModelUUID string `yaml:"model-uuid"`
	Lease     string `yaml:"lease"`
}

// SnapshotEntry defines the format of a lease entry in a snapshot.
type SnapshotEntry struct {
	Holder   string        `yaml:"holder"`
	Start    time.Time     `yaml:"start"`
	Duration time.Duration `yaml:"duration"`
}

// Command captures the details of an operation to be run on the FSM.
type Command struct {
	// Version of the command format in case it changes,
	// and we need to handle multiple formats.
	Version int `yaml:"version"`

	// Operation is one of claim, extend, expire or setTime.
	Operation string `yaml:"operation"`

	// Namespace is the kind of lease.
	Namespace string `yaml:"namespace,omitempty"`

	// ModelUUID identifies the model the lease belongs to.
	ModelUUID string `yaml:"model-uuid,omitempty"`

	// Lease is the name of the lease the command affects.
	Lease string `yaml:"lease,omitempty"`

	// Holder is the name of the party claiming or extending the
	// lease.
	Holder string `yaml:"holder,omitempty"`

	// Duration is how long the lease should last.
	Duration time.Duration `yaml:"duration,omitempty"`

	// OldTime is the previous time for time updates (to avoid
	// applying stale ones).
	OldTime time.Time `yaml:"old-time,omitempty"`

	// NewTime is the time to store as the global time.
	NewTime time.Time `yaml:"new-time,omitempty"`

	// PinEntity is a tag representing an entity concerned
	// with a pin or unpin operation.
	PinEntity string `yaml:"pin-entity,omitempty"`
}

// Validate checks that the command describes a valid state change.
func (c *Command) Validate() error {
	// For now there's only version 1.
	if c.Version != 1 {
		return errors.NotValidf("version %d", c.Version)
	}
	switch c.Operation {
	case OperationClaim, OperationExtend:
		if err := c.validateLeaseKey(); err != nil {
			return err
		}
		if err := c.validateNoTime(); err != nil {
			return err
		}
		if c.Holder == "" {
			return errors.NotValidf("%s with empty holder", c.Operation)
		}
		if c.Duration == 0 {
			return errors.NotValidf("%s with zero duration", c.Operation)
		}
		if c.PinEntity != "" {
			return errors.NotValidf("%s with pin entity", c.Operation)
		}
	case OperationRevoke:
		if err := c.validateLeaseKey(); err != nil {
			return err
		}
		if err := c.validateNoTime(); err != nil {
			return err
		}
		if c.Duration != 0 {
			return errors.NotValidf("%s with duration", c.Operation)
		}
	case OperationPin, OperationUnpin:
		if err := c.validateLeaseKey(); err != nil {
			return err
		}
		if err := c.validateNoTime(); err != nil {
			return err
		}
		if c.Duration != 0 {
			return errors.NotValidf("%s with duration", c.Operation)
		}
		if c.PinEntity == "" {
			return errors.NotValidf("%s with empty pin entity", c.Operation)
		}
	case OperationSetTime:
		// An old time of 0 is valid when starting up.
		var zeroTime time.Time
		if c.NewTime == zeroTime {
			return errors.NotValidf("setTime with zero new time")
		}
		if c.Holder != "" {
			return errors.NotValidf("setTime with holder")
		}
		if c.Duration != 0 {
			return errors.NotValidf("setTime with duration")
		}
		if c.Namespace != "" {
			return errors.NotValidf("setTime with namespace")
		}
		if c.ModelUUID != "" {
			return errors.NotValidf("setTime with model UUID")
		}
		if c.Lease != "" {
			return errors.NotValidf("setTime with lease")
		}
		if c.PinEntity != "" {
			return errors.NotValidf("setTime with pin entity")
		}
	default:
		return errors.NotValidf("operation %q", c.Operation)
	}
	return nil
}

func (c *Command) validateLeaseKey() error {
	if c.Namespace == "" {
		return errors.NotValidf("%s with empty namespace", c.Operation)
	}
	if c.ModelUUID == "" {
		return errors.NotValidf("%s with empty model UUID", c.Operation)
	}
	if c.Lease == "" {
		return errors.NotValidf("%s with empty lease", c.Operation)
	}
	return nil
}

func (c *Command) validateNoTime() error {
	var zeroTime time.Time
	if c.OldTime != zeroTime {
		return errors.NotValidf("%s with old time", c.Operation)
	}
	if c.NewTime != zeroTime {
		return errors.NotValidf("%s with new time", c.Operation)
	}
	return nil
}

// LeaseKey makes a lease key from the fields in the command.
func (c *Command) LeaseKey() lease.Key {
	return lease.Key{
		Namespace: c.Namespace,
		ModelUUID: c.ModelUUID,
		Lease:     c.Lease,
	}
}

// Marshal converts this command to a byte slice.
func (c *Command) Marshal() ([]byte, error) {
	return yaml.Marshal(c)
}

// String implements fmt.Stringer for the Command type.
func (c *Command) String() string {
	switch c.Operation {
	case OperationSetTime:
		return fmt.Sprintf(
			"Command(ver: %d, op: %s, old time: %v, new time: %v)",
			c.Version, c.Operation, c.OldTime, c.NewTime,
		)
	case OperationPin, OperationUnpin:
		return fmt.Sprintf(
			"Command(ver: %d, op: %s, ns: %s, model: %s, lease: %s, holder: %s, pin entity: %s)",
			c.Version, c.Operation, c.Namespace, c.ModelUUID, c.Lease, c.Holder, c.PinEntity,
		)
	default:
		return fmt.Sprintf(
			"Command(ver: %d, op: %s, ns: %s, model: %.6s, lease: %s, holder: %s)",
			c.Version, c.Operation, c.Namespace, c.ModelUUID, c.Lease, c.Holder,
		)
	}
}

func invalidResponse() *response {
	return &response{err: lease.ErrInvalid}
}

func alreadyHeldResponse() *response {
	return &response{err: lease.ErrHeld}
}
