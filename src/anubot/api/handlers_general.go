package api

// pingHandler responds with pong data frame used for testing connectivity.
func pingHandler(e event, s *session) {
	s.Send(event{
		Cmd:       "pong",
		RequestID: e.RequestID,
	})
}
