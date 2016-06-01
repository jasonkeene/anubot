package bot_test

import (
	. "anubot/bot"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageDispatcher", func() {
	var (
		dispatcher *MessageDispatcher
		testMsg    Message
	)

	BeforeEach(func() {
		dispatcher = NewMessageDispatcher()
		testMsg = Message{
			Channel: "#test-chan",
			Body:    "hello world",
		}
	})

	It("accepts messages", func() {
		dispatcher.Dispatch(testMsg)
		msgs := dispatcher.Messages("#test-chan")
		Expect(msgs).To(Receive(Equal(testMsg)))
	})

	It("sends new messages to readers as they come in", func() {
		msgs := dispatcher.Messages("#test-chan")
		dispatcher.Dispatch(testMsg)
		Expect(msgs).To(Receive(Equal(testMsg)))
	})

	It("removes channels that are no longer needed", func() {
		msgs := dispatcher.Messages("#test-chan")
		dispatcher.Remove(msgs)
		Expect(msgs).To(BeClosed())
		dispatcher.Dispatch(testMsg) // this should not panic
	})
})
