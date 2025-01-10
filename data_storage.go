package fsm

import (
	"encoding/json"
	"sync"
)

// defaultDataStorage is a type for default data storage
type dataStorage struct {
	mu      sync.Mutex
	storage map[int64]map[any]any
}

// initialDataStorage creates in memory storage for user's data
func initialDataStorage() *dataStorage {
	return &dataStorage{
		storage: make(map[int64]map[any]any),
	}
}

// Set sets user's data to data storage
func (d *dataStorage) Set(userID int64, key, value any) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	s, ok := d.storage[userID]
	if !ok {
		s = make(map[any]any)
		d.storage[userID] = s
	}

	s[key] = value

	return nil
}

// Get gets user's data from data storage
func (d *dataStorage) Get(userID int64, key any) (any, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, ok := d.storage[userID]
	if !ok {
		return nil, nil
	}

	return d.storage[userID][key], nil
}

// Delete deletes user's data from data storage
func (d *dataStorage) Delete(userID int64, key any) error {
	d.mu.Lock()
	delete(d.storage, userID)
	d.mu.Unlock()
	return nil
}

// MarshalJSON implements json.Marshaler
func (d *dataStorage) MarshalJSON() ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	return json.Marshal(d.storage)
}

// UnmarshalJSON implements json.Unmarshaler
func (d *dataStorage) UnmarshalJSON(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	var response map[int64]map[any]any
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	d.storage = response

	return nil
}
