package store

const (
	// Streamer is the primary Twitch user that is streaming.
	Streamer TwitchUser = iota
	// Bot is the Twitch user that represents the bot.
	Bot
)

// TwitchUser is a value that represets either a Stream or Bot user on Twitch.
type TwitchUser int
