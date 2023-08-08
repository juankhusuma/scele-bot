package handlers

import (
	"os"

	dg "github.com/bwmarrin/discordgo"
)

func MessageCreate(s *dg.Session, m *dg.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}
	if m.Content == "!scele" {
		_, err := s.ApplicationCommandBulkOverwrite(os.Getenv("APP_ID"), m.GuildID, []*dg.ApplicationCommand{
			{
				Name:        "register",
				Description: "Register your SSO credentials before using Scele Bot",
			},
			{
				Name:        "deadlines",
				Description: "Get your deadlines from Scele",
				Options: []*dg.ApplicationCommandOption{
					{Name: "period", Description: "The period of the deadlines with the format of (dd-mm-yyy)", Type: dg.ApplicationCommandOptionString, Required: true},
				},
			},
		})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		}
		s.ChannelMessageSend(m.ChannelID, "Scele Bot initialized!")
	}
}
