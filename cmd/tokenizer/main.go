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
		hydraClientID, _ := cmd.Flags().GetString("hydra-client-id")
		hydraScope, _ := cmd.Flags().GetString("hydra-scope")
		oidHost, _ := cmd.Flags().GetString("oid-host")
		redirectURI, _ := cmd.Flags().GetString("redirect-uri")
		loginURI, _ := cmd.Flags().GetString("login-url")
		userLogin, _ := cmd.Flags().GetString("user-login")
		userPassword, _ := cmd.Flags().GetString("user-password")
		consentURI, _ := cmd.Flags().GetString("consent-url")

		return token.Get(oidHost, hydraClientID, hydraScope, redirectURI, loginURI, userLogin, userPassword, consentURI)
	},
}

func init() {
	getCmd.Flags().String("hydra-client-id", "", "Hydra client ID")
	getCmd.Flags().String("hydra-scope", "openid offline offline_access", "Hydra scope")
	getCmd.Flags().String("oid-host", "oid.dev1.kassirplus.ru", "OpenID host to authenticate")
	getCmd.Flags().String("redirect-uri", "https://kirov-next.dev1.kassirplus.ru", "URI to redirect")
	getCmd.Flags().String("login-url", "https://api.dev1.kassirplus.ru/kassir.idm.Idm/LogIn", "URI to log in")
	getCmd.Flags().String("user-login", "", "user login")
	getCmd.Flags().String("user-password", "", "user password")
	getCmd.Flags().String("consent-url", "https://api.dev1.kassirplus.ru/kassir.idm.Idm/Consent", "URI to consent")
}

func main() {
	rootCmd.AddCommand(getCmd)
	rootCmd.Execute()
}
