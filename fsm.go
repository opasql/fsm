package fsm

import (
	"encoding/json"
	"fmt"
)

// StateID is a type for state identifier
type StateID string

// Callback is a function that will be called on state transition
type Callback func(f *FSM, args ...any)

// FSM is a finite state machine
type FSM struct {
	initialStateID StateID
	callbacks      map[StateID]Callback
	userStates     UserStateStorage
	storage        DataStorage
}

// UserStateStorage is an interface for user state storage
type UserStateStorage interface {
	Set(userID int64, stateID StateID) error
	Exists(userID int64) (bool, error)
	Get(userID int64) (StateID, error)
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

// DataStorage is an interface for data storage
type DataStorage interface {
	Set(userID int64, key, value any) error
	Get(userID int64, key any) (any, error)
	Delete(userID int64, key any) error
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

// New creates a new FSM
func New(initialStateName StateID, callbacks map[StateID]Callback, opts ...Option) *FSM {
	s := &FSM{
		initialStateID: initialStateName,
		callbacks:      make(map[StateID]Callback),
		userStates:     initialUserStateStorage(),
		storage:        initialDataStorage(),
	}

	for stateID, callback := range callbacks {
		s.callbacks[stateID] = callback
	}

	for _, opt := range opts {
		opt(s)
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
func (f *FSM) Transition(userID int64, stateID StateID, args ...any) error {
	err := f.userStates.Set(userID, stateID)
	if err != nil {
		return fmt.Errorf("failed to set user state: %w", err)
	}

	cb, okCb := f.callbacks[stateID]
	if okCb {
		cb(f, args...)
	}

	return nil
}

// Current returns the current state of the user
func (f *FSM) Current(userID int64) (StateID, error) {
	ok, err := f.userStates.Exists(userID)
	if err != nil {
		return "", fmt.Errorf("failed to check user state: %w", err)
	}
	if !ok {
		err = f.userStates.Set(userID, f.initialStateID)
		if err != nil {
			return "", fmt.Errorf("failed to set user state to initial: %w", err)
		}

		return f.initialStateID, nil
	}

	state, err := f.userStates.Get(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user state: %w", err)
	}

	return state, nil
}

// Reset resets the state of the user to the initial state
func (f *FSM) Reset(userID int64) error {
	return f.userStates.Set(userID, f.initialStateID)
}

// MarshalJSON marshals the FSM to JSON
func (f *FSM) MarshalJSON() ([]byte, error) {

	type response struct {
		InitialStateID StateID          `json:"initial_state_id"`
		UserStates     UserStateStorage `json:"user_states"`
		Storage        DataStorage      `json:"storage"`
	}

	return json.Marshal(response{
		InitialStateID: f.initialStateID,
		UserStates:     f.userStates,
		Storage:        f.storage,
	})
}

// UnmarshalJSON unmarshals the FSM from JSON
func (f *FSM) UnmarshalJSON(data []byte) error {

	type response struct {
		InitialStateID StateID          `json:"initial_state_id"`
		UserStates     UserStateStorage `json:"user_states"`
		Storage        DataStorage      `json:"storage"`
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

// Set sets a value to data storage by userID and key
func (f *FSM) Set(userID int64, key, value any) error {
	err := f.storage.Set(userID, key, value)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	return nil
}

// Get gets a value from data storage by userID and key
func (f *FSM) Get(userID int64, key any) (any, error) {
	v, err := f.storage.Get(userID, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	return v, nil
}

// Delete deletes a value from data storage by userID and key
func (f *FSM) Delete(userID int64, key any) error {
	err := f.storage.Delete(userID, key)
	if err != nil {
		return fmt.Errorf("failed to delete user data: %w", err)
	}

	return nil
}
