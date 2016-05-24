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

		// figure out the port that was assigned
		_, sport, err := net.SplitHostPort(listener.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		port, err := strconv.Atoi(sport)
		Expect(err).ToNot(HaveOccurred())

		// create the bot and server
		bot = &Bot{}
		connConfig = &ConnConfig{
			BotUsername: "test_user",
			BotPassword: "test_password",
			Host:        "127.0.0.1",
			Port:        port,
			TLSConfig:   clientTLSConfig,
		}
		fakeIRCServer = newFakeIRCServer(listener)
		go func() {
			defer GinkgoRecover()
			fakeIRCServer.Start()
		}()
	})

	AfterEach(func() {
		bot.Disconnect()
		fakeIRCServer.Stop()
	})

	Describe("Connect/Disconnect", func() {
		BeforeEach(func() {
			fakeIRCServer.Respond(":127.0.0.1 001 test_user :GLHF!")
		})

		It("attempts to connect over TLS", func() {
			_, err := bot.Connect(connConfig)
			Expect(err).ToNot(HaveOccurred())
			assertConnected(fakeIRCServer)
		})

		It("can disconnect", func() {
			disconnected, err := bot.Connect(connConfig)
			Expect(err).ToNot(HaveOccurred())
			assertConnected(fakeIRCServer)
			bot.Disconnect()
			Eventually(disconnected).Should(BeClosed())
		})
	})

	Describe("Join", func() {
		Context("with a connected bot", func() {
			BeforeEach(func() {
				fakeIRCServer.Respond(":127.0.0.1 001 test_user :GLHF!")

				_, err := bot.Connect(connConfig)
				Expect(err).ToNot(HaveOccurred())
				assertConnected(fakeIRCServer)
			})

			It("can join multiple channels", func() {
				fakeIRCServer.Respond(
					":test_user!test_user@test_user.127.0.0.1 JOIN #test_chan1",
					":test_user.127.0.0.1 353 test_user = #test_chan1 :test_user",
					":test_user.127.0.0.1 366 test_user #test_chan1 :End of /NAMES list",
				)
				fakeIRCServer.Respond(
					":test_user!test_user@test_user.127.0.0.1 JOIN #test_chan2",
					":test_user.127.0.0.1 353 test_user = #test_chan2 :test_user",
					":test_user.127.0.0.1 366 test_user #test_chan2 :End of /NAMES list",
				)

				bot.Join("#test_chan1")
				bot.Join("#test_chan2")
				Eventually(fakeIRCServer.Received, 3).Should(EqualLines(
					"JOIN #test_chan1",
					"JOIN #test_chan2",
				))
			})
		})
	})
})

func assertConnected(fakeIRCServer *fakeIRCServer) {
	Eventually(fakeIRCServer.Received).Should(EqualLines(
		"PASS test_password",
		"NICK test_user",
		"USER anubot 12 * :test_user",
	))
	Eventually(fakeIRCServer.Sent).Should(EqualLines(
		":127.0.0.1 001 test_user :GLHF!",
	))
	fakeIRCServer.Clear()
}
