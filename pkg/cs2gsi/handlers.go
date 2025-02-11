package cs2gsi

import (
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/checkers"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"
)

/*
	Main struct returned is GameEventDetails.
	Alongside the event type raised, it can (anonymously) contain extra data about the event
*/

type GameEventChecker func(gsiEvent *structs.GSIEvent) *events.GameEventDetails

var eventCheckers = []GameEventChecker{
	checkers.CheckEventHeartbeat,
	checkers.CheckEventWeaponsChanged,
	checkers.CheckEventPlayerActivityChanged,
	checkers.CheckEventPlayerAliveStatusChanged,
	checkers.CheckEventPlayerHealthChanged,
	checkers.CheckEventPlayerArmourChanged,
}

/*
Function for subscribing to GameEvents
*/
type gameEventHandlerCallback func(*structs.GSIEvent, events.GameEventDetails)

var gameEventHandlers = make(map[events.GameEvent][]gameEventHandlerCallback)

func RegisterEventHandler(event events.GameEvent, handler gameEventHandlerCallback) {
	gameEventHandlers[event] = append(gameEventHandlers[event], handler)
}

func findEvents(gsiEvent *structs.GSIEvent) []events.GameEventDetails {
	gameEvents := make([]events.GameEventDetails, 0)
	for _, checker := range eventCheckers {
		res := checker(gsiEvent)
		if res != nil {
			gameEvents = append(gameEvents, *res)
		}
	}
	return gameEvents
}
