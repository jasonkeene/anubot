package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/fluffle/goirc/logging/golog"
	"github.com/spf13/viper"

	"anubot/stream"
	"anubot/twitch"
)

func init() {
	golog.Init()
}

type fakeDispatcher struct {
	mu              sync.Mutex
	lastDiscordChan string
}

func (fd *fakeDispatcher) Dispatch(topic string, m stream.RXMessage) {
	switch m.Type {
	case stream.Twitch:
		fmt.Printf("TWITCH: %s\n", m.Twitch.Line.Text())
	case stream.Discord:
		fmt.Printf("DISCORD: %s\n", m.Discord.MessageCreate.Content)
		fd.mu.Lock()
		defer fd.mu.Unlock()
		fd.lastDiscordChan = m.Discord.MessageCreate.ChannelID
	}
}

func (fd *fakeDispatcher) LastDiscordChannel() string {
	fd.mu.Lock()
	defer fd.mu.Unlock()
	return fd.lastDiscordChan
}

func main() {
	// load config
	v := viper.New()
	v.SetEnvPrefix("anubot")
	v.AutomaticEnv()

	c := "#" + os.Getenv("TWITCH_USER_USER")
	u := os.Getenv("TWITCH_BOT_USER")
	p := os.Getenv("TWITCH_BOT_PASS")
	t := os.Getenv("DISCORD_BOT_PASS")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	d := &fakeDispatcher{}
	twitch := twitch.New(
		v.GetString("twitch_api_url"),
		v.GetString("twitch_oauth_client_id"),
	)

	manager := stream.NewManager(d, twitch)
	manager.ConnectTwitch(u, p, c)
	manager.ConnectDiscord(t)

	go func() {
		r := bufio.NewReader(os.Stdin)
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				panic(err)
			}
			if strings.TrimSpace(string(line)) == "" {
				continue
			}
			manager.Send(stream.TXMessage{
				Type: stream.Twitch,
				Twitch: &stream.TXTwitch{
					Username: u,
					To:       c,
					Message:  string(line),
				},
			})
			dchan := d.LastDiscordChannel()
			if dchan != "" {
				manager.Send(stream.TXMessage{
					Type: stream.Discord,
					Discord: &stream.TXDiscord{
						To:      dchan,
						Message: string(line),
					},
				})
			}
		}
	}()

	<-interrupt
	wait := manager.DisconnectTwitch(u)
	wait()
	wait = manager.DisconnectDiscord()
	wait()
}
