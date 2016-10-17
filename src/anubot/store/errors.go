package store

import "errors"

var (
	// ErrUnknownUserID is returned when providing a user ID that does not exist.
	ErrUnknownUserID = errors.New("user id does not exists")
	// ErrUsernameTaken is returned when attempting to register with a username
	// that is already taken.
	ErrUsernameTaken = errors.New("username was already taken")
	// ErrUnknownUsername is returned when providing a username that does not
	// exist.
	ErrUnknownUsername = errors.New("username does not exists")

	// ErrUnknownNonce is returned when providing a nonce that does not exist.
	ErrUnknownNonce = errors.New("nonce does not exists")

	// ErrInvalidTwitchUserType is returned when providing an invalid twitch
	// user type.
	ErrInvalidTwitchUserType = errors.New("invalid twitch user type")
)
