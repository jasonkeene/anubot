package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/fluffle/goirc/logging/golog"
	"github.com/pebbe/zmq4"

	"anubot/bot"
	"anubot/dispatch"
	"anubot/stream"
)

func init() {
	golog.Init()
}

type fakeStore struct{}

func (fakeStore) TwitchUser(channelName string) (string, error) {
	println("twitch users called with:", channelName)
	return "postcrypt", nil
}

func (fakeStore) DiscordUsers(channelID string) []string {
	println("discord users called with:", channelID)
	return []string{"postcrypt"}
}

func main() {
	uu := os.Getenv("TWITCH_USER_USER")
	up := os.Getenv("TWITCH_USER_PASS")
	c := "#" + uu

	u := os.Getenv("TWITCH_BOT_USER")
	p := os.Getenv("TWITCH_BOT_PASS")

	t := os.Getenv("DISCORD_BOT_PASS")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	s := fakeStore{}
	d := dispatch.New(
		[]string{"inproc://pub"},
		[]string{"inproc://push"},
		s,
	)
	manager := stream.NewManager(d)
	manager.ConnectTwitch(u, p, c)
	manager.ConnectTwitch(uu, up, c)
	manager.ConnectDiscord(t)

	pull, err := zmq4.NewSocket(zmq4.PULL)
	if err != nil {
		log.Panicf("pull not created, got err: %s", err)
	}
	err = pull.Connect("inproc://push")
	if err != nil {
		log.Panicf("pull not able to connect, got err: %s", err)
	}

	b, err := bot.New(
		"postcrypt",
		[]string{"inproc://pub"},
	)
	if err != nil {
		panic(err)
	}
	go b.Start()
	defer b.Stop()

	f := bot.NewEchoFeature("!echo", u, manager)
	b.SetFeature("echo", f)

	go func() {
		for {
			rb, err := pull.RecvBytes(0)
			if err != nil {
				log.Printf("messages not read, got err: %s", err)
				continue
			}
			var message stream.RXMessage
			err = json.Unmarshal(rb, &message)
			if err != nil {
				log.Printf("could not unmarshal, got err: %s", err)
				continue
			}
		}
	}()

	<-interrupt
	wait := manager.DisconnectTwitch(u)
	wait()
	wait = manager.DisconnectDiscord()
	wait()
}
