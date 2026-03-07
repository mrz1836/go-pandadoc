package pandadoc

// OAuthTokenRequest is used for create/refresh token operations.
type OAuthTokenRequest struct {
	GrantType    string
	ClientID     string
	ClientSecret string
	Code         string
	RefreshToken string
	Scope        string
	RedirectURI  string
}

// OAuthTokenResponse is the OAuth access-token response.
type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
}
