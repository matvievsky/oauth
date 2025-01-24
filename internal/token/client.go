package token

import (
	"net/http"
	"net/http/cookiejar"
)

type Client struct {
	*http.Client
	oidHost, redirectURI, hydraClientID string
}

func NewClient(oidHost, redirectURI, hydraClientID string) *Client {
	jar, _ := cookiejar.New(nil)

	client := &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}

	return &Client{client, oidHost, redirectURI, hydraClientID}
}
