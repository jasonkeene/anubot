package store

import "anubot/stream"

// Store is the interface all storage backends implement.
type Store interface {
	// Close cleans up the resources associated with the storage backend.
	Close() (err error)

	// RegisterUser registers a new user returning the user ID.
	RegisterUser(username, password string) (userID string, err error)

	// AuthenticateUser checks to see if the given user credentials are valid.
	// If they are the user ID is returned with a bool to indicate success.
	AuthenticateUser(username, password string) (userID string, success bool)

	// CreateOauthNonce creates and returns a unique oauth nonce.
	CreateOauthNonce(userID string, tu TwitchUser) (nonce string, err error)

	// OauthNonceExists tells you if the provided nonce was recently created
	// and not yet finished.
	OauthNonceExists(nonce string) (exists bool)

	// FinishOauthNonce completes the oauth flow, removing the nonce and
	// storing the oauth data.
	FinishOauthNonce(nonce, username string, userID int, od OauthData) (err error)

	// TwitchStreamerAuthenticated tells you if the user has authenticated with
	// twitch and that we have valid oauth credentials.
	TwitchStreamerAuthenticated(userID string) (authenticated bool)

	// TwitchStreamerCredentials gives you the credentials for the streamer
	// user.
	TwitchStreamerCredentials(userID string) (username, password string, twitchUserID int)

	// TwitchBotAuthenticated tells you if the user has authenticated his bot
	// with twitch and that we have valid oauth credentials.
	TwitchBotAuthenticated(userID string) (authenticated bool)

	// TwitchBotCredentials gives you the credentials for the streamer user.
	TwitchBotCredentials(userID string) (username, password string, twitchUserID int)

	// TwitchAuthenticated tells you if the user has authenticated his bot and
	// his streamer user with twitch and that we have valid oauth credentials.
	TwitchAuthenticated(userID string) (authenticated bool)

	// TwitchClearAuth removes all the auth data for twitch for the user.
	TwitchClearAuth(userID string)

	// StoreMessage stores a message for a given user for later searching and
	// scrollback history.
	StoreMessage(msg stream.RXMessage) (err error)

	// FetchRecentMessages gets the recent messages for the user's channel.
	FetchRecentMessages(userID string) (msgs []stream.RXMessage, err error)

	// QueryMessages allows the user to search for messages that match a
	// search string.
	QueryMessages(userID, search string) (msgs []stream.RXMessage, err error)
}
