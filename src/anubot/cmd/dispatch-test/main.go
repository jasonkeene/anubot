package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/fluffle/goirc/logging/golog"
	"github.com/pebbe/zmq4"
	"github.com/spf13/viper"

	"anubot/dispatch"
	"anubot/stream"
	"anubot/twitch"
)

func init() {
	golog.Init()
}

func main() {
	v := viper.New()
	v.SetEnvPrefix("anubot")
	v.AutomaticEnv()
	twitchStreamerUser := os.Getenv("TWITCH_USER_USER")
	twitchStreamerPass := os.Getenv("TWITCH_USER_PASS")
	twitchStreamerChannel := "#" + twitchStreamerUser
	twitchStreamerTopic := "twitch:" + twitchStreamerUser
	twitchBotUser := os.Getenv("TWITCH_BOT_USER")
	twitchBotPass := os.Getenv("TWITCH_BOT_PASS")
	discordBotPass := os.Getenv("DISCORD_BOT_PASS")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	dispatch.Start()
	twitch := twitch.New(
		v.GetString("twitch_api_url"),
		v.GetString("twitch_oauth_client_id"),
	)
	manager := stream.NewManager(twitch)
	manager.ConnectTwitch(
		twitchBotUser,
		twitchBotPass,
		twitchStreamerChannel,
	)
	manager.ConnectTwitch(
		twitchStreamerUser,
		twitchStreamerPass,
		twitchStreamerChannel,
	)
	manager.ConnectDiscord(discordBotPass)

	sub := createSub("inproc://dispatch-pub", "")
	topicSub := createSub("inproc://dispatch-pub", twitchStreamerTopic)
	pull := createPull("inproc://dispatch-push")

	go readFromSub(sub)
	go readFromSub(topicSub)
	go readFromPull(pull)

	<-interrupt
	twitchWait := manager.DisconnectTwitch(twitchBotUser)
	discordWait := manager.DisconnectDiscord()
	twitchWait()
	discordWait()
}

func readFromSub(sub *zmq4.Socket) {
	for {
		rb, err := sub.RecvMessageBytes(0)
		if err != nil {
			log.Printf("messages not read, got err: %s", err)
			continue
		}
		topic := rb[0]
		var message stream.RXMessage
		err = json.Unmarshal(rb[1], &message)
		if err != nil {
			log.Printf("could not unmarshal, got err: %s", err)
			continue
		}
		fmt.Printf("sub with topic: %s got message: %#v\n", topic, message)
	}
}

func readFromPull(pull *zmq4.Socket) {
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
		fmt.Printf("pull got message: %#v\n", message)
	}
}

func createSub(connect, topic string) *zmq4.Socket {
	sub, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		log.Panicf("sub not created, got err: %s", err)
	}
	err = sub.Connect(connect)
	if err != nil {
		log.Panicf("sub not able to connect, got err: %s", err)
	}
	err = sub.SetSubscribe(topic)
	if err != nil {
		log.Panicf("sub not able to set topic, got err: %s", err)
	}
	return sub
}

func createPull(connect string) *zmq4.Socket {
	pull, err := zmq4.NewSocket(zmq4.PULL)
	if err != nil {
		log.Panicf("pull not created, got err: %s", err)
	}
	err = pull.Connect(connect)
	if err != nil {
		log.Panicf("pull not able to connect, got err: %s", err)
	}
	return pull
}
