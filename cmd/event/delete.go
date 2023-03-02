/*
Copyright © 2023 Paul Huckabee <paul.huckabee@gmail.com>

*/
package event

import (
	"discordctl/cmd"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

var confDelete bool
var event_id string
var startIndex int
var index int

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(c *cobra.Command, args []string) {
		var deleter deleter
		d := del{}
		deleter = &del{&d}
		if confDelete {
			deleter = &confirmDelete{deleter}
		}
		if startIndex > -1 {
			deleter = &fromIndex{deleter, startIndex}
		}
		if index > -1 {
			deleter = &byIndex{deleter, index}
		}

		deleter.deleteEvent(cmd.Session, guild_id, event_id)
	},
}

type deleter interface {
	deleteEvent(
		s *discordgo.Session,
		guild_id string,
		event_id string,
	)
}

type del struct {
	deleter
}

func (d *del) deleteEvent(
	s *discordgo.Session,
	guild_id string,
	event_id string,
) {
	err := s.GuildScheduledEventDelete(guild_id, event_id)
	if err != nil {
		log.Fatal("Failed to delete event:", err)
	}
	log.Printf("%v deleted.\n", event_id)
}

type confirmDelete struct {
	deleter
}

func (c *confirmDelete) deleteEvent(
	s *discordgo.Session,
	guild_id string,
	event_id string,
) {
	var ans string
	fmt.Printf("Confirm %v delete (YES/no) ", event_id)
	fmt.Scanln(&ans)
	if ans != "YES" {
		fmt.Println("Delete aborted.")
		os.Exit(-1)
	}
	c.deleter.deleteEvent(s, guild_id, event_id)
}

type byIndex struct {
	deleter
	index int
}

func (i *byIndex) deleteEvent(
	s *discordgo.Session,
	guild_id string,
	event_id string,
) {
	var lister lister
	l := list{}
	lister = &list{&l}
	events := lister.getEvents(s, guild_id)
	i.deleter.deleteEvent(s, guild_id, events[i.index].ID)
}

type fromIndex struct {
	deleter
	index int
}

func (ind *fromIndex) deleteEvent(
	s *discordgo.Session,
	guild_id string,
	event_id string,
) {
	var lister lister
	l := list{}
	lister = &list{&l}
	events := lister.getEvents(s, guild_id)

	for i := ind.index; i < len(events); i++ {
		ind.deleter.deleteEvent(s, guild_id, events[i].ID)
	}
}

func init() {
	eventCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(
		&event_id,
		"eventID",
		"e",
		"",
		"Discord event ID to delete",
	)
	deleteCmd.Flags().IntVar(
		&startIndex,
		"fromIndex",
		-1,
		"Index to start deleting at",
	)
	deleteCmd.Flags().IntVar(
		&index,
		"index",
		-1,
		"Index to deleting",
	)
	deleteCmd.Flags().BoolVar(
		&confDelete,
		"confirm",
		true,
		"Confirm delete",
	)
}
