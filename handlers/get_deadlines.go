package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	dg "github.com/bwmarrin/discordgo"
	"github.com/juankhusuma/scele-bot/handlers/auth"
)

type Payload struct {
	Message string    `json:"message"`
	Status  int       `json:"status"`
	Data    Deadlines `json:"data"`
}

type Deadlines struct {
	Period string `json:"period"`
	Items  []Item `json:"items"`
}

type Item struct {
	Title string `json:"title"`
	Due   string `json:"due"`
	Unix  int    `json:"unix"`
	Url   string `json:"url"`
}

func GetDeadlines(s *dg.Session, i *dg.InteractionCreate, period string) {
	s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
		Type: dg.InteractionResponseChannelMessageWithSource,
		Data: &dg.InteractionResponseData{
			Content: "Getting your deadlines...",
			Flags:   dg.MessageFlagsEphemeral,
		},
	})
	token, err := auth.Login(s, i)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &dg.WebhookParams{
			Content: "Error: " + err.Error(),
		})
		return
	}
	client := http.Client{}
	req, err := http.NewRequest("GET", os.Getenv("PROXY_URL")+"/deadlines/"+period, nil)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &dg.WebhookParams{
			Content: "Error: " + err.Error(),
		})
		return
	}
	req.Header.Set("X-Moodle-Session", token)
	res, err := client.Do(req)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &dg.WebhookParams{
			Content: "Error: " + err.Error(),
		})
		return
	}
	var data Payload
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	json.Unmarshal(b, &data)
	embeds := []*dg.MessageEmbedField{}
	for _, item := range data.Data.Items {
		embeds = append(embeds, &dg.MessageEmbedField{
			Name:   item.Title,
			Value:  fmt.Sprintf("Due: %s\n[Link](%s)", item.Due, item.Url),
			Inline: false,
		})
	}
	s.FollowupMessageCreate(i.Interaction, true, &dg.WebhookParams{
		Content: "Here are your deadlines",
		Embeds: []*dg.MessageEmbed{
			{
				Title:       fmt.Sprintf("Deadlines for %s", i.Member.User.Username),
				Description: fmt.Sprintf("Period: %s", data.Data.Period),
				Fields:      embeds,
				Timestamp:   time.Now().Format(time.RFC3339),
				Author: &dg.MessageEmbedAuthor{
					Name:    "Scele Bot",
					IconURL: "https://cdn.discordapp.com/attachments/820295359158681603/1138442585854185502/uwiw.png",
				},
				Color: 0xffd500,
				Footer: &dg.MessageEmbedFooter{
					Text:    "Scraped from Scele",
					IconURL: "https://cdn.discordapp.com/attachments/820295359158681603/1138442585854185502/uwiw.png",
				},
			},
		},
		Username: "Scele Bot",
	})
}
