package checkers

import (
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"

	"github.com/rs/zerolog/log"
)

func CheckEventPlayerActivityChanged(gsiEvent *structs.GSIEvent) *events.GameEventDetails {
	if gsiEvent.Previous == nil || gsiEvent.Previous.Player == nil {
		return nil
	}

	if gsiEvent.Previous.Player.Activity != "" {
		switch gsiEvent.Player.Activity {
		case structs.PlayerActivityPlaying:
			return &events.GameEventDetails{EventType: events.EventPlayerPlaying}
		case structs.PlayerActivityPaused:
			return &events.GameEventDetails{EventType: events.EventPlayerPaused}
		case structs.PlayerActivityInTextInput:
			return &events.GameEventDetails{EventType: events.EventPlayerInTextInput}
		}
	} else {
		return nil
	}

	originalRequest := gsiEvent.GetOriginalRequestFlat()
	log.Warn().
		Str("event", originalRequest).
		Msg("Unknown Player Activity change event")
	return nil
}
