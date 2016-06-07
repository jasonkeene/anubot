package api_test

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/websocket"

	. "github.com/apoydence/eachers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "anubot/api"
	"anubot/bot"
)

var _ = Describe("APIServer", func() {
	var (
		mockStore *mockStore
		mockBot   *mockBot
		listener  net.Listener
		api       *APIServer
		server    *http.Server
		client    *websocket.Conn
	)

	BeforeEach(func() {
		mockStore = newMockStore()
		mockBot = newMockBot()

		// spin up server
		var err error
		listener, err = net.Listen("tcp", ":0")
		Expect(err).ToNot(HaveOccurred())

		api = New(mockStore, mockBot, nil)
		server = &http.Server{
			Handler:        websocket.Handler(api.Serve),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		go func() {
			defer GinkgoRecover()
			server.Serve(listener)
		}()
		_, port, err := net.SplitHostPort(listener.Addr().String())
		Expect(err).ToNot(HaveOccurred())

		// connect client
		origin := "http://localhost/"
		url := fmt.Sprintf("ws://localhost:%s/", port)
		client, err = websocket.Dial(url, "", origin)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		client.Close()
		listener.Close()
	})

	It("responds to ping commands", func() {
		event := Event{Cmd: "ping"}
		websocket.JSON.Send(client, &event)
		websocket.JSON.Receive(client, &event)
		Expect(event.Cmd).To(Equal("pong"))
	})

	It("can respond concurrently", func() {
		event := Event{Cmd: "ping"}
		websocket.JSON.Send(client, &event)
		websocket.JSON.Send(client, &event)
		websocket.JSON.Receive(client, &event)
		Expect(event.Cmd).To(Equal("pong"))
		websocket.JSON.Receive(client, &event)
		Expect(event.Cmd).To(Equal("pong"))
	})

	It("ignores invalid commands", func() {
		event := Event{Cmd: "invalid-cmd"}
		websocket.JSON.Send(client, &event)

		receive := make(chan *Event, 1)
		go func() {
			websocket.JSON.Receive(client, &event)
			receive <- &event
		}()
		Consistently(receive).ShouldNot(Receive())
	})

	It("can manage credentials", func() {
		// make sure credentails are not set
		mockStore.HasCredentialsOutput.Has <- false
		event := Event{
			Cmd:     "has-credentials-set",
			Payload: "user",
		}
		websocket.JSON.Send(client, &event)
		websocket.JSON.Receive(client, &event)
		Expect(event.Cmd).To(Equal("has-credentials-set"))
		payload := event.Payload.(map[string]interface{})
		kind := payload["kind"].(string)
		result := payload["result"].(bool)
		Expect(kind).To(Equal("user"))
		Expect(result).To(BeFalse())
		Expect(mockStore.HasCredentialsCalled).To(Receive())

		// set credentials
		mockStore.SetCredentialsOutput.Err <- nil
		event = Event{
			Cmd: "set-credentials",
			Payload: map[string]string{
				"kind":     "user",
				"username": "test-username",
				"password": "test-password",
			},
		}
		websocket.JSON.Send(client, &event)
		Eventually(mockStore.SetCredentialsInput).Should(BeCalled(
			With("user", "test-username", "test-password"),
		))

		// make sure they were set correctly
		mockStore.HasCredentialsOutput.Has <- true
		event = Event{
			Cmd:     "has-credentials-set",
			Payload: "user",
		}
		websocket.JSON.Send(client, &event)
		websocket.JSON.Receive(client, &event)
		Expect(event.Cmd).To(Equal("has-credentials-set"))
		payload = event.Payload.(map[string]interface{})
		kind = payload["kind"].(string)
		result = payload["result"].(bool)
		Expect(kind).To(Equal("user"))
		Expect(result).To(BeTrue())
		Expect(mockStore.HasCredentialsCalled).To(Receive())
	})

	It("ignores credential kinds other than user and bot", func() {
		// attempt to set credentials
		event := Event{
			Cmd: "set-credentials",
			Payload: map[string]string{
				"kind":     "bad-kind",
				"username": "test-username",
				"password": "test-password",
			},
		}
		websocket.JSON.Send(client, &event)
		Consistently(mockStore.SetCredentialsInput).ShouldNot(BeCalled())

		// make sure has credentials doesn't call into the store either
		event = Event{
			Cmd:     "has-credentials-set",
			Payload: "bad-kind",
		}
		websocket.JSON.Send(client, &event)
		Consistently(mockStore.HasCredentialsInput).ShouldNot(BeCalled())
		websocket.JSON.Receive(client, &event)
		Expect(event.Cmd).To(Equal("has-credentials-set"))
		payload := event.Payload.(map[string]interface{})
		kind := payload["kind"].(string)
		result := payload["result"].(bool)
		Expect(kind).To(Equal("bad-kind"))
		Expect(result).To(BeFalse())
	})

	It("can have the bot connect and disconnect", func() {
		// connect
		mockStore.CredentialsOutput.User <- "user-test-username"
		mockStore.CredentialsOutput.Pass <- "user-test-password"
		mockStore.CredentialsOutput.Err <- nil
		mockStore.CredentialsOutput.User <- "bot-test-username"
		mockStore.CredentialsOutput.Pass <- "bot-test-password"
		mockStore.CredentialsOutput.Err <- nil

		mockBot.ConnectOutput.Disconnected <- nil
		mockBot.ConnectOutput.Err <- nil

		event := Event{Cmd: "connect"}
		websocket.JSON.Send(client, &event)
		websocket.JSON.Receive(client, &event)
		Expect(event.Cmd).To(Equal("connect"))
		connected, ok := event.Payload.(bool)
		Expect(ok).To(BeTrue())
		Expect(connected).To(BeTrue())
		Expect(mockStore.CredentialsInput).To(BeCalled(
			With("user"),
			With("bot"),
		))

		var connConfig *bot.ConnConfig
		Eventually(mockBot.ConnectInput.ConnConfig).Should(Receive(&connConfig))
		Expect(connConfig.StreamerUsername).To(Equal("user-test-username"))
		Expect(connConfig.StreamerPassword).To(Equal("user-test-password"))
		Expect(connConfig.BotUsername).To(Equal("bot-test-username"))
		Expect(connConfig.BotPassword).To(Equal("bot-test-password"))
		Expect(connConfig.Host).To(Equal("irc.chat.twitch.tv"))
		Expect(connConfig.Port).To(Equal(443))

		// disconnect
		event = Event{Cmd: "disconnect"}
		websocket.JSON.Send(client, &event)
		Eventually(mockBot.DisconnectCalled).Should(Receive())
	})
})
