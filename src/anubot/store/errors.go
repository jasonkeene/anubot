package store

import "errors"

var (
	UsernameTaken = errors.New("username was already taken")
	BadNonce      = errors.New("nonce does not exists")
)
