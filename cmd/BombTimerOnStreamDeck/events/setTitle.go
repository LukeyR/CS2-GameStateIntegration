package events

type SetTitlePayload struct {
	Title  string `json:"title"`
	Target int    `json:"target"`
}

type SetTitle struct {
	Context string          `json:"context"`
	Event   string          `json:"event"`
	Payload SetTitlePayload `json:"payload"`
}
