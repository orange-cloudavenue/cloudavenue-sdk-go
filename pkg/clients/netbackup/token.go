package clientnetbackup

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
)

type token struct {
	baererToken string
	expiresAt   time.Time

	username string
	password string

	endpoint string

	debug bool
}

// IsExpired - Returns true if the token is expired
func (t *token) IsExpired() bool {
	return t.expiresAt.Before(time.Now())
}

// IsSet - Returns true if the token is set
func (t *token) IsSet() bool {
	return t.baererToken != ""
}

// GetToken - Returns the token
func (t *token) GetToken() string {
	return t.baererToken
}

// RefreshToken - Refreshes the token
func (t *token) RefreshToken() error {
	if !t.IsSet() || t.IsExpired() {
		c := resty.New().SetBaseURL(t.endpoint)

		criteria := url.Values{
			"grant_type": {"password"},
			"username":   {t.username},
			"password":   {t.password},
		}

		r, err := c.R().
			SetDebug(t.debug).
			SetHeader("Content-Type", "application/x-www-form-url-encoded").
			SetHeader("Accept", "application/json").
			SetResult(&authTokenResponse{}).
			SetError(&apiAuthTokenErrorResponse{}).
			SetFormDataFromValues(criteria).
			Post("/auth/token")
		if err != nil {
			return err
		}

		if r.IsError() {
			return fmt.Errorf("authentification failed : ErrorCode:%s - %s ", r.Status(), r.Error().(*apiAuthTokenErrorResponse).FormatError())
		}

		refreshedToken := r.Result().(*authTokenResponse)
		t.baererToken = refreshedToken.AccessToken
		t.expiresAt = time.Now().Add(time.Duration(refreshedToken.ExpiresIn) * time.Second)
	}
	return nil
}

type authTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	UserName     string `json:"userName"`
	UserGUID     string `json:"userGuid"`
	Issued       string `json:".issued"`
	Expires      string `json:".expires"`
}

type apiAuthTokenErrorResponse struct {
	Error string `json:"error"`
}

func (e *apiAuthTokenErrorResponse) FormatError() string {
	return fmt.Sprintf("Error:%s", e.Error)
}
