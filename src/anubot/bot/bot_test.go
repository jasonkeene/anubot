package bot_test

import (
	"crypto/tls"
	"net"
	"strconv"

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
			It("can join a channel", func() {
				fakeIRCServer.Respond(0,
					":test-streamer-user!test-streamer-user@test-streamer-user.127.0.0.1 JOIN #test_chan",
					":test-streamer-user.127.0.0.1 353 test-streamer-user = #test_chan :test-streamer-user",
					":test-streamer-user.127.0.0.1 366 test-streamer-user #test_chan :End of /NAMES list",
				)
				fakeIRCServer.Respond(1,
					":test-bot-user!test-bot-user@test-bot-user.127.0.0.1 JOIN #test_chan",
					":test-bot-user.127.0.0.1 353 test-bot-user = #test_chan :test-bot-user",
					":test-bot-user.127.0.0.1 366 test-bot-user #test_chan :End of /NAMES list",
				)

				bot.Join("#test_chan")

				Eventually(fakeIRCServer.Received(0)).Should(ContainLines(
					"JOIN #test_chan",
				))
				Eventually(fakeIRCServer.Received(1)).Should(ContainLines(
					"JOIN #test_chan",
				))
			})
		})

		Describe("Send", func() {
			It("sends messages to the server", func() {
				bot.Send("streamer", "test-streamer-message")
				Eventually(fakeIRCServer.Received(0), 3).Should(EqualLines(
					"test-streamer-message",
				))
				bot.Send("bot", "test-bot-message")
				Eventually(fakeIRCServer.Received(1), 3).Should(EqualLines(
					"test-bot-message",
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
