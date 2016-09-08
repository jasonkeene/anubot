package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/davecgh/go-spew/spew"
	"github.com/fluffle/goirc/logging/golog"
	"github.com/pebbe/zmq4"

	"anubot/dispatch"
	"anubot/stream"
)

func init() {
	golog.Init()
}

type fakeStore struct{}

func (fakeStore) TwitchUser(channelName string) (string, error) {
	println("twitch users called with:", channelName)
	return "foo", nil
}

func (fakeStore) DiscordUsers(channelID string) []string {
	println("discord users called with:", channelID)
	return []string{"foo"}
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
	d := dispatch.New([]string{"inproc://pub"}, []string{"inproc://push"}, s)
	manager := stream.NewManager(d)
	manager.ConnectTwitch(u, p, c)
	manager.ConnectTwitch(uu, up, c)
	manager.ConnectDiscord(t)

	sub, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		log.Panicf("sub not created, got err: %s", err)
	}
	err = sub.Connect("inproc://pub")
	if err != nil {
		log.Panicf("sub not able to connect, got err: %s", err)
	}
	err = sub.SetSubscribe("")
	if err != nil {
		log.Panicf("topicSub not able to set topic, got err: %s", err)
	}

	topicSub, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		log.Panicf("topicSub not created, got err: %s", err)
	}
	err = topicSub.Connect("inproc://pub")
	if err != nil {
		log.Panicf("topicSub not able to connect, got err: %s", err)
	}
	err = topicSub.SetSubscribe("cat")
	if err != nil {
		log.Panicf("topicSub not able to set topic, got err: %s", err)
	}

	pull, err := zmq4.NewSocket(zmq4.PULL)
	if err != nil {
		log.Panicf("pull not created, got err: %s", err)
	}
	err = pull.Connect("inproc://push")
	if err != nil {
		log.Panicf("pull not able to connect, got err: %s", err)
	}

	//go func() {
	//	for {
	//		rb, err := sub.RecvMessageBytes(0)
	//		if err != nil {
	//			log.Printf("messages not read, got err: %s", err)
	//			continue
	//		}
	//		var message stream.RXMessage
	//		err = json.Unmarshal(rb[1], &message)
	//		if err != nil {
	//			log.Printf("could not unmarshal, got err: %s", err)
	//			continue
	//		}
	//		fmt.Print("sub got message:")
	//		spew.Dump(message)
	//	}
	//}()

	//go func() {
	//	for {
	//		rb, err := topicSub.RecvMessageBytes(0)
	//		if err != nil {
	//			log.Printf("messages not read, got err: %s", err)
	//			continue
	//		}
	//		var message stream.RXMessage
	//		err = json.Unmarshal(rb[1], &message)
	//		if err != nil {
	//			log.Printf("could not unmarshal, got err: %s", err)
	//			continue
	//		}
	//		fmt.Print("topicSub got message:")
	//		spew.Dump(message)
	//	}
	//}()

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
			fmt.Print("pull got message:")
			spew.Dump(message)
		}
	}()

	<-interrupt
	wait := manager.DisconnectTwitch(u)
	wait()
	wait = manager.DisconnectDiscord()
	wait()
}
