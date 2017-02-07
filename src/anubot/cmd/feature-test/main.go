package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/fluffle/goirc/logging/golog"
	"github.com/pebbe/zmq4"
	"github.com/spf13/viper"

	"anubot/bot"
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
	twitchBotUser := os.Getenv("TWITCH_BOT_USER")
	twitchBotPass := os.Getenv("TWITCH_BOT_PASS")
	discordUserID := os.Getenv("DISCORD_USER_ID")
	discordBotPass := os.Getenv("DISCORD_BOT_PASS")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	dispatch.Start()

	// create puller that reads off of dispatcher's push endpoint
	pull := createPull("inproc://discord-push")
	go readFromPull(pull)

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
	manager.ConnectDiscord("Bot " + discordBotPass)

	b, err := bot.New(
		[]string{
			"twitch:" + twitchBotUser,
			"discord:" + discordUserID,
		},
	)
	if err != nil {
		panic(err)
	}
	go b.Start()
	defer b.Stop()

	f := bot.NewEchoFeature("!echo", twitchBotUser, manager)
	b.SetFeature("echo", f)

	<-interrupt
	twitchWait := manager.DisconnectTwitch(twitchBotUser)
	discordWait := manager.DisconnectDiscord()
	twitchWait()
	discordWait()
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
	}
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
