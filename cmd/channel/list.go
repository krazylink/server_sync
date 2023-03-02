/*
Copyright © 2023 Paul Huckabee <paul.huckabee@gmail.com>

*/
package channel

import (
	"discordctl/cmd"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

var printList bool

// channelCmd represents the channel command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list channels",
	Long:  `list channels`,
	Run: func(c *cobra.Command, args []string) {
		var lister lister
		l := list{}
		lister = &list{&l}
		if printList {
			lister = &andPrint{lister}
		}
		lister.getChannels(cmd.Session, guild_id)
	},
}

func init() {
	channelCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(
		&printList,
		"print",
		true,
		"Print list to stdout",
	)
}

type lister interface {
	getChannels(s *discordgo.Session, guild_id string) []*discordgo.Channel
}

type list struct {
	lister
}

func (d *list) getChannels(s *discordgo.Session, guild_id string) []*discordgo.Channel {
	channels, err := s.GuildChannels(guild_id)
	if err != nil {
		log.Fatal("Failed to get channel list: ", err)
	}
	return channels
}

type andPrint struct {
	lister
}

func (a *andPrint) getChannels(s *discordgo.Session, guild_id string) []*discordgo.Channel {
	channels := a.lister.getChannels(s, guild_id)
	for i, channel := range channels {
		fmt.Println(i, channel.Name, channel.ID, channel.Type, channel.ParentID)
	}
	return channels
}
