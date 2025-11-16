package events

type LogMessagePayload struct {
	Message string `json:"message"`
}

type LogMessage struct {
	Event   string            `json:"event"`
	Payload LogMessagePayload `json:"payload"`
}
