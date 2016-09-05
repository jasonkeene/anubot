package stream

import "github.com/fluffle/goirc/client"

// Type is the type of the stream (Twitch, Discord).
type Type int

const (
	// Twitch is the the type for Twitch streams.
	Twitch Type = iota
	// Discord is the the type for Discord streams.
	Discord
)

// TXMessage is the data writen out to a stream source.
type TXMessage struct {
	To      string
	Message string
}

// RXMessage is the data read from a stream source.
type RXMessage struct {
	Type    Type
	Twitch  *RXTwitch
	Discord *RXDiscord
}

// RXTwitch contains information received from Twitch.
type RXTwitch struct {
	Line *client.Line
}

// RXDiscord contains information received from Discord.
type RXDiscord struct {
}

// Dispatcher dispoatches messages from a stream source.
type Dispatcher interface {
	Dispatch(message RXMessage)
}
