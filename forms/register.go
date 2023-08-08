package forms

import dg "github.com/bwmarrin/discordgo"

func CreateRegisterForm() []dg.MessageComponent {
	return []dg.MessageComponent{
		dg.ActionsRow{
			Components: []dg.MessageComponent{
				dg.TextInput{
					CustomID:    "sso_username",
					Label:       "SSO Username",
					Placeholder: "Enter your SSO username here...",
					Style:       dg.TextInputShort,
					Required:    true,
				},
			},
		},
		dg.ActionsRow{
			Components: []dg.MessageComponent{
				dg.TextInput{
					CustomID:    "sso_password",
					Label:       "SSO Password",
					Placeholder: "Enter your SSO password here...",
					Style:       dg.TextInputShort,
					Required:    true,
				},
			},
		},
	}
}
