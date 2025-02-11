package checkers

import (
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"
)

func CheckEventBombPlanted(gsiEvent *structs.GSIEvent) *events.GameEventDetails {
	if gsiEvent.Added != nil && gsiEvent.Added.Round != nil && gsiEvent.Added.Round.Bomb && gsiEvent.Round.Bomb == "planted" {
		return &events.GameEventDetails{EventType: events.EventBombPlanted}
	}
	return nil
}

func CheckEventBombExploded(gsiEvent *structs.GSIEvent) *events.GameEventDetails {
	if gsiEvent.Added != nil && gsiEvent.Added.Round != nil && gsiEvent.Added.Round.WinTeam && gsiEvent.Round.Bomb == "exploded" {
		return &events.GameEventDetails{EventType: events.EventBombExploded}
	}
	return nil
}

func CheckEventBombDefused(gsiEvent *structs.GSIEvent) *events.GameEventDetails {
	if gsiEvent.Added != nil && gsiEvent.Added.Round != nil && gsiEvent.Added.Round.WinTeam && gsiEvent.Round.Bomb == "defused" {
		return &events.GameEventDetails{EventType: events.EventBombDefused}
	}
	return nil
}
