package main

import (
	"os"

	"github.com/fluffle/goirc/logging"
	logginglib "github.com/fluffle/golog/logging"

	anubot "anubot/bot"
)

const (
	twitchHost = "irc.chat.twitch.tv"
	twitchPort = 443
)

func main() {
	twitchUser := os.Getenv("TWITCH_USER")
	twitchPass := os.Getenv("TWITCH_PASS")
	twitchChannel := os.Getenv("TWITCH_CHANNEL")
	initLogging()
	bot := anubot.New(
		twitchUser,
		twitchPass,
		twitchHost,
		twitchPort,
	)
	err, disconnected := bot.Connect(nil)
	if err != nil {
		panic(err)
	}
	bot.Join(twitchChannel)

	<-disconnected
}

func initLogging() {
	logger := logginglib.NewFromFlags()
	logger.SetLogLevel(logginglib.LogDebug)
	logging.SetLogger(logger)
}
