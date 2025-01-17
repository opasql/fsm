package fsm

import (
	"encoding/json"
	"fmt"
	"sync"
)

// userStateStorage is a type for default user's state storage
type userStateStorage struct {
	mu      sync.RWMutex
	Storage map[int64]StateID `json:"storage"`
}

// initialUserStateStorage creates in memory storage for user's state
func initialUserStateStorage() *userStateStorage {
	return &userStateStorage{
		Storage: make(map[int64]StateID),
	}
}

// Set sets user's state to state storage
func (u *userStateStorage) Set(userID int64, stateID StateID) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.Storage[userID] = stateID

	return nil
}

// Exists checks whether any user's state exist in state storage
func (u *userStateStorage) Exists(userID int64) (bool, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	_, ok := u.Storage[userID]

	return ok, nil
}

// Get gets user's state from state storage
func (u *userStateStorage) Get(userID int64) (StateID, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	s, ok := u.Storage[userID]
	if !ok {
		return "", fmt.Errorf("%w: userID: %d", errNoUserState, userID)
	}

	return s, nil
}

// MarshalJSON implements json.Marshaler
func (u *userStateStorage) MarshalJSON() ([]byte, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	return json.Marshal(u.Storage)
}

// UnmarshalJSON implements json.Unmarshaler
func (u *userStateStorage) UnmarshalJSON(data []byte) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	return json.Unmarshal(data, &u.Storage)
}
