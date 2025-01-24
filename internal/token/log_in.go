package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) LogIn(loginURI, userLogin, userPassword, loginChallenge string) (resp *http.Response, err error) {
	jsonData, err := json.Marshal(map[string]any{
		"email":                   userLogin,
		"password":                userPassword,
		"remember":                false,
		"challenge":               loginChallenge,
		"subscribedToNewsletters": false,
		"dataProcessingAgreement": false,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshalling data: %v", err)
	}

	return http.Post(loginURI, "application/json", bytes.NewBuffer(jsonData))
}
