package api

// eventHandlers is a map of event names to their Handlers.
var eventHandlers map[string]eventHandler

// An eventHandler is a func that can handle events from a websocket
// connection.
type eventHandler interface {
	HandleEvent(event, *session)
}

// handlerFunc is a wrapper type that allows event Handlers to be defined as a
// func.
type handlerFunc func(event, *session)

// HandleEvent simply dispatches the event to the underlying func.
func (f handlerFunc) HandleEvent(e event, s *session) {
	f(e, s)
}

func init() {
	eventHandlers = make(map[string]eventHandler)

	// public
	{
		// general
		eventHandlers["ping"] = handlerFunc(pingHandler)

		// authentication
		eventHandlers["register"] = handlerFunc(registerHandler)
		eventHandlers["authenticate"] = handlerFunc(authenticateHandler)
	}

	// authenticated
	{
		// twitch auth
		eventHandlers["twitch-oauth-start"] = authenticateWrapper(handlerFunc(twitchOauthStartHandler))
		eventHandlers["twtitch-user-details"] = authenticateWrapper(handlerFunc(twitchUserDetailsHandler))

		// twitch chat
		eventHandlers["twitch-send-message"] = authenticateWrapper(handlerFunc(twitchSendMessageHandler))
		eventHandlers["twitch-update-chat-description"] = authenticateWrapper(handlerFunc(twitchUpdateChatDescriptionHandler))
	}
}
