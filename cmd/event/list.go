/*
Copyright © 2023 Paul Huckabee <paul.huckabee@gmail.com>

*/
package event

import (
	"discordctl/cmd"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

var printList bool

// deleteCmd represents the delete command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(c *cobra.Command, args []string) {
		var lister lister
		l := list{}
		lister = &list{&l}
		if printList {
			lister = &andPrint{lister}
		}
		lister.getEvents(cmd.Session, guild_id)

	},
}

func getEvent(s *discordgo.Session, guild_id string, event_id string) *discordgo.GuildScheduledEvent {
	event, err := s.GuildScheduledEvent(guild_id, event_id, false)
	if err != nil {
		log.Fatal("Failed to get event:", err)
	}
	return event
}

type lister interface {
	getEvents(s *discordgo.Session, guild_id string) []*discordgo.GuildScheduledEvent
}

type list struct {
	lister
}

func (l *list) getEvents(s *discordgo.Session, guild_id string) []*discordgo.GuildScheduledEvent {
	events, err := s.GuildScheduledEvents(guild_id, false)
	if err != nil {
		log.Fatal("Failed to get events:", err)
	}
	return events
}

type andPrint struct {
	lister
}

func (ap *andPrint) getEvents(s *discordgo.Session, guild_id string) []*discordgo.GuildScheduledEvent {
	events := ap.lister.getEvents(s, guild_id)
	fmt.Println("Index | ID | Name | Start | End")
	for i, event := range events {
		fmt.Println(i, event.ID, event.Name, event.ScheduledStartTime, event.ScheduledEndTime)
	}
	return events
}

func init() {
	eventCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(
		&printList,
		"print",
		true,
		"Print list",
	)

}
