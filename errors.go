package fsm

import "errors"

var (
	errNoUserData  = errors.New("no user data")
	errNoUserState = errors.New("no user state")
)
