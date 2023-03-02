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

var conf bool
var channel string
var FromIndex int

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	Long:  `Delete a channel`,
	Run: func(c *cobra.Command, args []string) {
		var deleter deleter
		d := defaultDeleter{}
		deleter = &defaultDeleter{&d}
		if conf {
			deleter = &confirm{deleter}
		}
		if FromIndex > -1 {
			deleter = &fromIndex{deleter, FromIndex}
		}
		deleter.deleteChannel(cmd.Session, channel)
	},
}

type deleter interface {
	deleteChannel(s *discordgo.Session, channel_id string)
}

type defaultDeleter struct {
	deleter
}

func (d *defaultDeleter) deleteChannel(
	s *discordgo.Session,
	channel_id string,
) {
	c, err := s.ChannelDelete(channel_id)
	if err != nil {
		log.Println("Delete failed", err)
	}
	log.Println("Channel deleted", c.ID)
}

type confirm struct {
	deleter
}

func (c *confirm) deleteChannel(
	s *discordgo.Session,
	channel_id string,
) {
	var ans string
	fmt.Printf("Delete channel %v (YES/no)? ", channel_id)
	fmt.Scanln(&ans)
	if ans == "YES" {
		c.deleter.deleteChannel(s, channel_id)
	}
}

type fromIndex struct {
	deleter
	index int
}

func (ind *fromIndex) deleteChannel(
	s *discordgo.Session,
	channel_id string,
) {
	var lister lister
	l := list{}
	lister = &list{&l}
	channels := lister.getChannels(s, guild_id)
	for i := ind.index; i < len(channels); i++ {
		ind.deleter.deleteChannel(s, channels[i].ID)
	}
}

func init() {

	channelCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(
		&channel,
		"channelID",
		"",
		"Channel id to be deleted",
	)
	deleteCmd.Flags().BoolVar(
		&conf,
		"confirm",
		true,
		"Confirm delete",
	)
	deleteCmd.Flags().IntVar(
		&FromIndex,
		"from_index",
		-1,
		"Index to start deleting from",
	)
}
