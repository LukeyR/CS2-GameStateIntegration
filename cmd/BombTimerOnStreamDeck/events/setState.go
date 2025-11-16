package events

type SetStatePayload struct {
	State int `json:"state"`
}

type SetState struct {
	Context string          `json:"context"`
	Event   string          `json:"event"`
	Payload SetStatePayload `json:"payload"`
}
