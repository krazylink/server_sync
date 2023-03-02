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

var conf bool

// deleteCmd represents the delete command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(c *cobra.Command, args []string) {
		event_id := ""
		params := discordgo.GuildScheduledEventParams{}
		var updater updater
		u := update{}
		updater = &update{&u}
		if conf {
			updater = &confirmUpdate{&u}
		}
		updater.updateEvent(cmd.Session, guild_id, event_id, &params)
	},
}

type updater interface {
	updateEvent(
		s *discordgo.Session,
		guild_id string,
		event_id string,
		params *discordgo.GuildScheduledEventParams,
	)
}

type update struct {
	updater
}

func (u *update) updateEvent(
	s *discordgo.Session,
	guild_id string,
	event_id string,
	params *discordgo.GuildScheduledEventParams,
) {
	e, err := s.GuildScheduledEventEdit(guild_id, event_id, params)
	if err != nil {
		log.Fatal("Failed to update event", event_id, err)
	}
	log.Printf(
		"Updated %v with id %v starting at %v ending at %v\n",
		e.Name,
		e.ID,
		e.ScheduledStartTime,
		e.ScheduledEndTime,
	)
}

type confirmUpdate struct {
	updater
}

func (c *confirmUpdate) updateEvent(
	s *discordgo.Session,
	guild_id string,
	event_id string,
	params *discordgo.GuildScheduledEventParams,
) {
	var ans string
	event := getEvent(s, guild_id, event_id)
	fmt.Printf("Confirm update %s (YES/no): ", event.ID)
	fmt.Scanln(&ans)

	if ans != "YES" {
		log.Fatalf("Couldn't update event %v\n.", event_id)
	}
	c.updater.updateEvent(s, guild_id, event_id, params)
}

type moveAllByDate struct {
	updater
	years  int //years to add
	months int //months to add
	days   int //days to add
}

func (m *moveAllByDate) updateEvent(
	s *discordgo.Session,
	guild_id string,
	event_id string,
	params *discordgo.GuildScheduledEventParams,
) {
	var lister lister
	l := list{}
	lister = &list{&l}
	events := lister.getEvents(s, guild_id)
	for _, event := range events {
		params := discordgo.GuildScheduledEventParams{}
		start := event.ScheduledStartTime.AddDate(m.years, m.months, m.days)
		end := event.ScheduledEndTime.AddDate(m.years, m.months, m.days)
		params.ScheduledStartTime = &start
		params.ScheduledEndTime = &end
		m.updater.updateEvent(s, guild_id, event.ID, &params)
	}
}

type moveAllFromIndex struct {
	updater
	index  int
	years  int
	months int
	days   int
}

func (m *moveAllFromIndex) updateEvent(
	s *discordgo.Session,
	guild_id string,
	event_id string,
	params *discordgo.GuildScheduledEventParams,
) {
	var lister lister
	l := list{}
	lister = &list{&l}

	events := lister.getEvents(s, guild_id)
	for i := m.index; i < len(events); i++ {
		params := discordgo.GuildScheduledEventParams{}
		start := events[i].ScheduledStartTime.AddDate(
			m.years,
			m.months,
			m.days,
		)
		end := events[i].ScheduledEndTime.AddDate(
			m.years,
			m.months,
			m.days,
		)
		params.ScheduledStartTime = &start
		params.ScheduledEndTime = &end
		m.updater.updateEvent(s, guild_id, events[i].ID, &params)
	}
}

func init() {
	eventCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	updateCmd.Flags().BoolVar(&conf, "confirm", true, "Confirm update")
}
