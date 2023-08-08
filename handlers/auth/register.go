package auth

import (
	"encoding/json"
	"fmt"
	"os"

	dg "github.com/bwmarrin/discordgo"
	"github.com/juankhusuma/scele-bot/forms"
	"github.com/juankhusuma/scele-bot/types"
	"github.com/juankhusuma/scele-bot/utils"
)

func RegisterSSO(s *dg.Session, i *dg.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
		Type: dg.InteractionResponseModal,
		Data: &dg.InteractionResponseData{
			CustomID:   "bot_register_" + i.Interaction.Member.User.ID,
			Title:      "Register Your SSO",
			Components: forms.CreateRegisterForm(),
			Flags:      dg.MessageFlagsEphemeral,
		},
	})
}

func HandleRegistrationSubmission(s *dg.Session, i *dg.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
		Type: dg.InteractionResponseChannelMessageWithSource,
		Data: &dg.InteractionResponseData{
			Content: "Thanks for submitting, please wait...",
			Flags:   dg.MessageFlagsEphemeral,
		},
	})
	data := i.ModalSubmitData()
	var mc types.MessageComponent
	sd := make(map[string]string)
	for _, c := range data.Components {
		j, _ := c.MarshalJSON()
		json.Unmarshal(j, &mc)
		sd[mc.Components[0].CustomID] = mc.Components[0].Value
	}
	db, err := utils.GetDB()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	rows, _ := db.Query("SELECT * FROM users")
	for rows.Next() {
		var id string
		var username string
		var password string
		rows.Scan(&id, &username, &password)
		if id == i.Interaction.Member.User.ID {
			s.FollowupMessageCreate(i.Interaction, true, &dg.WebhookParams{
				Content: "You have already registered previously, your old credentials will be overwritten!",
			})
		}
	}
	key := i.Interaction.Member.User.ID + os.Getenv("SALT")
	sd["sso_password"], err = utils.Encrypt(key, sd["sso_password"])
	if err != nil {
		fmt.Println(err)
	}
	sd["sso_username"], err = utils.Encrypt(key, sd["sso_username"])
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET username = $2, password = $3",
		i.Interaction.Member.User.ID, sd["sso_username"], sd["sso_password"])
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &dg.WebhookParams{
			Content: "Error: " + err.Error(),
		})
	}

	s.FollowupMessageCreate(i.Interaction, true, &dg.WebhookParams{
		Content: "Successfully registered!",
	})
}
