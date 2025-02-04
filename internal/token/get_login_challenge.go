package token

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
)

func (c *Client) GetLoginChallenge(scope string) (*http.Response, error) {
	rb := make([]byte, 32)
	_, err := rand.Read(rb)
	if err != nil {
		return nil, err
	}

	data := url.Values{}
	data.Add("client_id", c.hydraClientID)
	data.Add("response_type", "code")
	data.Add("scope", scope)
	data.Add("redirect_uri", c.redirect.String())
	data.Add("state", base64.StdEncoding.EncodeToString(rb))
	data.Add("prompt", "login")

	u := url.URL{
		Scheme: c.oid.Scheme,
		Host:   c.oid.Host,
		Path:   "oauth2/auth",
	}

	u.RawQuery = data.Encode()

	return c.Client.Get(u.String())
}
