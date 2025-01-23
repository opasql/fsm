package fsm

import (
	"context"
	"fmt"
)

// StateID is a type for state identifier
type StateID string

// Callback is a function that will be called on state transition
type Callback func(ctx context.Context, args ...any) error

// FSM is a finite state machine
type FSM[K comparable, V any] struct {
	initialStateID StateID
	callbacks      map[StateID]Callback
	userStates     UserStateStorage
	storage        DataStorage[K, V]
}

// UserStateStorage is an interface for user state storage
type UserStateStorage interface {
	Set(userID int64, stateID StateID) error
	Exists(userID int64) (bool, error)
	Get(userID int64) (StateID, error)
}

// DataStorage is an interface for data storage
type DataStorage[K comparable, V any] interface {
	Set(userID int64, key K, value V) error
	Get(userID int64, key K) (any, error)
	Delete(userID int64, key K) error
}

// New creates a new FSM
func New[K comparable, V any](initialStateName StateID, callbacks map[StateID]Callback, opts ...Option[K, V]) *FSM[K, V] {
	s := &FSM[K, V]{
		initialStateID: initialStateName,
		callbacks:      make(map[StateID]Callback),
		userStates:     initialUserStateStorage(),
		storage:        initialDataStorage[K, V](),
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
func (f *FSM[K, V]) AddCallback(stateID StateID, callback Callback) {
	f.callbacks[stateID] = callback
}

// AddCallbacks adds callbacks for states
func (f *FSM[K, V]) AddCallbacks(cb map[StateID]Callback) {
	for stateID, callback := range cb {
		f.callbacks[stateID] = callback
	}
}

// Transition transitions the user to a new state
func (f *FSM[K, V]) Transition(ctx context.Context, userID int64, stateID StateID, args ...any) error {
	err := f.userStates.Set(userID, stateID)
	if err != nil {
		return fmt.Errorf("failed to set user state: %w", err)
	}

	cb, okCb := f.callbacks[stateID]
	if okCb {
		err = cb(ctx, args...)
		if err != nil {
			return fmt.Errorf("failed to execute callback: %w", err)
		}
	}

	return nil
}

// Current returns the current state of the user
func (f *FSM[K, V]) Current(userID int64) (StateID, error) {
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
func (f *FSM[K, V]) Reset(userID int64) error {
	return f.userStates.Set(userID, f.initialStateID)
}

// Set sets a value to data storage by userID and comparable
func (f *FSM[K, V]) Set(userID int64, key K, value V) error {
	err := f.storage.Set(userID, key, value)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	return nil
}

// Get gets a value from data storage by userID and comparable
func (f *FSM[K, V]) Get(userID int64, key K) (any, error) {
	v, err := f.storage.Get(userID, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	return v, nil
}

// Delete deletes a value from data storage by userID and comparable
func (f *FSM[K, V]) Delete(userID int64, key K) error {
	err := f.storage.Delete(userID, key)
	if err != nil {
		return fmt.Errorf("failed to delete user data: %w", err)
	}

	return nil
}
