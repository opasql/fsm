package fsm

import (
	"encoding/json"
	"fmt"
	"sync"
)

// dataStorage is a type for default data storage
type dataStorage struct {
	mu      sync.Mutex
	Storage map[int64]map[any]any `json:"storage"`
}

// initialDataStorage creates in memory storage for user's data
func initialDataStorage() *dataStorage {
	return &dataStorage{
		Storage: make(map[int64]map[any]any),
	}
}

// Set sets user's data to data storage
func (d *dataStorage) Set(userID int64, key, value any) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	s, ok := d.Storage[userID]
	if !ok {
		s = make(map[any]any)
		d.Storage[userID] = s
	}

	s[key] = value

	return nil
}

// Get gets user's data from data storage
func (d *dataStorage) Get(userID int64, key any) (any, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.Storage[userID]; !ok {
		return nil, fmt.Errorf("%w, userID:%d, key:%v", errNoUserData, userID, key)
	}

	return d.Storage[userID][key], nil
}

// Delete deletes user's data from data storage
func (d *dataStorage) Delete(userID int64, key any) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.Storage[userID]; !ok {
		return nil
	}

	delete(d.Storage[userID], key)

	return nil
}

// MarshalJSON implements json.Marshaler
func (d *dataStorage) MarshalJSON() ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	return json.Marshal(d.Storage)
}

// UnmarshalJSON implements json.Unmarshaler
func (d *dataStorage) UnmarshalJSON(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	return json.Unmarshal(data, &d.Storage)
}
