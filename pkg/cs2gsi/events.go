package cs2gsi

import "github.com/rs/zerolog/log"

type GameEventDetails struct {
	EventType GameEvent
	EventPlayerWeaponAmmoChange
}
type GameEvent int

const (
	EventHeartbeat GameEvent = iota
	EventPlayerWeaponUse
	EventPlayerWeaponReload
)

type gameEventHandlerCallback func(GSIEvent, GameEventDetails)

var gameEventHandlers = make(map[GameEvent][]gameEventHandlerCallback)

func RegisterEventHandler(event GameEvent, handler gameEventHandlerCallback) {
	gameEventHandlers[event] = append(gameEventHandlers[event], handler)
}

type EventPlayerWeaponAmmoChange struct {
	OldAmmoAmount int
	NewAmmoAmount int
}

type GameEventChecker func(gsiEvent GSIEvent) *GameEventDetails

var eventCheckers = []GameEventChecker{
	checkEventHeartbeat,
	checkEventWeaponAmmoChange,
}

func findEvents(gsiEvent GSIEvent) []GameEventDetails {
	events := make([]GameEventDetails, 0)
	for _, checker := range eventCheckers {
		res := checker(gsiEvent)
		if res != nil {
			events = append(events, *res)
		}
	}
	return events
}

func checkEventHeartbeat(gsiEvent GSIEvent) *GameEventDetails {
	if gsiEvent.Previous == nil {
		return &GameEventDetails{EventType: EventHeartbeat}
	}
	return nil
}

func checkEventWeaponAmmoChange(gsiEvent GSIEvent) *GameEventDetails {
	if gsiEvent.Previous == nil ||
		gsiEvent.Previous.Player == nil ||
		len(gsiEvent.Previous.Player.Weapons) == 0 {

		return nil
	}

	keyForChangedWeapon := ""
	for weaponKey, weapon := range gsiEvent.Previous.Player.Weapons {
		if weapon.AmmoClip != nil {
			keyForChangedWeapon = weaponKey
		}
	}

	if keyForChangedWeapon == "" {
		return nil
	}

	ammoEvent := EventPlayerWeaponAmmoChange{
		*gsiEvent.Previous.Player.Weapons[keyForChangedWeapon].AmmoClip,
		*gsiEvent.Player.Weapons[keyForChangedWeapon].AmmoClip,
	}

	if ammoEvent.NewAmmoAmount >= ammoEvent.OldAmmoAmount {
		return &GameEventDetails{
			EventType:                   EventPlayerWeaponReload,
			EventPlayerWeaponAmmoChange: ammoEvent,
		}
	} else if ammoEvent.OldAmmoAmount > ammoEvent.NewAmmoAmount {
		return &GameEventDetails{
			EventType:                   EventPlayerWeaponUse,
			EventPlayerWeaponAmmoChange: ammoEvent,
		}

	} else {
		log.Error().
			Int("OldAmmo", ammoEvent.OldAmmoAmount).
			Int("NewAmmo", ammoEvent.NewAmmoAmount).
			Msg("Unknown Ammo change event")
		return nil
	}
}
