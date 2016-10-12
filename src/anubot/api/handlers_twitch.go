package api

import (
	"anubot/store"
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

	s.Send(event{
		Cmd:       "twitch-oauth-start",
		RequestID: e.RequestID,
		Payload:   oauth.URL(s.TwitchOauthClientID(), userID, tu, s.Store()),
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
	username, password := s.Store().TwitchStreamerCredentials(s.userID)
	s.api.sm.ConnectTwitch(username, "oauth:"+password, "#"+username)
	mw, err := newMessageWriter("twitch:"+username, s.api.pubEndpoints, s.ws)
	if err != nil {
		log.Printf("unable to stream messages: %s", err)
		return
	}
	go mw.start()
}

// twitchSendMessageHandler accepts messages to send via Twitch chat.
func twitchSendMessageHandler(e event, s *session) {
	//data, ok := e.Payload.(map[string]interface{})
	//if !ok {
	//	return
	//}
	//user, ok := data["user"].(string)
	//if !ok {
	//	return
	//}
	//message, ok := data["message"].(string)
	//if !ok {
	//	return
	//}
	//s.bot.ChatFeature().Send(user, message)
}

// twitchUpdateChatDescriptionHandler updates the chat description for Twitch.
func twitchUpdateChatDescriptionHandler(e event, s *session) {
	//payload := e.Payload.(map[string]interface{})
	//game := payload["game"].(string)
	//status := payload["status"].(string)
	//user, pass, err := s.store.Credentials("user")
	//if err != nil {
	//	fmt.Println("bad creds!")
	//	return
	//}
	//err = updateDescription(game, status, user, pass)
	//if err != nil {
	//	fmt.Println(err)
	//}
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
