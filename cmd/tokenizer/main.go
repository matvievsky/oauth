package main

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/matvievsky/oauth/internal/flags"
	"github.com/matvievsky/oauth/internal/token"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Short: "OAuth 2.0 tokenizer",
		Long: `
Ory Hydra based OAuth 2.0 tokenizer:

get - Receive bearer access token
get-login-challenge - Provides new login challenge only
update - Updates existing token
`,
	}

	commonFlags = flags.Init().
			WithFlag("oid-host", "", "OpenID host to authenticate").
			WithFlag("hydra-client-id", "", "Hydra client ID").
			WithFlag("redirect-uri", "", "URI to redirect")
)

var commandToFlags = map[*cobra.Command]flags.Flags{
	{
		Use:   "get",
		Short: "Provides new token",
		RunE: func(cmd *cobra.Command, args []string) error {
			return token.NewClient(
				viper.GetString("oid-host"),
				viper.GetString("redirect-uri"),
				viper.GetString("hydra-client-id"),
			).Get(
				viper.GetString("hydra-scope"),
				viper.GetString("login-url"),
				viper.GetString("user-login"),
				viper.GetString("user-password"),
				viper.GetString("consent-url"),
			)
		},
	}: append(commonFlags, flags.Init().
		WithFlag("hydra-scope", "", "Hydra scope").
		WithFlag("login-url", "", "URI to log in").
		WithFlag("user-login", "", "user login").
		WithFlag("user-password", "", "user password").
		WithFlag("consent-url", "", "URI to consent")...),
	{
		Use:   "get-login-challenge",
		Short: "Provides new login challenge",
		RunE: func(cmd *cobra.Command, args []string) error {
			loginChallengeResp, err := token.NewClient(
				viper.GetString("oid-host"),
				viper.GetString("redirect-uri"),
				viper.GetString("hydra-client-id"),
			).GetLoginChallenge(viper.GetString("hydra-scope"))
			if err != nil {
				return err
			}
			defer loginChallengeResp.Body.Close()

			if !(loginChallengeResp.StatusCode >= 300 && loginChallengeResp.StatusCode < 400) {
				return fmt.Errorf("response status code: %d", loginChallengeResp.StatusCode)
			}

			parsedURL, err := url.Parse(loginChallengeResp.Header.Get(token.Location))
			if err != nil {
				return fmt.Errorf("error parsing redirect URL: %v", err)
			}

			loginChallenge := parsedURL.Query().Get("login_challenge")
			if loginChallenge == "" {
				return fmt.Errorf("login challenge not found in redirect URL")
			}

			fmt.Println(loginChallenge)
			return nil
		},
	}: append(commonFlags, nil...),
	{
		Use:   "update",
		Short: "Updates existing token",
		RunE: func(cmd *cobra.Command, args []string) error {
			return token.NewClient(
				viper.GetString("oid-host"),
				viper.GetString("redirect-uri"),
				viper.GetString("hydra-client-id"),
			).Update(viper.GetString("refresh-token"))
		},
	}: append(commonFlags, flags.Init().
		WithFlag("refresh-token", "", "Token to refresh")...),
}

func init() {
	viper.SetEnvPrefix("TOKENIZER")
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		slog.Info("No .env file found")
	}
	for _, key := range viper.AllKeys() {
		slog.Info("Loaded environment variable", "key", key, "value", viper.GetString(key))
	}
	for command, flags := range commandToFlags {
		if err := flags.Bind(command); err != nil {
			slog.Error("can't bind flags to command", slog.Any(command.Short, err))
		}
		rootCmd.AddCommand(command)
	}
}

func main() {
	rootCmd.Execute()
}
