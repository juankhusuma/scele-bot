package handlers

import (
	"strings"

	dg "github.com/bwmarrin/discordgo"
	"github.com/juankhusuma/scele-bot/handlers/auth"
)

func InteractionCreate(s *dg.Session, i *dg.InteractionCreate) {
	if i.Type == dg.InteractionModalSubmit {
		modalID := i.ModalSubmitData().CustomID
		if strings.Contains(modalID, "bot_register") {
			auth.HandleRegistrationSubmission(s, i)
			return
		}
		return
	}

	data := i.ApplicationCommandData()
	switch data.Name {
	case "register":
		auth.RegisterSSO(s, i)
	case "deadlines":
		GetDeadlines(s, i, data.Options[0].StringValue())
	}
}
