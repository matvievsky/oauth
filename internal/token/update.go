package token

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) Update(refreshToken string) error {
	exchangeResp, err := c.exchange(map[string]string{
		"grant_type":   "refresh_token",
		"refresh_toke": refreshToken,
	})
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer exchangeResp.Body.Close()

	if exchangeResp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code: %v", exchangeResp.StatusCode)
	}

	result := map[string]any{}

	err = json.NewDecoder(exchangeResp.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("error decoding response body: %v", err)
	}

	// TODO: Add os.SetEnv to the updated access token
	fmt.Printf("%q", result)

	return nil
}
