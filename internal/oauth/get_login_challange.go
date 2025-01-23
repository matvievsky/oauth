package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
)

func (c *Client) GetLoginChallenge(oidHost, hydraClientID, redirectURI string) (*http.Response, error) {
	rb := make([]byte, 32)
	_, err := rand.Read(rb)
	if err != nil {
		return nil, err
	}

	data := url.Values{}
	data.Add("client_id", hydraClientID)
	data.Add("response_type", "code")
	data.Add("scope", "offline_access offline openid")
	data.Add("redirect_uri", redirectURI)
	data.Add("state", base64.StdEncoding.EncodeToString(rb))
	data.Add("prompt", "login")

	u := url.URL{
		Scheme: "https",
		Host:   oidHost,
		Path:   "oauth2/auth",
	}

	u.RawQuery = data.Encode()

	return c.Get(u.String())
}
