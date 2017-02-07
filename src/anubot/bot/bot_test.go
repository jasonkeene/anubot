package bot_test

import (
	"anubot/bot"
	"anubot/stream"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/a8m/expect"
	"github.com/fluffle/goirc/client"
	"github.com/pebbe/zmq4"
)

func TestBotDispatchesMessagesToFeatures(t *testing.T) {
	expect := expect.New(t)

	pubTopic, subTopic := "test-topic", "test-topic"

	f := newMockFeature()
	pub, endpoints := setupPubSocket(expect)
	defer func() {
		err := pub.Close()
		if err != nil {
			log.Printf("got err while closing pub socket: %s", err)
		}
	}()
	expected, toSend := testMessage(expect, "test-message")

	b, err := bot.New([]string{subTopic}, bot.WithSubEndpoints(endpoints))
	expect(err).To.Be.Nil()
	b.SetFeature("test-feature", f)
	go b.Start()
	defer b.Stop()

	_, err = pub.SendMessage(pubTopic, toSend)
	expect(err).To.Be.Nil()

	select {
	case actual := <-f.HandleMessageInput.Ms:
		expect(expected.Type).To.Equal(actual.Type)
		expect(expected.Twitch.Line.Raw).To.Equal(actual.Twitch.Line.Raw)
	case <-time.After(3 * time.Second):
		fmt.Println("timed out waiting for bot to dispatch message")
		t.Fail()
	}
}

func TestBotDoesNotDispatchMessagesIfTopicDoesNotMatch(t *testing.T) {
	expect := expect.New(t)

	pubTopic, subTopic := "test-a", "test-b"

	f := newMockFeature()
	pub, endpoints := setupPubSocket(expect)
	defer func() {
		err := pub.Close()
		if err != nil {
			log.Printf("got err while closing pub socket: %s", err)
		}
	}()
	_, badBytes := testMessage(expect, "test-message")
	expected, finalBytes := testMessage(expect, "final-message")

	b, err := bot.New([]string{subTopic}, bot.WithSubEndpoints(endpoints))
	expect(err).To.Be.Nil()
	b.SetFeature("test-feature", f)
	go b.Start()
	defer b.Stop()

	_, err = pub.SendMessage(pubTopic, badBytes)
	expect(err).To.Be.Nil()
	_, err = pub.SendMessage(subTopic, finalBytes)
	expect(err).To.Be.Nil()

	select {
	case actual := <-f.HandleMessageInput.Ms:
		expect(expected.Type).To.Equal(actual.Type)
		expect(expected.Twitch.Line.Raw).To.Equal(actual.Twitch.Line.Raw)
	case <-time.After(3 * time.Second):
		fmt.Println("timed out waiting for bot to dispatch message")
		t.Fail()
	}
}

func setupPubSocket(expect func(v interface{}) *expect.Expect) (*zmq4.Socket, []string) {
	endpoint := "inproc://test-pub-" + randString()
	pub, err := zmq4.NewSocket(zmq4.PUB)
	expect(err).To.Be.Nil()
	expect(pub.Bind(endpoint)).To.Be.Nil()
	return pub, []string{endpoint}
}

func randString() string {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	if err != nil {
		log.Panicf("unable to read randomness %s:", err)
	}
	return fmt.Sprintf("%x", b)
}

func testMessage(expect func(v interface{}) *expect.Expect, raw string) (stream.RXMessage, []byte) {
	msg := stream.RXMessage{
		Type: stream.Twitch,
		Twitch: &stream.RXTwitch{
			Line: &client.Line{
				Raw: raw,
			},
		},
	}
	msgBytes, err := json.Marshal(msg)
	expect(err).To.Be.Nil()
	return msg, msgBytes
}
