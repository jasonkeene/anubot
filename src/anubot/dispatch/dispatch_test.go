package dispatch_test

import (
	"crypto/rand"
	"fmt"
	"log"
	"testing"

	"github.com/a8m/expect"
	"github.com/pebbe/zmq4"

	"anubot/dispatch"
)

func TestDispatcher(t *testing.T) {
	expect := expect.New(t)

	ep := endpoints()
	dispatch.Start(
		dispatch.WithPullEndpoints([]string{ep["pull"]}),
		dispatch.WithPubEndpoints([]string{ep["pub"]}),
		dispatch.WithPushEndpoints([]string{ep["push"]}),
	)

	push := setupPushSocket(expect, ep["pull"])
	sub := setupSubSocket(expect, "test-topic", ep["pub"])
	pull := setupPullSocket(expect, ep["push"])

	_, err := push.SendBytes([]byte("test-topic"), zmq4.SNDMORE)
	expect(err).To.Be.Nil().Else.FailNow()
	_, err = push.SendBytes([]byte("test-content"), 0)
	expect(err).To.Be.Nil().Else.FailNow()

	parts, err := sub.RecvMessageBytes(0)
	expect(err).To.Be.Nil().Else.FailNow()
	expect(len(parts)).To.Equal(2)
	expect(parts[0]).To.Equal([]byte("test-topic"))
	expect(parts[1]).To.Equal([]byte("test-content"))

	pullContent, err := pull.RecvBytes(0)
	expect(err).To.Be.Nil().Else.FailNow()
	expect(pullContent).To.Equal([]byte("test-content"))
}

func endpoints() map[string]string {
	return map[string]string{
		"pull": "inproc://test-dispatch-pull-" + randString(),
		"pub":  "inproc://test-dispatch-pub-" + randString(),
		"push": "inproc://test-dispatch-push-" + randString(),
	}
}

func randString() string {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	if err != nil {
		log.Panicf("unable to read randomness %s:", err)
	}
	return fmt.Sprintf("%x", b)
}

func setupPushSocket(expect expect.Expectation, endpoint string) *zmq4.Socket {
	push, err := zmq4.NewSocket(zmq4.PUSH)
	expect(err).To.Be.Nil().Else.FailNow()
	expect(push.Connect(endpoint)).To.Be.Nil().Else.FailNow()
	return push
}

func setupSubSocket(expect expect.Expectation, topic, endpoint string) *zmq4.Socket {
	sub, err := zmq4.NewSocket(zmq4.SUB)
	expect(err).To.Be.Nil().Else.FailNow()
	err = sub.SetSubscribe(topic)
	expect(err).To.Be.Nil().Else.FailNow()
	expect(sub.Connect(endpoint)).To.Be.Nil().Else.FailNow()
	return sub
}

func setupPullSocket(expect expect.Expectation, endpoint string) *zmq4.Socket {
	pull, err := zmq4.NewSocket(zmq4.PULL)
	expect(err).To.Be.Nil().Else.FailNow()
	expect(pull.Connect(endpoint)).To.Be.Nil().Else.FailNow()
	return pull
}
