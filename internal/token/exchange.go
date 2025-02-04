package token

import (
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) exchange(query map[string]string) (resp *http.Response, err error) {
	data := url.Values{}
	data.Add("client_id", c.hydraClientID)
	data.Add("redirect_uri", c.redirect.String())

	for key, val := range query {
		data.Add(key, val)
	}

	u := url.URL{
		Scheme: c.oid.Scheme,
		Host:   c.oid.Host,
		Path:   "oauth2/token",
	}

	return c.Post(u.String(), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}
