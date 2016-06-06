package bot_test

import (
	"crypto/tls"
	"net"
	"strconv"
	"sync"

	"github.com/fluffle/goirc/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "anubot/bot"
)

var _ = Describe("Bot", func() {
	var (
		bot           *Bot
		connConfig    *ConnConfig
		fakeIRCServer *fakeIRCServer
	)

	BeforeEach(func() {
		// create tls listener for the fake server
		listener, err := tls.Listen("tcp", ":0", serverTLSConfig)
		Expect(err).ToNot(HaveOccurred())

		// start up fake IRC server
		fakeIRCServer = newFakeIRCServer(listener)
		go func() {
			defer GinkgoRecover()
			fakeIRCServer.Start()
		}()

		// figure out the port that was assigned
		_, sport, err := net.SplitHostPort(listener.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		port, err := strconv.Atoi(sport)
		Expect(err).ToNot(HaveOccurred())

		// create the bot and server
		bot = &Bot{}
		connConfig = &ConnConfig{
			StreamerUsername: "test-streamer-user",
			StreamerPassword: "test-streamer-password",
			BotUsername:      "test-bot-user",
			BotPassword:      "test-bot-password",
			Host:             "127.0.0.1",
			Port:             port,
			TLSConfig:        clientTLSConfig,
			Flood:            true,
		}
	})

	AfterEach(func() {
		bot.Disconnect()
		fakeIRCServer.Stop()
	})

	Describe("Connect/Disconnect", func() {
		BeforeEach(func() {
			fakeIRCServer.Respond(0, ":127.0.0.1 001 test-streamer-user :GLHF!")
			fakeIRCServer.Respond(1, ":127.0.0.1 001 test-bot-user :GLHF!")
		})

		It("attempts to connect over TLS", func() {
			_, err := bot.Connect(connConfig)
			Expect(err).ToNot(HaveOccurred())
			assertConnected(0, "test-streamer-user", "test-streamer-password", fakeIRCServer)
			assertConnected(1, "test-bot-user", "test-bot-password", fakeIRCServer)
			fakeIRCServer.Clear()
		})

		It("joins the streamer's channel on connect", func() {
			_, err := bot.Connect(connConfig)
			Expect(bot.Channel()).To(Equal("#test-streamer-user"))
			Expect(err).ToNot(HaveOccurred())

			assertConnected(0, "test-streamer-user", "test-streamer-password", fakeIRCServer)
			assertConnected(1, "test-bot-user", "test-bot-password", fakeIRCServer)

			Eventually(fakeIRCServer.Received(0)).Should(ContainLines(
				"JOIN #test-streamer-user",
			))
			Eventually(fakeIRCServer.Received(1)).Should(ContainLines(
				"JOIN #test-streamer-user",
			))
		})

		It("can disconnect", func() {
			disconnected, err := bot.Connect(connConfig)
			Expect(err).ToNot(HaveOccurred())
			assertConnected(0, "test-streamer-user", "test-streamer-password", fakeIRCServer)
			assertConnected(1, "test-bot-user", "test-bot-password", fakeIRCServer)
			bot.Disconnect()
			Eventually(disconnected).Should(BeClosed())
		})
	})

	Context("with a connected bot", func() {
		BeforeEach(func() {
			fakeIRCServer.Respond(0, ":127.0.0.1 001 test-streamer-user :GLHF!")
			fakeIRCServer.Respond(1, ":127.0.0.1 001 test-bot-user :GLHF!")

			_, err := bot.Connect(connConfig)
			Expect(err).ToNot(HaveOccurred())
			assertConnected(0, "test-streamer-user", "test-streamer-password", fakeIRCServer)
			assertConnected(1, "test-bot-user", "test-bot-password", fakeIRCServer)
			fakeIRCServer.Clear()
		})

		Describe("Join", func() {
			It("can join a different channel", func() {
				bot.Join("#test_chan")
				Expect(bot.Channel()).To(Equal("#test_chan"))

				Eventually(fakeIRCServer.Received(0)).Should(ContainLines(
					"JOIN #test_chan",
				))
				Eventually(fakeIRCServer.Received(1)).Should(ContainLines(
					"JOIN #test_chan",
				))
			})

			It("parts old channel when it joins a different channel", func() {
				bot.Join("#test_chan")

				Eventually(fakeIRCServer.Received(0)).Should(ContainLines(
					"PART #test-streamer-user",
				))
				Eventually(fakeIRCServer.Received(1)).Should(ContainLines(
					"PART #test-streamer-user",
				))

				bot.Join("#test_chan2")
				Expect(bot.Channel()).To(Equal("#test_chan2"))

				Eventually(fakeIRCServer.Received(0)).Should(ContainLines(
					"PART #test_chan",
				))
				Eventually(fakeIRCServer.Received(1)).Should(ContainLines(
					"PART #test_chan",
				))
			})
		})

		Describe("Send", func() {
			It("sends raw messages to the specified connection", func() {
				bot.Send("streamer", "test-streamer-message")
				Eventually(fakeIRCServer.Received(0)).Should(ContainLines(
					"test-streamer-message",
				))
				bot.Send("bot", "test-bot-message")
				Eventually(fakeIRCServer.Received(1)).Should(ContainLines(
					"test-bot-message",
				))
			})
		})

		Describe("Privmsg", func() {
			It("sends chat messages to the specified connection", func() {
				bot.Privmsg("streamer", "#test-streamer-user", "test-streamer-message")
				Eventually(fakeIRCServer.Received(0)).Should(ContainLines(
					"PRIVMSG #test-streamer-user :test-streamer-message",
				))
				bot.Privmsg("bot", "#test-streamer-user", "test-bot-message")
				Eventually(fakeIRCServer.Received(1)).Should(ContainLines(
					"PRIVMSG #test-streamer-user :test-bot-message",
				))
			})
		})

		Describe("HandleFunc", func() {
			var (
				mu     sync.Mutex
				called bool
			)

			It("can register functions to handle events for the streamer", func() {
				fakeIRCServer.Respond(0, ":thebossreturns!thebossreturns@thebossreturns.tmi.twitch.tv PRIVMSG #postcrypt :hello world")
				bot.HandleFunc("streamer", "PRIVMSG", func(*client.Conn, *client.Line) {
					mu.Lock()
					defer mu.Unlock()
					called = true
				})
				Eventually(func() bool {
					mu.Lock()
					defer mu.Unlock()
					return called
				}).Should(BeTrue())
			})

			It("can register functions to handle events for the bot", func() {
				fakeIRCServer.Respond(1, ":thebossreturns!thebossreturns@thebossreturns.tmi.twitch.tv PRIVMSG #postcrypt :hello world")
				bot.HandleFunc("bot", "PRIVMSG", func(*client.Conn, *client.Line) {
					mu.Lock()
					defer mu.Unlock()
					called = true
				})
				Eventually(func() bool {
					mu.Lock()
					defer mu.Unlock()
					return called
				}).Should(BeTrue())
			})
		})

		Describe("InitChatFeature", func() {
			It("creates a chat feature from a message dispatcher", func() {
				bot.InitChatFeature(nil)
				cf := bot.ChatFeature()
				cf.Send("streamer", "test-cf-message")

				Eventually(fakeIRCServer.Received(0)).Should(ContainLines(
					"PRIVMSG #test-streamer-user :test-cf-message",
				))
			})
		})
	})
})

func assertConnected(connIndex int, username, password string, fakeIRCServer *fakeIRCServer) {
	Eventually(fakeIRCServer.Received(connIndex)).Should(ContainLines(
		"PASS "+password,
		"NICK "+username,
		"USER anubot 12 * :"+username,
	))
	Eventually(fakeIRCServer.Sent(connIndex)).Should(ContainLines(
		":127.0.0.1 001 " + username + " :GLHF!",
	))
}
