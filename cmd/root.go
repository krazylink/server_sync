/*
Copyright © 2023 Paul Huckabee <paul.huckabee@gmail.com>

*/
package cmd

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

var AuthToken string
var Session *discordgo.Session

var RootCmd = &cobra.Command{
	Use:   "discordctl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		tok, _ := cmd.Flags().GetString("auth_token")
		Session = AuthSession(tok)
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(
		&AuthToken,
		"auth_token",
		"a",
		"",
		"OAuth2 token to use.",
	)
	RootCmd.MarkPersistentFlagRequired("auth_token")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
