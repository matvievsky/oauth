package oauth

import (
	"net/http"
	"net/http/cookiejar"
)

type Client struct {
	*http.Client
}

func NewClient() *Client {
	jar, _ := cookiejar.New(nil)

	return &Client{&http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}}
}
