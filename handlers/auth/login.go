package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	dg "github.com/bwmarrin/discordgo"
	"github.com/juankhusuma/scele-bot/utils"
)

type LoginReponse struct {
	Message string    `json:"message"`
	Status  int       `json:"status"`
	Data    LoginData `json:"data"`
}

type LoginData struct {
	MoodleSession string `json:"moodle_session"`
}

func Login(s *dg.Session, i *dg.InteractionCreate) (string, error) {
	db, err := utils.GetDB()
	if err != nil {
		return "", err
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users WHERE id = $1", i.Interaction.Member.User.ID)
	if err != nil {
		return "", err
	}
	var id string
	var username string
	var password string
	for rows.Next() {
		rows.Scan(&id, &username, &password)
	}
	if password == "" {
		s.FollowupMessageCreate(i.Interaction, true, &dg.WebhookParams{
			Content: "You have not registered yet. Please register first.",
		})
		return "", nil
	}
	username, err = utils.Decrypt(i.Interaction.Member.User.ID+os.Getenv("SALT"), username)
	if err != nil {
		fmt.Println(err)
	}
	password, err = utils.Decrypt(i.Interaction.Member.User.ID+os.Getenv("SALT"), password)
	if err != nil {
		fmt.Println(err)
	}

	client := http.Client{}
	creds := map[string]string{
		"username": username,
		"password": password,
	}
	jval, err := json.Marshal(creds)
	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Post(os.Getenv("PROXY_URL")+"/auth/login", "application/json", bytes.NewBuffer(jval))
	if err != nil {
		fmt.Println(err)
	}
	b, _ := io.ReadAll(res.Body)
	var data LoginReponse
	json.Unmarshal(b, &data)

	if data.Status == 200 {
		s.FollowupMessageCreate(i.Interaction, true, &dg.WebhookParams{
			Content: fmt.Sprintf("Logged in as %s", username),
		})
		return data.Data.MoodleSession, nil
	}
	return "", nil
}
