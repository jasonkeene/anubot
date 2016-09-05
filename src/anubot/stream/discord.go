package stream

type discordConn struct {
	d Dispatcher
}

func (c *discordConn) send(m TXMessage) {

}

func (c *discordConn) close() error {
	return nil
}

func connectDiscord(u, p, c string, d Dispatcher) (*discordConn, error) {
	return nil, nil
}
