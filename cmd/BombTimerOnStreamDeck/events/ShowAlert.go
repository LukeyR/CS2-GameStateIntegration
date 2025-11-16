package events

type ShowAlert struct {
	Context string `json:"context"`
	Event   string `json:"event"`
}
