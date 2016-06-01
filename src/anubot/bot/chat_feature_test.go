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
		It("dispatches messages it got from the IRC server", func() {
			now := time.Now()
			chatFeature.ChatHandler(nil, &client.Line{
				Time: now,
				Args: []string{
					"#test-chan",
					"test-message",
				},
			})
			Expect(mockDispatcher.DispatchInput).To(BeCalled(
				With(Message{
					Channel: "#test-chan",
					Body:    "test-message",
					Time:    now,
				}),
			))
		})
	})
})
