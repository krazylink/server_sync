/*
Copyright © 2023 Paul Huckabee <paul.huckabee@gmail.com>

*/
package channel

import (
	"discordctl/cmd"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
)

var createFile string
var Name string
var Type discordgo.ChannelType = discordgo.ChannelTypeGuildText
var Topic string
var Bitrate int
var Userlimit int
var RateLimitPerUser int
var Position int

//permission_overwrites
var ParentID string
var Nsfw bool

// channelCmd represents the channel command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create channels",
	Long:  `create channels`,
	Run: func(c *cobra.Command, args []string) {
		var data discordgo.GuildChannelCreateData
		var creator creator
		cre := defaultChannelCreator{}
		//
		if createFile != "" {
			creator = &fileCreator{&cre, createFile}
		} else {
			creator = &defaultChannelCreator{&cre}
			data.Name, _ = c.Flags().GetString("name")
			data.Type = Type
			if data.Type == discordgo.ChannelTypeGuildText || data.Type == discordgo.ChannelTypeGuildNews || data.Type == discordgo.ChannelTypeGuildForum {
				data.Topic, _ = c.Flags().GetString("topic")
			}
			data.Bitrate, _ = c.Flags().GetInt("bitrate")
			data.UserLimit, _ = c.Flags().GetInt("user_limit")
			data.RateLimitPerUser, _ = c.Flags().GetInt("rate_limit_per_user")
			data.Position, _ = c.Flags().GetInt("position")
			//TODO(krazylink): impliment this flag
			//data.PermissionOverwrites = c.Flags().
			if data.Type != discordgo.ChannelTypeGuildCategory {
				data.ParentID, _ = c.Flags().GetString("parent_id")
			} else {
				log.Println("Type Category can not have parent. Ignoring parent_id.")
			}
			data.NSFW, _ = c.Flags().GetBool("nsfw")
		}
		creator.createChannel(cmd.Session, guild_id, data)
	},
}

type creator interface {
	createChannel(
		s *discordgo.Session,
		guild_id string,
		data discordgo.GuildChannelCreateData,
	)
}

type defaultChannelCreator struct {
	creator
}

func (d *defaultChannelCreator) createChannel(
	s *discordgo.Session,
	guild_id string,
	data discordgo.GuildChannelCreateData,
) {
	_, err := s.GuildChannelCreateComplex(guild_id, data)
	if err != nil {
		log.Fatal("Failed to create channel", err)
	}
}

type Channels struct {
	Channels   []discordgo.GuildChannelCreateData `json:"channels"`
	Categories []Category                         `json:"categories,omitempty"`
}
type Category struct {
	Name     string                             `json:"name"`
	Channels []discordgo.GuildChannelCreateData `json:"channels,omitempty"`
}

type fileCreator struct {
	creator
	fileName string
}

func (f *fileCreator) createChannel(
	s *discordgo.Session,
	guild_id string,
	data discordgo.GuildChannelCreateData,
) {
	file, err := os.Open(f.fileName)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	channels := Channels{}
	json.Unmarshal(b, &channels)

	//create unparented channels.
	for _, c := range channels.Channels {
		f.creator.createChannel(s, guild_id, c)
	}
	for _, c := range channels.Categories {
		data.Name = c.Name
		data.Type = 4
		f.creator.createChannel(s, guild_id, data)
		for _, ch := range c.Channels {
			ch.ParentID = getCategoryID(s, guild_id, c.Name)
			f.creator.createChannel(s, guild_id, ch)
		}
	}
}

func getCategoryID(s *discordgo.Session, guild_id string, category string) string {
	channels, err := s.GuildChannels(guild_id)
	if err != nil {
		log.Println("Failed to get channels:", err)
	}
	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildCategory {
			if channel.Name == category {
				return channel.ID
			}
		}
	}
	return "not found"
}

func hasCategory(s *discordgo.Session, guild_id, category string) bool {
	channels, err := s.GuildChannels(guild_id)
	if err != nil {
		log.Fatal("Failed to get channels:", err)
	}
	for _, channel := range channels {
		if channel.Name == category {
			return true
		}
	}
	return false
}

var ChannelTypeIDs = map[discordgo.ChannelType][]string{
	discordgo.ChannelTypeGuildText:     {"Text"},
	discordgo.ChannelTypeGuildVoice:    {"Voice"},
	discordgo.ChannelTypeGuildCategory: {"Category"},
}

func init() {
	channelCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(
		&createFile,
		"from_file",
		"",
		"File to use for input",
	)
	createCmd.Flags().StringVarP(
		&Name,
		"name",
		"n",
		"",
		"Channel name.",
	)
	createCmd.Flags().VarP(
		enumflag.New(&Type, "type", ChannelTypeIDs, enumflag.EnumCaseSensitive),
		"type",
		"t",
		"Channel type (0: text channel, 2: voice channel, :4 category).",
	)
	createCmd.Flags().StringVarP(
		&Topic,
		"topic",
		"T",
		"",
		"Channel Topic",
	)
	createCmd.Flags().IntVarP(
		&Bitrate,
		"bitrate",
		"b",
		8000,
		"bitrate (in bits) of the voice channel; min 8000",
	)
	createCmd.Flags().IntVar(
		&Userlimit,
		"user_limit",
		0,
		"User limit of the voice channel",
	)
	createCmd.Flags().IntVar(
		&RateLimitPerUser,
		"rate_limit_per_user",
		0,
		"Amount of seconds a user has to wait before sending another message (0-21600)",
	)
	createCmd.Flags().IntVar(
		&Position,
		"Position",
		0,
		"Sorting position of the channel",
	)
	//permission_overwrites
	createCmd.Flags().StringVarP(
		&ParentID,
		"parent_id",
		"p",
		"",
		"ID of the parent catagory for a channel",
	)
	createCmd.Flags().BoolVar(
		&Nsfw,
		"nsfw",
		false,
		"Whether the channel is NSFW",
	)
}
