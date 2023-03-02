/*
Copyright © 2023 Paul Huckabee <paul.huckabee@gmail.com>

*/
package event

import (
	"discordctl/cmd"

	"github.com/spf13/cobra"
)

var Confirm bool
var guild_id string

// eventCmd represents the event command
var eventCmd = &cobra.Command{
	Use:              "event",
	Short:            "Discord scheduled events",
	Long:             `list, create, update and delete discord scheduled events`,
	TraverseChildren: true,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println(cmd.Help())
	//	},
}

func init() {
	cmd.RootCmd.AddCommand(eventCmd)
	eventCmd.PersistentFlags().StringVarP(
		&guild_id,
		"guildID",
		"g",
		"1072385994097168394",
		"Discord guild (server) ID to use",
	)
	eventCmd.PersistentFlags().BoolVar(
		&Confirm,
		"confirm",
		Confirm,
		"Confirm action",
	)
}
