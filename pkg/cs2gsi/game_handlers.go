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

type gameEventChecker func(gsiEvent *structs.GSIEvent) *events.GameEventDetails

var eventCheckers = []gameEventChecker{
	checkers.CheckEventHeartbeat,
	checkers.CheckEventWeaponsChanged,
	checkers.CheckEventPlayerActivityChanged,
	checkers.CheckEventPlayerAliveStatusChanged,
	checkers.CheckEventPlayerHealthChanged,
	checkers.CheckEventPlayerArmourChanged,
	checkers.CheckEventBombPlanted,
	checkers.CheckEventBombExploded,
	checkers.CheckEventBombDefused,
}

/*
Function for subscribing to GameEvents
*/
type gameEventHandlerCallback func(*structs.GSIEvent, events.GameEventDetails)

var gameEventHandlers = make(map[events.GameEvent][]gameEventHandlerCallback)
var gameNonEventHandlers = make([]func(event *structs.GSIEvent), 0)

// RegisterEventHandler is the main way to register a handler callback for when an events.GameEvent is raised
func RegisterEventHandler(event events.GameEvent, handler gameEventHandlerCallback) {
	gameEventHandlers[event] = append(gameEventHandlers[event], handler)
}

// RegisterGlobalHandler is a function that will call handler anytime a RECOGNISED event is raised
func RegisterGlobalHandler(handler gameEventHandlerCallback) {
	for enum := range events.EnumToEventName {
		RegisterEventHandler(enum, handler)
	}
}

// RegisterNonEventHandler calls the handler callback provided if a game event happens that the library does not currently support
// Useful for using the package as a base should you wish to extend it, or for debugging
func RegisterNonEventHandler(handler func(event *structs.GSIEvent)) {
	gameNonEventHandlers = append(gameNonEventHandlers, handler)
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
