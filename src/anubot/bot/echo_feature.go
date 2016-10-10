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
		msg := e.matchMessage(in.Twitch.Line.Args[1])
		if msg == "" {
			return
		}
		out.Twitch = &stream.TXTwitch{
			Username: e.twitchUsername,
			To:       "#jtv",
			Message:  fmt.Sprintf("/w %s %s", in.Twitch.Line.Nick, msg),
		}
	case stream.Discord:
		msg := e.matchMessage(in.Discord.MessageCreate.Content)
		if msg == "" {
			return
		}
		out.Discord = &stream.TXDiscord{
			Type:    stream.Private,
			To:      in.Discord.MessageCreate.Author.ID,
			Message: msg,
		}
	}
	e.sman.Send(out)
}

func (e *EchoFeature) matchMessage(msg string) string {
	if !strings.HasPrefix(msg, e.cmd+" ") {
		return ""
	}
	return msg[len(e.cmd)+1:]
}

// Start is a NOOP.
func (e *EchoFeature) Start() {}

// Stop is a NOOP.
func (e *EchoFeature) Stop() {}
