package fsm

import "errors"

var (
	ErrNoUserData  = errors.New("no user data")
	ErrNoUserState = errors.New("no user state")
)
