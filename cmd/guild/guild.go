/*
Copyright © 2023 Paul Huckabee <paul.huckabee@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// guildCmd represents the guild command
var guildCmd = &cobra.Command{
	Use:   "guild",
	Short: "Discord guilds (servers)",
	Long:  `list, create, update, and delete discord servers`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Help())
	},
}

func init() {
	RootCmd.AddCommand(guildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// guildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// guildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
