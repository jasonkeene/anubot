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
	// load config
	v := viper.New()
	v.SetEnvPrefix("anubot")
	v.AutomaticEnv()

	twitchUserUsername := os.Getenv("TWITCH_USER_USER")
	twitchUserPassword := os.Getenv("TWITCH_USER_PASS")
	twitchChannel := "#" + twitchUserUsername

	twitchBotUsername := os.Getenv("TWITCH_BOT_USER")
	twitchBotPassword := os.Getenv("TWITCH_BOT_PASS")

	discrodUserID := os.Getenv("DISCORD_USER_ID")

	discordBotPassword := os.Getenv("DISCORD_BOT_PASS")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	d := dispatch.New(
		[]string{"inproc://pub"},
		[]string{"inproc://push"},
	)
	twitch := twitch.New(
		v.GetString("twitch_api_url"),
		v.GetString("twitch_oauth_client_id"),
	)
	manager := stream.NewManager(d, twitch)
	manager.ConnectTwitch(twitchBotUsername, twitchBotPassword, twitchChannel)
	manager.ConnectTwitch(twitchUserUsername, twitchUserPassword, twitchChannel)
	manager.ConnectDiscord("Bot " + discordBotPassword)

	pull, err := zmq4.NewSocket(zmq4.PULL)
	if err != nil {
		log.Panicf("pull not created, got err: %s", err)
	}
	err = pull.Connect("inproc://push")
	if err != nil {
		log.Panicf("pull not able to connect, got err: %s", err)
	}

	b, err := bot.New(
		[]string{
			"twitch:" + twitchBotUsername,
			"discord:" + discrodUserID,
		},
		[]string{"inproc://pub"},
	)
	if err != nil {
		panic(err)
	}
	go b.Start()
	defer b.Stop()

	f := bot.NewEchoFeature("!echo", twitchBotUsername, manager)
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
	wait := manager.DisconnectTwitch(twitchBotUsername)
	wait()
	wait = manager.DisconnectDiscord()
	wait()
}
