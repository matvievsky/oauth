package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) consentChallenge(consentURI, consentChallenge, scope string) (*http.Response, error) {
	jsonData, err := json.Marshal(map[string]any{
		"challenge": consentChallenge,
		"scope":     scope,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshalling data: %v", err)
	}

	return http.Post(consentURI, "application/json", bytes.NewBuffer(jsonData))
}
