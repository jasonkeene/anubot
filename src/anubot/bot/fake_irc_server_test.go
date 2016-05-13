package anubot_test

import (
	"errors"
	"net"
	"strings"
	"sync"

	. "github.com/onsi/gomega"
)

type fakeIRCServer struct {
	listener    net.Listener
	connections []net.Conn
	received    []byte
	sent        []byte
	responses   [][]byte

	stop    chan struct{}
	stopped chan struct{}
	connWg  sync.WaitGroup
	mu      sync.Mutex
}

func newFakeIRCServer(listener net.Listener) *fakeIRCServer {
	return &fakeIRCServer{
		listener: listener,

		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

func (a *fakeIRCServer) Start() {
	defer close(a.stopped)
	for {
		select {
		case <-a.stop:
			return
		default:
		}
		conn, _ := a.listener.Accept()
		select {
		case <-a.stop:
			return
		default:
		}
		a.mu.Lock()
		a.connections = append(a.connections, conn)
		a.connWg.Add(1)
		a.mu.Unlock()
		go a.handleConn(conn)
	}
}

func (a *fakeIRCServer) handleConn(conn net.Conn) {
	defer a.connWg.Done()
	for {
		select {
		case <-a.stop:
			return
		default:
		}
		data := make([]byte, 256)
		i, _ := conn.Read(data)
		a.mu.Lock()
		a.received = append(a.received, data[:i]...)
		if strings.HasPrefix(string(data[:i]), "QUIT") {
			conn.Close()
			a.removeConnection(conn)
			a.mu.Unlock()
			return
		}
		if len(a.responses) > 0 {
			data, a.responses = a.responses[0], a.responses[1:]
			conn.Write(data)
			a.sent = append(a.sent, data...)
		}
		a.mu.Unlock()
	}
}

func (a *fakeIRCServer) Stop() {
	close(a.stop)
	err := a.listener.Close()
	Expect(err).ToNot(HaveOccurred())
	println("entering lock")
	a.mu.Lock()
	println("lock entered")
	connections := a.connections
	a.mu.Unlock()
	println("exited lock")
	for _, conn := range connections {
		err = conn.Close()
		Expect(err).ToNot(HaveOccurred())
	}
	<-a.stopped
	a.connWg.Wait()
}

func (a *fakeIRCServer) Respond(data ...string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	message := []byte(strings.Join(data, "\r\n") + "\r\n")
	a.responses = append(a.responses, message)
}

func (a *fakeIRCServer) Connections() []net.Conn {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.connections
}

func (a *fakeIRCServer) Received() []byte {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.received
}

func (a *fakeIRCServer) Sent() []byte {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.sent
}

func (a *fakeIRCServer) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.received = nil
	a.sent = nil
}

func (a *fakeIRCServer) removeConnection(conn net.Conn) {
	i, err := a.findConnectionIndex(conn)
	if err == nil {
		a.connections = append(a.connections[:i], a.connections[i+1:]...)
	}
}

func (a *fakeIRCServer) findConnectionIndex(conn net.Conn) (int, error) {
	var (
		i int
		c net.Conn
	)
	for i, c = range a.connections {
		if c == conn {
			return i, nil
		}
	}
	return 0, errors.New("connection missing")
}
