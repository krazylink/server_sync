package cmd

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func AuthSession(token string) *discordgo.Session {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Failed to created authenticated session:", err)
	}
	return s
}
