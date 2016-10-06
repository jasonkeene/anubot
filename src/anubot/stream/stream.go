package stream

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fluffle/goirc/client"
)

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
	Type    Type
	Twitch  *TXTwitch
	Discord *TXDiscord
}

// TXTwitch contains information to send to Twitch.
type TXTwitch struct {
	Username string
	To       string
	Message  string
}

// TXDiscord contains information to send to Discord.
type TXDiscord struct {
	To      string
	Message string
}

// RXMessage is the data read from a stream source.
type RXMessage struct {
	Type    Type       `json:"type"`
	Twitch  *RXTwitch  `json:"twitch"`
	Discord *RXDiscord `json:"discord"`
}

// RXTwitch contains information received from Twitch.
type RXTwitch struct {
	Line *client.Line `json:"line"`
}

// RXDiscord contains information received from Discord.
type RXDiscord struct {
	MessageCreate *discordgo.MessageCreate `json:"message_create"`
}

// Dispatcher dispoatches messages from a stream source.
type Dispatcher interface {
	Dispatch(topic string, message RXMessage)
}
