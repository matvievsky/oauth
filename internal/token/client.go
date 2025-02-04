package token

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Client struct {
	*http.Client
	oid, redirect url.URL
	hydraClientID string
}

func NewClient(oidHost, redirectURI, hydraClientID string) *Client {
	jar, _ := cookiejar.New(nil)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}

	oidURI, err := url.ParseRequestURI(oidHost)
	if err != nil {
		panic(err)
	}

	redirectURL, err := url.Parse(redirectURI)
	if err != nil {
		panic(err)
	}

	return &Client{client, *oidURI, *redirectURL, hydraClientID}
}
