package api

// pingHandler responds with pong data frame used for testing connectivity.
func pingHandler(e event, s *session) {
	s.Send(event{
		Cmd:       "pong",
		RequestID: e.RequestID,
	})
}

// methodsHandler responds with a list of methods the API supports.
func methodsHandler(e event, s *session) {
	methods := []string{}
	for m := range eventHandlers {
		methods = append(methods, m)
	}
	s.Send(event{
		Cmd:       "methods",
		RequestID: e.RequestID,
		Payload:   methods,
	})
}
