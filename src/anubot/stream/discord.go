package stream

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type discordConn struct {
	d  Dispatcher
	dg *discordgo.Session
}

func (c *discordConn) send(m TXMessage) {
	_, err := c.dg.ChannelMessageSend(m.Discord.To, m.Discord.Message)
	if err != nil {
		log.Printf("discordConn.send: error occured: %s", err)
	}
}

func (c *discordConn) close() error {
	return c.dg.Close()
}

func connectDiscord(token string, d Dispatcher) (*discordConn, error) {
	dg, err := discordgo.New(token)
	dg.LogLevel = discordgo.LogInformational
	if err != nil {
		return nil, err
	}
	err = dg.Open()
	if err != nil {
		return nil, err
	}
	dc := &discordConn{
		d:  d,
		dg: dg,
	}
	dg.AddHandler(dc.messageCreate)
	return dc, nil
}

func (c *discordConn) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO: figure out topic to dispatch
	c.d.Dispatch("", RXMessage{
		Type: Discord,
		Discord: &RXDiscord{
			MessageCreate: m,
		},
	})
}
