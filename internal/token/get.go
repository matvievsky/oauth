package token

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/matvievsky/oauth/internal/oauth"

	"github.com/google/uuid"
)

const location = "Location"

func Get(oidHost, hydraClientID, redirectURI, loginURI, userLogin, userPassword, consentURI string) error {
	_, err := uuid.Parse(hydraClientID)
	if err != nil {
		return fmt.Errorf("client ID is broken: %w", err)
	}

	client := oauth.NewClient()

	challengeResp, err := client.GetLoginChallenge(oidHost, hydraClientID, redirectURI)
	if err != nil {
		return err
	}
	defer challengeResp.Body.Close()

	if !(challengeResp.StatusCode >= 300 && challengeResp.StatusCode < 400) {
		return fmt.Errorf("response status code: %d", challengeResp.StatusCode)
	}

	parsedURL, err := url.Parse(challengeResp.Header.Get(location))
	if err != nil {
		return fmt.Errorf("error parsing redirect URL: %v", err)
	}

	loginChallenge := parsedURL.Query().Get("login_challenge")
	if loginChallenge == "" {
		return fmt.Errorf("login challenge not found in redirect URL")
	}

	logInResp, err := client.LogIn(loginURI, userLogin, userPassword, loginChallenge)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer logInResp.Body.Close()

	if logInResp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code: %v", logInResp.StatusCode)
	}

	redirect, err := client.GetRedirect(logInResp.Body)
	if err != nil {
		return fmt.Errorf("can't get redirect: %w", err)
	}

	parsedURL, err = url.Parse(redirect.Header.Get(location))
	if err != nil {
		return fmt.Errorf("error parsing redirect URL: %v", err)
	}

	consentChallenge := parsedURL.Query().Get("consent_challenge")
	if consentChallenge == "" {
		return fmt.Errorf("consent challenge not found in redirect URL")
	}

	consentChallengeResp, err := client.ConsentChallenge(consentURI, consentChallenge)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer consentChallengeResp.Body.Close()

	if consentChallengeResp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code: %v", consentChallengeResp.StatusCode)
	}

	redirect, err = client.GetRedirect(consentChallengeResp.Body)
	if err != nil {
		return fmt.Errorf("can't get redirect: %w", err)
	}

	parsedURL, err = url.Parse(redirect.Header.Get(location))
	if err != nil {
		return fmt.Errorf("error parsing redirect URL: %v", err)
	}

	code := parsedURL.Query().Get("code")
	if code == "" {
		return fmt.Errorf("code not found in redirect URL")
	}

	exchangeResp, err := client.ExchangeToken(redirectURI, oidHost, hydraClientID, code)
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

	accessToken := fmt.Sprintf("%q", result["access_token"].(string))

	fmt.Println(accessToken[1 : len(accessToken)-1])

	return nil
}
