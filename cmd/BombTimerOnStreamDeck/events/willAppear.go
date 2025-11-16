package events

type WillAppearEvent struct {
	Action  string
	Context string
	Device  string
	Event   string
	Payload struct {
		Controller  string
		Coordinates struct {
			Column int
			Row    int
		}
		IsInMultiAction bool
	}
}
