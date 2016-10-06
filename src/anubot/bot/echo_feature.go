package bot

import (
	"anubot/stream"
	"fmt"
	"strings"
)

// EchoFeature echos messages back to the user.
type EchoFeature struct {
	cmd            string
	twitchUsername string
	sman           *stream.Manager
}

// NewEchoFeature returns a new echo feature.
func NewEchoFeature(cmd, twitchUsername string, sman *stream.Manager) *EchoFeature {
	return &EchoFeature{
		cmd:            cmd,
		twitchUsername: twitchUsername,
		sman:           sman,
	}
}

// HandleMessage echos back the body from the message.
func (e *EchoFeature) HandleMessage(in stream.RXMessage) {
	out := stream.TXMessage{
		Type: in.Type,
	}
	switch in.Type {
	case stream.Twitch:
		if in.Twitch.Line.Cmd != "PRIVMSG" {
			return
		}
		if len(in.Twitch.Line.Args) < 2 {
			return
		}
		msg := in.Twitch.Line.Args[1]
		if !strings.HasPrefix(msg, e.cmd+" ") {
			return
		}
		msg = msg[len(e.cmd)+1:]
		out.Twitch = &stream.TXTwitch{
			Username: e.twitchUsername,
			To:       "#jtv",
			Message:  fmt.Sprintf("/w %s %s", in.Twitch.Line.Nick, msg),
		}
	case stream.Discord:
		// TODO: validation of discord in message
		out.Discord = &stream.TXDiscord{
			To:      in.Discord.MessageCreate.Author.Username,
			Message: in.Discord.MessageCreate.Content,
		}
	}
	e.sman.Send(out)
}

// Start is a NOOP.
func (e *EchoFeature) Start() {}

// Stop is a NOOP.
func (e *EchoFeature) Stop() {}
