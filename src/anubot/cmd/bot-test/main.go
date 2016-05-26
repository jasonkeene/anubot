package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/fluffle/goirc/logging"
	logginglib "github.com/fluffle/golog/logging"

	"anubot/bot"
)

const (
	twitchHost = "irc.chat.twitch.tv"
	twitchPort = 443
)

func main() {
	initLogging()
	connConfig := &bot.ConnConfig{
		UserUsername: os.Getenv("TWITCH_USER_USER"),
		UserPassword: os.Getenv("TWITCH_USER_PASS"),
		BotUsername:  os.Getenv("TWITCH_BOT_USER"),
		BotPassword:  os.Getenv("TWITCH_BOT_PASS"),
		Host:         twitchHost,
		Port:         twitchPort,
	}
	bot := &bot.Bot{}
	disconnected, err := bot.Connect(connConfig)
	if err != nil {
		panic(err)
	}
	bot.Join("#" + connConfig.UserUsername)

	// read from stdin and send to irc server
	reader := bufio.NewReader(os.Stdin)
	go func() {
		for {
			text, err := reader.ReadString('\n')
			if err != nil {
				continue
			}
			text = strings.Trim(text, "\r\n")
			if len(text) > 0 {
				bot.Send(text)
			}
		}
	}()
	<-disconnected
}

func initLogging() {
	logger := logginglib.NewFromFlags()
	logger.SetLogLevel(logginglib.LogDebug)
	logging.SetLogger(logger)
}
