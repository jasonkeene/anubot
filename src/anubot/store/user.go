package store

const (
	Streamer TwitchUser = iota
	Bot
)

type TwitchUser int
