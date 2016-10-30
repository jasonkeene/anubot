package api

import (
	"anubot/bttv"
	"log"
)

// bttvEmojiHandler returns emoji from BTTV. If the user has authenticated
// their streamer user with Twitch it will also include channel specific
// emoji.
func bttvEmojiHandler(e event, s *session) {
	var streamerUsername string
	streamerAuthenticated := s.Store().TwitchStreamerAuthenticated(s.userID)
	if streamerAuthenticated {
		streamerUsername, _ = s.Store().TwitchStreamerCredentials(s.userID)
	}
	payload, err := bttv.Emoji(streamerUsername)
	if err != nil {
		log.Printf("error getting bttv emoji: %s", err)
		s.Send(event{
			Cmd:       e.Cmd,
			RequestID: e.RequestID,
			Error:     bttvUnavailable,
		})
		return
	}
	s.Send(event{
		Cmd:       e.Cmd,
		RequestID: e.RequestID,
		Payload:   payload,
	})
}
