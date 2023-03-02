/*
Copyright © 2023 Paul Huckabee <paul.huckabee@gmail.com>

*/
package event

import (
	"discordctl/cmd"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

var filename string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create discord scheduled event",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(c *cobra.Command, args []string) {
		params := discordgo.GuildScheduledEventParams{}

		var creator creator
		cre := create{}
		creator = &create{&cre}
		if filename != "" {
			creator = &fromFile{&cre, filename}
		}
		creator.createEvent(cmd.Session, guild_id, &params)
	},
}

type creator interface {
	createEvent(
		s *discordgo.Session,
		guild_id string,
		params *discordgo.GuildScheduledEventParams,
	)
}

type create struct {
	creator
}

func (c *create) createEvent(
	s *discordgo.Session,
	guild_id string,
	params *discordgo.GuildScheduledEventParams,
) {
	event, err := s.GuildScheduledEventCreate(guild_id, params)
	if err != nil {
		log.Fatal("Failed to create event:", err)
	}
	log.Printf(
		"Created %v %v starting at %v ending at %v\n",
		event.Name,
		event.Description,
		event.ScheduledStartTime,
		event.ScheduledEndTime,
	)
}

type Events struct {
	Events []discordgo.GuildScheduledEventParams `json:"events"`
}

type fromFile struct {
	creator
	fromFile string
}

func (f *fromFile) createEvent(
	s *discordgo.Session,
	guild_id string,
	params *discordgo.GuildScheduledEventParams,
) {
	file, err := os.Open(f.fromFile)
	if err != nil {
		log.Fatal("Failed to open file:", err)
	}

	b, _ := ioutil.ReadAll(file)
	var events Events

	json.Unmarshal(b, &events)
	for _, event := range events.Events {
		f.creator.createEvent(s, guild_id, &event)
	}
}
func init() {
	eventCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(
		&filename,
		"fromFile",
		"",
		"JSON file to read events from",
	)

}
