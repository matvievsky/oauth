package token

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const redirectTo = "redirectTo"

func (c *Client) getRedirect(body io.ReadCloser) (*http.Response, error) {
	var result map[string]any

	err := json.NewDecoder(body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding response body: %w", err)
	}

	value, ok := result[redirectTo]
	if !ok {
		return nil, fmt.Errorf("%s not found in response", redirectTo)
	}

	redirectTo, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%s not found in response", redirectTo)
	}

	resp, err := c.Client.Get(redirectTo)
	if err != nil {
		return nil, fmt.Errorf("error following redirect: %w", err)
	}
	defer resp.Body.Close()

	return resp, nil
}
