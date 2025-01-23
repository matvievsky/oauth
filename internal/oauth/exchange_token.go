package oauth

import (
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) ExchangeToken(redirectURI, oidHost, hydraClientID, code string) (resp *http.Response, err error) {
	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("client_id", hydraClientID)
	data.Add("code", code)
	data.Add("redirect_uri", redirectURI)

	u := url.URL{
		Scheme: "https",
		Host:   oidHost,
		Path:   "oauth2/token",
	}

	return c.Post(u.String(), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}
