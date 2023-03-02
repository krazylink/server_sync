/*
Copyright © 2023 Paul Huckabee <paul.huckabee@gmail.com>

*/
package channel

import (
	"discordctl/cmd"

	"github.com/spf13/cobra"
)

var guild_id string

var channelCmd = &cobra.Command{
	Use:   "channel",
	Short: "create update delete discord channels",
	Long:  `Commands to list, create, update and delete discord channels.`,
}

func init() {
	cmd.RootCmd.AddCommand(channelCmd)
	channelCmd.PersistentFlags().StringVarP(
		&guild_id,
		"guild_id",
		"g",
		"",
		"Discord guild (server) id",
	)
	channelCmd.MarkPersistentFlagRequired("guild_id")
}
