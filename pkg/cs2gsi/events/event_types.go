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
)

var eventNames = map[GameEvent]string{
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
}

func (e GameEvent) String() string {
	if name, exists := eventNames[e]; exists {
		return name
	}
	log.Error().Msgf("Unknown GameEvent %v", strconv.Itoa(int(e)))
	return "Unknown Event"
}
