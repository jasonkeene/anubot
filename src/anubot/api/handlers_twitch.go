package api

import (
	"anubot/store"
	"anubot/twitch/oauth"
)

const (
	twitchHost = "irc.chat.twitch.tv"
	twitchPort = 443
)

// twitchOauthStartHandler responds with a URL to start the Twitch oauth flow.
func twitchOauthStartHandler(e event, s *session) {
	// grab the twitch user string from payload
	tus, ok := e.Payload.(string)
	if !ok {
		s.Send(event{
			Cmd:       "twitch-oauth-start",
			RequestID: e.RequestID,
			Error:     invalidPayload,
		})
		return
	}

	// validate twitch user
	var tu store.TwitchUser
	switch tus {
	case "streamer":
		tu = store.Streamer
	case "bot":
		tu = store.Bot
	default:
		s.Send(event{
			Cmd:       "twitch-oauth-start",
			RequestID: e.RequestID,
			Error:     invalidPayload,
		})
		return
	}

	userID, _ := s.Authenticated()
	s.Send(event{
		Cmd:       "twitch-oauth-start",
		RequestID: e.RequestID,
		Payload:   oauth.URL(s.TwitchOauthClientID(), userID, tu, s.Store()),
	})
}

// twitchUserDetailsHandler provides information on the Twitch streamer and
// bot users.
func twitchUserDetailsHandler(e event, s *session) {
	//resp := &event{
	//	Cmd: "twitch-user-details",
	//	Payload: map[string]interface{}{
	//		"authenticated": false,
	//		"streamer":      "",
	//		"bot":           "",
	//		"status":        "",
	//		"game":          "",
	//	},
	//}
	//defer websocket.JSON.Send(s.ws, resp)

	//streamerUser, streamerPass, err := s.store.Credentials("user")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//botUser, botPass, err := s.store.Credentials("bot")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//if streamerUser == "" || streamerPass == "" ||
	//	botUser == "" || botPass == "" {
	//	println("bad creds")
	//	return
	//}

	//status, game, _ := fetchStreamInfo(streamerUser)
	//resp.Payload = map[string]interface{}{
	//	"authenticated": true,
	//	"streamer":      streamerUser,
	//	"bot":           botUser,
	//	"status":        status,
	//	"game":          game,
	//}
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
