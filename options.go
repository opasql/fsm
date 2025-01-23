package fsm

// Option is a type for FSM options
type Option[K comparable, V any] func(*FSM[K, V])

// WithUserStateStorage sets userStateStorage FSM
func WithUserStateStorage[K comparable, V any](storage UserStateStorage) Option[K, V] {
	return func(fsm *FSM[K, V]) {
		fsm.userStates = storage
	}
}

// WithDataStorage sets a data storage for FSM
func WithDataStorage[K comparable, V any](storage DataStorage[K, V]) Option[K, V] {
	return func(fsm *FSM[K, V]) {
		fsm.storage = storage
	}
}
