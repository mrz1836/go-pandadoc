package pandadoc

import (
	"context"
	"net/http"
	"net/url"
)

// oauthService implements OAuthService.
type oauthService struct {
	client *Client
}

// Token exchanges either an authorization code or refresh token for an access token.
func (s *oauthService) Token(ctx context.Context, req *OAuthTokenRequest) (*OAuthTokenResponse, error) {
	if req == nil {
		return nil, ErrNilRequest
	}

	form := url.Values{}
	if req.GrantType != "" {
		form.Set("grant_type", req.GrantType)
	}
	if req.ClientID != "" {
		form.Set("client_id", req.ClientID)
	}
	if req.ClientSecret != "" {
		form.Set("client_secret", req.ClientSecret)
	}
	if req.Code != "" {
		form.Set("code", req.Code)
	}
	if req.RefreshToken != "" {
		form.Set("refresh_token", req.RefreshToken)
	}
	if req.Scope != "" {
		form.Set("scope", req.Scope)
	}
	if req.RedirectURI != "" {
		form.Set("redirect_uri", req.RedirectURI)
	}

	var out OAuthTokenResponse
	err := s.client.decodeJSON(ctx, &request{
		method:      http.MethodPost,
		path:        "/oauth2/access_token",
		requireAuth: false,
		formBody:    form,
	}, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
