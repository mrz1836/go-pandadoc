package pandadoc

// OAuthTokenRequest is used for create/refresh token operations.
type OAuthTokenRequest struct {
	GrantType    string
	ClientID     string
	ClientSecret string //nolint:gosec // G117: intentional OAuth credential field
	Code         string
	RefreshToken string //nolint:gosec // G117: intentional OAuth credential field
	Scope        string
	RedirectURI  string
}

// OAuthTokenResponse is the OAuth access-token response.
type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`            //nolint:gosec // G117: intentional OAuth credential field
	RefreshToken string `json:"refresh_token,omitempty"` //nolint:gosec // G117: intentional OAuth credential field
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
}
