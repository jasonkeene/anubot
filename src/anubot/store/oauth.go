package store

// OauthData contains the data returned from Twitch when finishing the Oauth
// flow.
type OauthData struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
}
