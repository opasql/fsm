package fsm

import (
	"encoding/json"
	"sync"
)

// userStateStorage is a type for default user's state storage
type userStateStorage struct {
	mu      sync.RWMutex
	storage map[int64]StateID
}

// initialUserStateStorage creates in memory storage for user's state
func initialUserStateStorage() *userStateStorage {
	return &userStateStorage{
		storage: make(map[int64]StateID),
	}
}

// Set sets user's state to state storage
func (u *userStateStorage) Set(userID int64, stateID StateID) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.storage[userID] = stateID
	return nil
}

// Exists checks whether any user's state exist in state storage
func (u *userStateStorage) Exists(userID int64) (bool, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	_, ok := u.storage[userID]
	return ok, nil
}

// Get gets user's state from state storage
func (u *userStateStorage) Get(userID int64) (StateID, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	s, ok := u.storage[userID]
	if !ok {
		return "", nil
	}

	return s, nil
}

// MarshalJSON implements json.Marshaler
func (u *userStateStorage) MarshalJSON() ([]byte, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	return json.Marshal(u.storage)
}

// UnmarshalJSON implements json.Unmarshaler
func (u *userStateStorage) UnmarshalJSON(data []byte) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	return json.Unmarshal(data, &u.storage)
}
