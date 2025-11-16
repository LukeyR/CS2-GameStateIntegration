package events

type SetImagePayload struct {
	Image  string `json:"image"`
	Target int    `json:"target"`
}

type SetImage struct {
	Context string          `json:"context"`
	Event   string          `json:"event"`
	Payload SetImagePayload `json:"payload"`
}
