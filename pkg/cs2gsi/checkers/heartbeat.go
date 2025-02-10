package checkers

import (
	"CS2-GameStateIntegration/pkg/cs2gsi/events"
	"CS2-GameStateIntegration/pkg/cs2gsi/structs"
)

/*
CheckEventHeartbeat
Event Checking functions. Should have signature `GameEventChecker`
*/
func CheckEventHeartbeat(gsiEvent *structs.GSIEvent) *events.GameEventDetails {
	if gsiEvent.Previous == nil {
		return &events.GameEventDetails{EventType: events.EventHeartbeat}
	}
	return nil
}
