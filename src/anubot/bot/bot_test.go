package bot_test

import (
	"anubot/bot"
	"anubot/stream"
	"encoding/json"
	"testing"
	"time"

	"github.com/fluffle/goirc/client"
	"github.com/pebbe/zmq4"
	"github.com/stretchr/testify/require"
)

func TestBotDispatchesMessagesToFeatures(t *testing.T) {
	require := require.New(t)

	pubTopic, subTopic := "test-topic", "test-topic"

	f := newMockFeature()
	pub := setupPubSocket(require)
	defer pub.Close()
	expected, toSend := testMessage(require, "test-message")

	b, err := bot.New(subTopic, []string{"inproc://test-pub"})
	require.Nil(err)
	b.SetFeature("test-feature", f)
	go b.Start()
	defer b.Stop()

	_, err = pub.SendMessage(pubTopic, toSend)
	require.Nil(err)

	select {
	case actual := <-f.HandleMessageInput.Ms:
		require.Equal(expected.Type, actual.Type)
		require.Equal(expected.Twitch.Line.Raw, actual.Twitch.Line.Raw)
	case <-time.After(3 * time.Second):
		require.Fail("timed out waiting for bot to dispatch message")
	}
}

func TestBotDoesNotDispatchMessagesIfTopicDoesNotMatch(t *testing.T) {
	require := require.New(t)

	pubTopic, subTopic := "test-a", "test-b"

	f := newMockFeature()
	pub := setupPubSocket(require)
	defer pub.Close()
	_, badBytes := testMessage(require, "test-message")
	expected, finalBytes := testMessage(require, "final-message")

	b, err := bot.New(subTopic, []string{"inproc://test-pub"})
	require.Nil(err)
	b.SetFeature("test-feature", f)
	go b.Start()
	defer b.Stop()

	_, err = pub.SendMessage(pubTopic, badBytes)
	require.Nil(err)
	_, err = pub.SendMessage(subTopic, finalBytes)
	require.Nil(err)

	select {
	case actual := <-f.HandleMessageInput.Ms:
		require.Equal(expected.Type, actual.Type)
		require.Equal(expected.Twitch.Line.Raw, actual.Twitch.Line.Raw)
		//require.Equal(expected, actual)
	case <-time.After(3 * time.Second):
		require.Fail("timed out waiting for bot to dispatch message")
	}
}

func setupPubSocket(require *require.Assertions) *zmq4.Socket {
	pub, err := zmq4.NewSocket(zmq4.PUB)
	require.Nil(err)
	require.Nil(pub.Bind("inproc://test-pub"))
	return pub
}

func testMessage(require *require.Assertions, raw string) (stream.RXMessage, []byte) {
	msg := stream.RXMessage{
		Type: stream.Twitch,
		Twitch: &stream.RXTwitch{
			Line: &client.Line{
				Raw: raw,
			},
		},
	}
	msgBytes, err := json.Marshal(msg)
	require.Nil(err)
	return msg, msgBytes
}
