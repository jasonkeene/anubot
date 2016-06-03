package bot_test

import (
	. "anubot/bot"
	"time"

	. "github.com/apoydence/eachers"
	"github.com/fluffle/goirc/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ChatFeature", func() {
	var (
		chatFeature       *ChatFeature
		mockFeatureWriter *mockFeatureWriter
		mockDispatcher    *mockDispatcher
	)

	BeforeEach(func() {
		mockFeatureWriter = newMockFeatureWriter()
		mockDispatcher = newMockDispatcher()
		chatFeature = NewChatFeature(mockFeatureWriter, mockDispatcher)
	})

	It("can register it's handler func to listen to chat messages", func() {
		chatFeature.Init()
		Expect(mockFeatureWriter.HandleFuncInput.User).To(Receive(Equal("bot")))
		Expect(mockFeatureWriter.HandleFuncInput.Command).To(Receive(Equal("PRIVMSG")))
		Expect(mockFeatureWriter.HandleFuncInput.User).To(Receive(Equal("streamer")))
		Expect(mockFeatureWriter.HandleFuncInput.Command).To(Receive(Equal("PRIVMSG")))
	})

	It("can write messages to the IRC server as the bot", func() {
		mockFeatureWriter.ChannelOutput.Channel <- "test-chan"

		chatFeature.Send("bot", "test-bot-message")
		Expect(mockFeatureWriter.PrivmsgInput).To(BeCalled(
			With("bot", "test-chan", "test-bot-message"),
		))
	})

	It("can write messages to the IRC server as the streamer", func() {
		mockFeatureWriter.ChannelOutput.Channel <- "test-chan"

		chatFeature.Send("streamer", "test-streamer-message")
		Expect(mockFeatureWriter.PrivmsgInput).To(BeCalled(
			With("streamer", "test-chan", "test-streamer-message"),
		))
	})

	Describe("ChatHandler", func() {
		var handler func(*client.Conn, *client.Line)

		BeforeEach(func() {
			mockFeatureWriter.ChannelOutput.Channel <- "test-user"
		})

		Context("streamer handler", func() {
			BeforeEach(func() {
				handler = chatFeature.ChatHandler("streamer")
			})

			It("dispatches private messages", func() {
				now := time.Now()
				handler(nil, &client.Line{
					Time: now,
					Args: []string{
						"test-target",
						"test-message",
					},
				})
				Expect(mockDispatcher.DispatchInput).To(BeCalled(
					With(Message{
						Target: "test-target",
						Body:   "test-message",
						Time:   now,
					}),
				))
			})

			It("doesn't dispatch messages sent to the current channel", func() {
				now := time.Now()
				handler(nil, &client.Line{
					Time: now,
					Args: []string{
						"test-user",
						"test-message",
					},
				})
				Expect(mockDispatcher.DispatchInput).ToNot(BeCalled())
			})
		})

		Context("bot handler", func() {
			BeforeEach(func() {
				handler = chatFeature.ChatHandler("bot")
			})

			It("dispatches private messages", func() {
				now := time.Now()
				handler(nil, &client.Line{
					Time: now,
					Args: []string{
						"test-target",
						"test-message",
					},
				})
				Expect(mockDispatcher.DispatchInput).To(BeCalled(
					With(Message{
						Target: "test-target",
						Body:   "test-message",
						Time:   now,
					}),
				))
			})

			It("dispatches messages sent to the current channel", func() {
				now := time.Now()
				handler(nil, &client.Line{
					Time: now,
					Args: []string{
						"test-user",
						"test-message",
					},
				})
				Expect(mockDispatcher.DispatchInput).To(BeCalled(
					With(Message{
						Target: "test-user",
						Body:   "test-message",
						Time:   now,
					}),
				))
			})
		})

	})
})
