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
	// Username is the username to send the message as.
	Username string
	// To is where to send the message to (channel or user).
	To string
	// Message is the content of the message.
	Message string
}

// TXDiscordType is the type of message to send to Discord ()
type TXDiscordType int

const (
	// Channel is used when sending messages to Discord channels.
	Channel TXDiscordType = iota
	// Private is used when sending direct messages to Discord.
	Private
)

// TXDiscord contains information to send to Discord.
type TXDiscord struct {
	// Type is the type of message to send (Channel or Private).
	Type TXDiscordType
	// To is where to send the message to (channel ID or user ID).
	To string
	// Message is the content of the message.
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
	// TODO: add other types
	MessageCreate *discordgo.MessageCreate `json:"message_create"`
}

// Dispatcher dispoatches messages from a stream source.
type Dispatcher interface {
	Dispatch(topic string, message RXMessage)
}
