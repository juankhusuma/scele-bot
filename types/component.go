package types

type MessageComponent struct {
	Components []Component `json:"components"`
}

type Component struct {
	CustomID string `json:"custom_id"`
	Label    string `json:"label"`
	Style    int    `json:"style"`
	Value    string `json:"value"`
	Required bool   `json:"required"`
	Type     int    `json:"type"`
}
