package api

import (
	"anubot/store"
	"anubot/stream"
	"anubot/twitch/oauth"
	"log"

	"golang.org/x/net/websocket"
)

const (
	twitchHost = "irc.chat.twitch.tv"
	twitchPort = 443
)

// twitchOauthStartHandler responds with a URL to start the Twitch oauth flow.
// The streamer user is required to be the first to begin the oauth flow,
// followed by the bot user.
func twitchOauthStartHandler(e event, s *session) {
	tus, ok := e.Payload.(string)
	if !ok {
		s.Send(event{
			Cmd:       "twitch-oauth-start",
			RequestID: e.RequestID,
			Error:     invalidPayload,
		})
		return
	}

	userID, _ := s.Authenticated()

	var tu store.TwitchUser
	switch tus {
	case "streamer":
		tu = store.Streamer
	case "bot":
		tu = store.Bot
		if !s.Store().TwitchStreamerAuthenticated(userID) {
			s.Send(event{
				Cmd:       "twitch-oauth-start",
				RequestID: e.RequestID,
				Error:     twitchOauthStartOrderError,
			})
			return
		}
	default:
		s.Send(event{
			Cmd:       "twitch-oauth-start",
			RequestID: e.RequestID,
			Error:     invalidPayload,
		})
		return
	}

	url, err := oauth.URL(s.TwitchOauthClientID(), userID, tu, s.Store())
	if err != nil {
		log.Printf("got an err trying to create oauth url: %s", err)
		s.Send(event{
			Cmd:       "twitch-oauth-start",
			RequestID: e.RequestID,
			Error:     unknownError,
		})
		return
	}

	s.Send(event{
		Cmd:       "twitch-oauth-start",
		RequestID: e.RequestID,
		Payload:   url,
	})
}

// twitchClearAuth clears all auth data for the user.
func twitchClearAuth(e event, s *session) {
	s.Store().TwitchClearAuth(s.userID)
	s.Send(event{
		Cmd:       "twitch-clear-auth",
		RequestID: e.RequestID,
	})
}

// twitchUserDetailsHandler provides information on the Twitch streamer and
// bot users.
func twitchUserDetailsHandler(e event, s *session) {
	p := map[string]interface{}{
		"streamer_authenticated": false,
		"streamer_username":      "",
		"streamer_status":        "",
		"streamer_game":          "",

		"bot_authenticated": false,
		"bot_username":      "",
	}
	resp := &event{
		Cmd:     "twitch-user-details",
		Payload: p,
	}
	defer websocket.JSON.Send(s.ws, resp)

	streamerAuthenticated := s.Store().TwitchStreamerAuthenticated(s.userID)
	if !streamerAuthenticated {
		return
	}
	streamerUsername, _ := s.Store().TwitchStreamerCredentials(s.userID)
	status, game, err := s.api.twitch.StreamInfo(streamerUsername)
	if err != nil {
		log.Printf("unable to fetch stream info for user %s: %s",
			streamerUsername, err)
		return
	}
	p["streamer_authenticated"] = streamerAuthenticated
	p["streamer_username"] = streamerUsername
	p["streamer_status"] = status
	p["streamer_game"] = game

	botAuthenticated := s.Store().TwitchBotAuthenticated(s.userID)
	if !botAuthenticated {
		return
	}

	p["bot_authenticated"] = botAuthenticated
	p["bot_username"], _ = s.Store().TwitchBotCredentials(s.userID)
}

// twitchStreamMessages writes chat messages to websocket connection.
func twitchStreamMessages(e event, s *session) {
	streamerUsername, streamerPassword := s.Store().TwitchStreamerCredentials(s.userID)
	s.api.sm.ConnectTwitch(streamerUsername, "oauth:"+streamerPassword, "#"+streamerUsername)

	botUsername, botPassword := s.Store().TwitchStreamerCredentials(s.userID)
	s.api.sm.ConnectTwitch(botUsername, "oauth:"+botPassword, "#"+streamerUsername)

	mw, err := newMessageWriter(
		streamerUsername,
		"twitch:"+streamerUsername,
		"twitch:"+botUsername,
		s.api.pubEndpoints,
		s.ws,
	)
	if err != nil {
		log.Printf("unable to stream messages: %s", err)
		return
	}
	go mw.startStreamer()
	go mw.startBot()
}

// twitchSendMessageHandler accepts messages to send via Twitch chat.
func twitchSendMessageHandler(e event, s *session) {
	data, ok := e.Payload.(map[string]interface{})
	if !ok {
		s.Send(event{
			Cmd:   e.Cmd,
			Error: invalidPayload,
		})
		return
	}
	userType, ok := data["user_type"].(string)
	if !ok {
		s.Send(event{
			Cmd:   e.Cmd,
			Error: invalidPayload,
		})
		return
	}
	message, ok := data["message"].(string)
	if !ok {
		s.Send(event{
			Cmd:   e.Cmd,
			Error: invalidPayload,
		})
		return
	}

	streamerUsername, streamerPassword := s.Store().TwitchStreamerCredentials(s.userID)
	var username, password string
	switch userType {
	case "streamer":
		username, password = streamerUsername, streamerPassword
	case "bot":
		username, password = s.Store().TwitchBotCredentials(s.userID)
	default:
		s.Send(event{
			Cmd:   e.Cmd,
			Error: invalidTwitchUserType,
		})
		return
	}
	s.api.sm.ConnectTwitch(username, "oauth:"+password, "#"+username)
	s.api.sm.Send(stream.TXMessage{
		Type: stream.Twitch,
		Twitch: &stream.TXTwitch{
			Username: username,
			To:       "#" + streamerUsername,
			Message:  message,
		},
	})
}

// twitchUpdateChatDescriptionHandler updates the chat description for Twitch.
func twitchUpdateChatDescriptionHandler(e event, s *session) {
	payload, ok := e.Payload.(map[string]interface{})
	if !ok {
		s.Send(event{
			Cmd:   e.Cmd,
			Error: invalidPayload,
		})
		return
	}
	status, ok := payload["status"].(string)
	if !ok {
		s.Send(event{
			Cmd:   e.Cmd,
			Error: invalidPayload,
		})
		return
	}
	game, ok := payload["game"].(string)
	if !ok {
		s.Send(event{
			Cmd:   e.Cmd,
			Error: invalidPayload,
		})
		return
	}

	user, pass := s.Store().TwitchStreamerCredentials(s.userID)
	err := s.api.twitch.UpdateDescription(status, game, user, pass)
	if err != nil {
		log.Println("unable to update chat description, got error:", err)
		s.Send(event{
			Cmd:   e.Cmd,
			Error: unknownError,
		})
	}
}

// twitchAuthenticateWrapper wraps a handler and makes sure the user attached
// to the session is properly authenticated with twitch.
func twitchAuthenticateWrapper(f handlerFunc) handlerFunc {
	return func(e event, s *session) {
		userID, _ := s.Authenticated()
		if s.Store().TwitchAuthenticated(userID) {
			f(e, s)
			return
		}
		s.Send(event{
			Cmd:       e.Cmd,
			RequestID: e.RequestID,
			Error:     twitchAuthenticationError,
		})
	}
}
