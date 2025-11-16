package events

type RegisterEvent struct {
	Event string `json:"event"`
	UUID  string `json:"uuid"`
}

type GenericReceiveEvent struct {
	Event string `json:"event"`
}
