package stream

import (
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
)

type discordConn struct {
	d  Dispatcher
	dg *discordgo.Session
}

func (c *discordConn) send(m TXMessage) {
	if m.Discord == nil {
		log.Printf("Attempted to TX a nil Discord message")
		return
	}

	switch m.Discord.Type {
	case Channel:
		_, err := c.dg.ChannelMessageSend(m.Discord.To, m.Discord.Message)
		if err != nil {
			log.Printf("discordConn.send: error sending channel message: %s", err)
		}
	case Private:
		ch, err := c.dg.UserChannelCreate(m.Discord.To)
		if err != nil {
			log.Printf("discordConn.send: error getting DM channel: %s", err)
		}
		_, err = c.dg.ChannelMessageSend(ch.ID, m.Discord.Message)
		if err != nil {
			log.Printf("discordConn.send: error sending DM message: %s", err)
		}
	default:
		log.Printf("discordConn.send: Attempted to TX a Discord message with unknown type: %d", m.Discord.Type)
		return
	}
}

func (c *discordConn) close() error {
	return c.dg.Close()
}

func connectDiscord(token string, d Dispatcher) (*discordConn, error) {
	dg, err := discordgo.New(token)
	dg.LogLevel = discordgo.LogDebug
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
	ownerID, err := c.resolveOwnerID(m)
	if err != nil {
		log.Printf("got err attempting to resolve discord topic: %s", err)
	}
	topic := "discord:" + ownerID
	c.d.Dispatch(topic, RXMessage{
		Type: Discord,
		Discord: &RXDiscord{
			OwnerID:       ownerID,
			MessageCreate: m,
		},
	})
}

func (c *discordConn) resolveOwnerID(m *discordgo.MessageCreate) (string, error) {
	ch, err := c.dg.Channel(m.ChannelID)
	if err != nil {
		return "", err
	}
	if ch.IsPrivate {
		return "", errors.New("Not possible to resolve private messages back to guild owner")
	}
	gld, err := c.dg.Guild(ch.GuildID)
	if err != nil {
		return "", err
	}
	return gld.OwnerID, nil
}
