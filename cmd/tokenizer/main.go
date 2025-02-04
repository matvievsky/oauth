package main

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/matvievsky/oauth/internal/flags"
	_ "github.com/matvievsky/oauth/internal/log"
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
			WithFlag("oid-uri", "", "OpenID URI to authenticate").
			WithFlag("hydra-client-id", "", "Hydra client ID").
			WithFlag("redirect-uri", "", "URI to redirect")
)

var commandToFlags = map[*cobra.Command]flags.Flags{
	{
		Use:   "get",
		Short: "Provides new token",
		RunE: func(cmd *cobra.Command, args []string) error {
			return token.NewClient(
				viper.GetString("oid-uri"),
				viper.GetString("redirect-uri"),
				viper.GetString("hydra-client-id"),
			).Get(
				viper.GetString("hydra-scope"),
				viper.GetString("login-uri"),
				viper.GetString("user-login"),
				viper.GetString("user-password"),
				viper.GetString("consent-uri"),
			)
		},
	}: append(commonFlags, flags.Init().
		WithFlag("hydra-scope", "", "Hydra scope").
		WithFlag("login-uri", "", "URI to log in").
		WithFlag("user-login", "", "user login").
		WithFlag("user-password", "", "user password").
		WithFlag("consent-uri", "", "URI to consent")...),
	{
		Use:   "get-login-challenge",
		Short: "Provides new login challenge",
		RunE: func(cmd *cobra.Command, args []string) error {
			loginChallengeResp, err := token.NewClient(
				viper.GetString("oid-uri"),
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

			parsedURI, err := url.ParseRequestURI(loginChallengeResp.Header.Get(token.Location))
			if err != nil {
				return fmt.Errorf("error parsing redirect URL: %v", err)
			}

			loginChallenge := parsedURI.Query().Get("login_challenge")
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
				viper.GetString("oid-uri"),
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
		slog.Debug("Loaded environment variable", "key", key, "value", viper.GetString(key))
	}
	for command, flags := range commandToFlags {
		if err := flags.Bind(command); err != nil {
			slog.Error("can't bind flags to command", slog.Any(command.Short, err))
			os.Exit(1)
		}
		rootCmd.AddCommand(command)
	}
}

func main() {
	rootCmd.Execute()
}
