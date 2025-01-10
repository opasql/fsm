package fsm

// Option is a type for FSM options
type Option func(*FSM)

// WithUserStateStorage sets userStateStorage FSM
func WithUserStateStorage(storage UserStateStorage) Option {
	return func(fsm *FSM) {
		fsm.userStates = storage
	}
}

// WithDataStorage sets a data storage for FSM
func WithDataStorage(storage DataStorage) Option {
	return func(fsm *FSM) {
		fsm.storage = storage
	}
}
