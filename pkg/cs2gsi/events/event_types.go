package events

import (
	"strconv"

	"github.com/rs/zerolog/log"
)

type GameEvent int

/*
All possible event that can be raised, as well as functions needed to check for them
*/
const (
	EventHeartbeat GameEvent = iota

	EventPlayerPaused
	EventPlayerPlaying
	EventPlayerInTextInput

	EventPlayerWeaponUse
	EventPlayerWeaponReloadStarted
	EventPlayerWeaponReloadFinished
	EventPlayerWeaponChanged
	EventPlayerWeaponAdded
	EventPlayerWeaponRemoved
	EventPlayerActiveWeaponSwitched

	EventPlayerHealthChanged
	EventPlayerArmourChanged
	EventPlayerAlivenessChanged

	EventBombPlanted
	EventBombDefused
	EventBombExploded
)

// EnumToEventName KEEP THIS LIST UP TO DATE WITH THE ABOVE ENUMS AS IT IS USED IN DOCS
var EnumToEventName = map[GameEvent]string{
	EventHeartbeat:                  "HeartBeat",
	EventPlayerPaused:               "PlayerPaused",
	EventPlayerPlaying:              "PlayerPlaying",
	EventPlayerInTextInput:          "PlayerInTextInput",
	EventPlayerWeaponUse:            "PlayerWeaponUse",
	EventPlayerWeaponReloadStarted:  "PlayerWeaponReloadStarted",
	EventPlayerWeaponReloadFinished: "PlayerWeaponReloadFinished",
	EventPlayerWeaponChanged:        "PlayerWeaponChanged",
	EventPlayerWeaponAdded:          "PlayerWeaponAdded",
	EventPlayerWeaponRemoved:        "PlayerWeaponRemoved",
	EventPlayerActiveWeaponSwitched: "PlayerActiveWeaponSwitched",
	EventPlayerHealthChanged:        "EventPlayerHealthChanged",
	EventPlayerArmourChanged:        "EventPlayerArmourChanged",
	EventPlayerAlivenessChanged:     "EventPlayerAlivenessChanged",
	EventBombPlanted:                "EventBombPlanted",
	EventBombDefused:                "EventBombDefused",
	EventBombExploded:               "EventBombExploded",
}
var EventNameToEnum map[string]GameEvent = func() map[string]GameEvent {
	m := make(map[string]GameEvent)
	for k, v := range EnumToEventName {
		m[v] = k
	}
	return m
}()

func (e GameEvent) String() string {
	if name, exists := EnumToEventName[e]; exists {
		return name
	}
	log.Error().Msgf("Unknown GameEvent %v", strconv.Itoa(int(e)))
	return "Unknown Event"
}
