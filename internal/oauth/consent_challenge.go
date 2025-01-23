package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) ConsentChallenge(consentURI, consentChallenge string) (*http.Response, error) {
	jsonData, err := json.Marshal(map[string]any{
		"challenge": consentChallenge,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshalling data: %v", err)
	}

	return http.Post(consentURI, "application/json", bytes.NewBuffer(jsonData))
}
