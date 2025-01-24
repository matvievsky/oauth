package main

import (
	"github.com/matvievsky/oauth/internal/token"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "OAuth 2.0 tokenizer",
	Long: `
Ory Hydra based OAuth 2.0 tokenizer:

get - receive bearer access token`,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Provides new token",
	RunE: func(cmd *cobra.Command, args []string) error {
		oidHost, _ := cmd.Flags().GetString("oid-host")
		redirectURI, _ := cmd.Flags().GetString("redirect-uri")
		hydraClientID, _ := cmd.Flags().GetString("hydra-client-id")

		hydraScope, _ := cmd.Flags().GetString("hydra-scope")
		loginURI, _ := cmd.Flags().GetString("login-url")
		userLogin, _ := cmd.Flags().GetString("user-login")
		userPassword, _ := cmd.Flags().GetString("user-password")
		consentURI, _ := cmd.Flags().GetString("consent-url")

		return token.NewClient(oidHost, redirectURI, hydraClientID).Get(hydraScope, loginURI, userLogin, userPassword, consentURI)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates existing token",
	RunE: func(cmd *cobra.Command, args []string) error {
		oidHost, _ := cmd.Flags().GetString("oid-host")
		redirectURI, _ := cmd.Flags().GetString("redirect-uri")
		hydraClientID, _ := cmd.Flags().GetString("hydra-client-id")

		refreshToken, _ := cmd.Flags().GetString("refresh-token")

		return token.NewClient(oidHost, redirectURI, hydraClientID).Update(refreshToken)
	},
}

func init() {
	rootCmd.AddCommand(getCmd, updateCmd)

	for _, cmd := range rootCmd.Commands() {
		cmd.Flags().String("oid-host", "oid.dev1.kassirplus.ru", "OpenID host to authenticate")
		cmd.Flags().String("hydra-client-id", "", "Hydra client ID")
		cmd.Flags().String("redirect-uri", "https://kirov-next.dev1.kassirplus.ru", "URI to redirect")
	}

	getCmd.Flags().String("hydra-scope", "openid offline offline_access", "Hydra scope")
	getCmd.Flags().String("login-url", "https://api.dev1.kassirplus.ru/kassir.idm.Idm/LogIn", "URI to log in")
	getCmd.Flags().String("user-login", "", "user login")
	getCmd.Flags().String("user-password", "", "user password")
	getCmd.Flags().String("consent-url", "https://api.dev1.kassirplus.ru/kassir.idm.Idm/Consent", "URI to consent")
	updateCmd.Flags().String("refresh-token", "", "Token to refresh")

}

func main() {
	rootCmd.Execute()
}
