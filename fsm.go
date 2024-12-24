package fsm

import (
	"encoding/json"
	"sync"
)

// StateID is a type for state identifier
type StateID string

// Callback is a function that will be called on state transition
type Callback func(f *FSM, args ...any)

// FSM is a finite state machine
type FSM struct {
	initialStateID StateID
	callbacks      map[StateID]Callback
	userStatesMu   sync.RWMutex
	userStates     map[int64]StateID
	storageMx      sync.Mutex
	storage        map[int64]map[any]any
}

// New creates a new FSM
func New(initialStateName StateID, callbacks map[StateID]Callback) *FSM {
	s := &FSM{
		initialStateID: initialStateName,
		callbacks:      make(map[StateID]Callback),
		userStates:     make(map[int64]StateID),
		storage:        make(map[int64]map[any]any),
	}

	for stateID, callback := range callbacks {
		s.callbacks[stateID] = callback
	}

	return s
}

// AddCallback adds a callback for a state
func (f *FSM) AddCallback(stateID StateID, callback Callback) {
	f.callbacks[stateID] = callback
}

// AddCallbacks adds callbacks for states
func (f *FSM) AddCallbacks(cb map[StateID]Callback) {
	for stateID, callback := range cb {
		f.callbacks[stateID] = callback
	}
}

// Transition transitions the user to a new state
func (f *FSM) Transition(userID int64, stateID StateID, args ...any) {
	f.userStatesMu.Lock()

	userStateID, okUserState := f.userStates[userID]
	if !okUserState {
		userStateID = f.initialStateID
		f.userStates[userID] = userStateID
	}
	f.userStates[userID] = stateID

	f.userStatesMu.Unlock()

	cb, okCb := f.callbacks[stateID]
	if okCb {
		cb(f, args...)
	}
}

// Current returns the current state of the user
func (f *FSM) Current(userID int64) StateID {
	f.userStatesMu.RLock()
	defer f.userStatesMu.RUnlock()

	userStateID, ok := f.userStates[userID]
	if !ok {
		f.userStates[userID] = f.initialStateID
		return f.initialStateID
	}

	return userStateID
}

// Reset resets the state of the user to the initial state
func (f *FSM) Reset(userID int64) {
	f.userStatesMu.Lock()
	delete(f.userStates, userID)
	f.userStatesMu.Unlock()
}

// MarshalJSON marshals the FSM to JSON
func (f *FSM) MarshalJSON() ([]byte, error) {
	f.userStatesMu.RLock()
	defer f.userStatesMu.RUnlock()

	type response struct {
		InitialStateID StateID               `json:"initial_state_id"`
		UserStates     map[int64]StateID     `json:"user_states"`
		Storage        map[int64]map[any]any `json:"storage"`
	}

	return json.Marshal(response{
		InitialStateID: f.initialStateID,
		UserStates:     f.userStates,
		Storage:        f.storage,
	})
}

// UnmarshalJSON unmarshals the FSM from JSON
func (f *FSM) UnmarshalJSON(data []byte) error {
	f.userStatesMu.Lock()
	defer f.userStatesMu.Unlock()

	type response struct {
		InitialStateID StateID               `json:"initial_state_id"`
		UserStates     map[int64]StateID     `json:"user_states"`
		Storage        map[int64]map[any]any `json:"storage"`
	}

	var r response
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	f.initialStateID = r.InitialStateID
	f.userStates = r.UserStates
	f.storage = r.Storage

	return nil
}

// Set sets a value for a key for a user
func (f *FSM) Set(userID int64, key, value any) {
	f.storageMx.Lock()
	defer f.storageMx.Unlock()
	s, ok := f.storage[userID]
	if !ok {
		s = make(map[any]any)
		f.storage[userID] = s
	}
	s[key] = value
}

// Get gets a value for a key for a user
func (f *FSM) Get(userID int64, key any) (any, bool) {
	f.storageMx.Lock()
	defer f.storageMx.Unlock()
	s, ok := f.storage[userID]
	if !ok {
		return nil, false
	}
	v, ok := s[key]
	return v, ok
}
