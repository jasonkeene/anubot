package stream

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
	"testing"
)

const (
	localhostCert = `-----BEGIN CERTIFICATE-----
MIICEzCCAXygAwIBAgIQMIMChMLGrR+QvmQvpwAU6zANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMCAXDTcwMDEwMTAwMDAwMFoYDzIwODQwMTI5MTYw
MDAwWjASMRAwDgYDVQQKEwdBY21lIENvMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCB
iQKBgQDuLnQAI3mDgey3VBzWnB2L39JUU4txjeVE6myuDqkM/uGlfjb9SjY1bIw4
iA5sBBZzHi3z0h1YV8QPuxEbi4nW91IJm2gsvvZhIrCHS3l6afab4pZBl2+XsDul
rKBxKKtD1rGxlG4LjncdabFn9gvLZad2bSysqz/qTAUStTvqJQIDAQABo2gwZjAO
BgNVHQ8BAf8EBAMCAqQwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUw
AwEB/zAuBgNVHREEJzAlggtleGFtcGxlLmNvbYcEfwAAAYcQAAAAAAAAAAAAAAAA
AAAAATANBgkqhkiG9w0BAQsFAAOBgQCEcetwO59EWk7WiJsG4x8SY+UIAA+flUI9
tyC4lNhbcF2Idq9greZwbYCqTTTr2XiRNSMLCOjKyI7ukPoPjo16ocHj+P3vZGfs
h1fIw3cSS2OolhloGw/XM6RWPWtPAlGykKLciQrBru5NAPvCMsb/I1DAceTiotQM
fblo6RBxUQ==
-----END CERTIFICATE-----`
	localhostKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDuLnQAI3mDgey3VBzWnB2L39JUU4txjeVE6myuDqkM/uGlfjb9
SjY1bIw4iA5sBBZzHi3z0h1YV8QPuxEbi4nW91IJm2gsvvZhIrCHS3l6afab4pZB
l2+XsDulrKBxKKtD1rGxlG4LjncdabFn9gvLZad2bSysqz/qTAUStTvqJQIDAQAB
AoGAGRzwwir7XvBOAy5tM/uV6e+Zf6anZzus1s1Y1ClbjbE6HXbnWWF/wbZGOpet
3Zm4vD6MXc7jpTLryzTQIvVdfQbRc6+MUVeLKwZatTXtdZrhu+Jk7hx0nTPy8Jcb
uJqFk541aEw+mMogY/xEcfbWd6IOkp+4xqjlFLBEDytgbIECQQDvH/E6nk+hgN4H
qzzVtxxr397vWrjrIgPbJpQvBsafG7b0dA4AFjwVbFLmQcj2PprIMmPcQrooz8vp
jy4SHEg1AkEA/v13/5M47K9vCxmb8QeD/asydfsgS5TeuNi8DoUBEmiSJwma7FXY
fFUtxuvL7XvjwjN5B30pNEbc6Iuyt7y4MQJBAIt21su4b3sjXNueLKH85Q+phy2U
fQtuUE9txblTu14q3N7gHRZB4ZMhFYyDy8CKrN2cPg/Fvyt0Xlp/DoCzjA0CQQDU
y2ptGsuSmgUtWj3NM9xuwYPm+Z/F84K6+ARYiZ6PYj013sovGKUFfYAqVXVlxtIX
qyUBnu3X9ps8ZfjLZO7BAkEAlT4R5Yl6cGhaJQYZHOde3JEMhNRcVFMO8dJDaFeo
f9Oeos0UUothgiDktdQHxdNEwLjQf7lJJBzV+5OtwswCWA==
-----END RSA PRIVATE KEY-----`
)

var serverTLSConfig *tls.Config

func init() {
	cert, err := tls.X509KeyPair([]byte(localhostCert), []byte(localhostKey))
	if err != nil {
		panic(fmt.Sprintf("unable to create cert for testing: %v", err))
	}
	serverTLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
}

type fakeIRCServer struct {
	ln net.Listener
}

func newFakeIRCServer(t *testing.T) *fakeIRCServer {
	ln, err := tls.Listen("tcp", ":0", serverTLSConfig)
	if err != nil {
		t.Fatal("unable to listen")
	}
	return &fakeIRCServer{
		ln: ln,
	}
}

func (s *fakeIRCServer) port() int {
	return s.ln.Addr().(*net.TCPAddr).Port
}

func (s *fakeIRCServer) accept() *ircConn {
	conn, err := s.ln.Accept()
	if err != nil {
		log.Panicf("fakeIRCServer.accept: got error wile accepting: %s", err)
	}
	log.Print("fakeIRCServer.accept: accepted connection")
	return &ircConn{
		conn: conn,
	}
}

func (s *fakeIRCServer) close() {
	err := s.ln.Close()
	if err != nil {
		log.Panic("fakeIRCServer.close: unable to close listener")
	}
}

type ircConn struct {
	conn net.Conn
}

func (c *ircConn) receive(cmd string) string {
	for {
		rxb := make([]byte, 2048)
		log.Print("ircConn.receive: waiting on read")
		n, err := c.conn.Read(rxb)
		if err != nil {
			log.Panicf("ircConn.receive: got error while reading from connection: %s", err)
		}
		result := strings.TrimRight(string(rxb[:n]), "\r\n")
		if strings.HasPrefix(result, cmd) {
			log.Printf(`ircConn.receive: read: "%s"`, result)
			return result
		}
	}
}

func (c *ircConn) send(line string) {
	txb := []byte(line + "\r\n")
	log.Print("ircConn.send: waiting on write")
	n, err := c.conn.Write(txb)
	if err != nil {
		log.Panicf("ircConn.send: got error while writing to connection: %s", err)
	}
	if n != len(txb) {
		log.Panic("ircConn.send: did not write the entire message to connection")
	}
	log.Printf(`ircConn.send: wrote: "%s"`, line)
}

func (c *ircConn) close() {
	err := c.conn.Close()
	if err != nil {
		log.Panic("ircConn.close: unable to close connection")
	}
}
