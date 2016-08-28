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
)
