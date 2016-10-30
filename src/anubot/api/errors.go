package api

type apiError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func (e apiError) Error() string {
	return e.Text
}

var (
	unknownError = &apiError{
		Code: 0,
		Text: "unknown error",
	}
	invalidCommand = &apiError{
		Code: 1,
		Text: "invalid command",
	}
	invalidPayload = &apiError{
		Code: 2,
		Text: "invalid payload data",
	}
	usernameTaken = &apiError{
		Code: 3,
		Text: "username has already been taken",
	}
	authenticationError = &apiError{
		Code: 4,
		Text: "authentication error",
	}
	twitchAuthenticationError = &apiError{
		Code: 5,
		Text: "authentication error with twitch",
	}
	twitchOauthStartOrderError = &apiError{
		Code: 6,
		Text: "unable to start oauth flow for bot, streamer not finished",
	}
	invalidTwitchUserType = &apiError{
		Code: 7,
		Text: "you specified an invalid user type",
	}
	bttvUnavailable = &apiError{
		Code: 8,
		Text: "unable to gather emoji from bttv api",
	}
)
