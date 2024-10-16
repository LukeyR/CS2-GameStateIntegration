package cs2gsi

import (
	"slices"

	"github.com/rs/zerolog/log"
)

/*
	Main struct returned is GameEventDetails.
	Alongside the event type raised, it can (anonymously) contain extra data about the event
*/

type EventPlayerWeaponAmmoChange_t struct {
	OldAmmoAmount int
	NewAmmoAmount int
}

type EventPlayerWeaponReloadStarted_t struct {
	WeaponKey string
}

type EventPlayerWeaponChange_t struct {
	OldWeaponKey  string
	NewWeaponKey  string
	OldWeaponName string
	NewWeaponName string
}

type EventPlayerActiveWeaponChange_t struct {
	OldWeaponKey string
	NewWeaponKey string
}

type EventPlayerWeaponAddedOrRemoved_t struct {
	Weapons WeaponCollection
}

type GameEventDetails struct {
	EventType                       GameEvent
	EventPlayerWeaponAmmoChange     *EventPlayerWeaponAmmoChange_t
	EventPlayerWeaponChange         *EventPlayerWeaponChange_t
	EventPlayerWeaponAddedOrRemoved *EventPlayerWeaponAddedOrRemoved_t
	EventPlayerActiveWeaponChange   *EventPlayerActiveWeaponChange_t
	EventPlayerWeaponReload         *EventPlayerWeaponReloadStarted_t
}

/*
All possible event that can be raised, as well as functions needed to check for them
*/
type GameEvent int

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

type GameEventChecker func(gsiEvent *GSIEvent) *GameEventDetails

var eventCheckers = []GameEventChecker{
	checkEventHeartbeat,
	checkEventWeaponsChanged,
	checkEventPlayerActivityChanged,
}

/*
Function for subscribing to GameEvents
*/
type gameEventHandlerCallback func(*GSIEvent, GameEventDetails)

var gameEventHandlers = make(map[GameEvent][]gameEventHandlerCallback)

func RegisterEventHandler(event GameEvent, handler gameEventHandlerCallback) {
	gameEventHandlers[event] = append(gameEventHandlers[event], handler)
}

func findEvents(gsiEvent *GSIEvent) []GameEventDetails {
	events := make([]GameEventDetails, 0)
	for _, checker := range eventCheckers {
		res := checker(gsiEvent)
		if res != nil {
			events = append(events, *res)
		}
	}
	return events
}

/*
Event Checking functions. Should have signature `GameEventChecker`
*/
func checkEventHeartbeat(gsiEvent *GSIEvent) *GameEventDetails {
	if gsiEvent.Previous == nil {
		return &GameEventDetails{EventType: EventHeartbeat}
	}
	return nil
}

func checkEventPlayerActivityChanged(gsiEvent *GSIEvent) *GameEventDetails {
	if gsiEvent.Previous == nil || gsiEvent.Previous.Player == nil {
		return nil
	}

	if gsiEvent.Previous.Player.Activity != "" {
		switch gsiEvent.Player.Activity {
		case PlayerActivityPlaying:
			return &GameEventDetails{EventType: EventPlayerPlaying}
		case PlayerActivityPaused:
			return &GameEventDetails{EventType: EventPlayerPaused}
		case PlayerActivityInTextInput:
			return &GameEventDetails{EventType: EventPlayerInTextInput}
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

func checkEventWeaponsChanged(gsiEvent *GSIEvent) *GameEventDetails {
	previousWeaponsPresent := !(gsiEvent.Previous == nil || gsiEvent.Previous.Player == nil || len(gsiEvent.Previous.Player.Weapons) == 0)
	newWeaponsPresent := !(gsiEvent.Previous == nil || gsiEvent.Previous.Player == nil || len(gsiEvent.Previous.Player.Weapons) == 0)

	if !previousWeaponsPresent && !newWeaponsPresent {
		return nil
	}

	if previousWeaponsPresent {
		changedWeapons := gsiEvent.Previous.Player.Weapons
		if len(changedWeapons) == 1 {
			// Player has reloading, or switched weapons of same position for each
			// other (e.g. swapped ak for awp, or glock for p250, etc. etc.)
			keyForChangedWeapon := ""
			var changedWeaponData Weapon
			for weaponKey, weapon := range gsiEvent.Previous.Player.Weapons {
				keyForChangedWeapon = weaponKey
				changedWeaponData = *weapon
			}
			newWeaponData := gsiEvent.Player.Weapons[keyForChangedWeapon]
			// keyForChangedWeapon is the weapon slot that has changed; lets figure out what has changed

			// Check if the name of the old weapon to the new on has changed - if so, they have changed weapons
			if changedWeaponData.Name != "" &&
				changedWeaponData.Name != newWeaponData.Name {
				return &GameEventDetails{
					EventType: EventPlayerWeaponChanged,
					EventPlayerWeaponChange: &EventPlayerWeaponChange_t{
						OldWeaponKey:  keyForChangedWeapon,
						NewWeaponKey:  keyForChangedWeapon,
						OldWeaponName: changedWeaponData.Name,
						NewWeaponName: newWeaponData.Name,
					},
				}
			}

			if newWeaponData.State == WeaponStateReloading {
				return &GameEventDetails{
					EventType: EventPlayerWeaponReloadStarted,
					EventPlayerWeaponReload: &EventPlayerWeaponReloadStarted_t{
						WeaponKey: keyForChangedWeapon,
					},
				}
			}

			if changedWeaponData.State == WeaponStateReloading {
				return &GameEventDetails{
					EventType: EventPlayerWeaponReloadFinished,
					EventPlayerWeaponReload: &EventPlayerWeaponReloadStarted_t{
						WeaponKey: keyForChangedWeapon,
					},
				}
			}

			// Add handling for grenades
			// {
			//	"map": {
			//		"round_wins": {
			//			"1": "ct_win_time",
			//			"2": "ct_win_time",
			//			"3": "ct_win_time",
			//			"4": "ct_win_time",
			//			"5": "t_win_bomb"
			//		},
			//		"mode": "competitive",
			//		"name": "de_dust2",
			//		"phase": "live",
			//		"round": 5,
			//		"team_ct": {
			//			"score": 4,
			//			"consecutive_round_losses": 1,
			//			"timeouts_remaining": 1,
			//			"matches_won_this_series": 0
			//		},
			//		"team_t": {
			//			"score": 1,
			//			"consecutive_round_losses": 3,
			//			"timeouts_remaining": 1,
			//			"matches_won_this_series": 0
			//		},
			//		"num_matches_to_win_series": 0
			//	},
			//	"player": {
			//		"steamid": "76561198117545744",
			//		"name": "Phoenix",
			//		"observer_slot": 0,
			//		"team": "T",
			//		"activity": "playing",
			//		"match_stats": {
			//			"kills": 0,
			//			"assists": 0,
			//			"deaths": 0,
			//			"mvps": 1,
			//			"score": 4
			//		},
			//		"state": {
			//			"health": 99,
			//			"armor": 100,
			//			"helmet": true,
			//			"flashed": 0,
			//			"smoked": 0,
			//			"burning": 0,
			//			"money": 27700,
			//			"round_kills": 0,
			//			"round_killhs": 0,
			//			"equip_value": 4700
			//		},
			//		"weapons": {
			//			"weapon_0": {
			//				"name": "weapon_knife_t",
			//				"paintkit": "default",
			//				"type": "Knife",
			//				"state": "holstered"
			//			},
			//			"weapon_1": {
			//				"name": "weapon_c4",
			//				"paintkit": "default",
			//				"type": "C4",
			//				"state": "holstered"
			//			},
			//			"weapon_2": {
			//				"name": "weapon_taser",
			//				"paintkit": "default",
			//				"ammo_clip": 1,
			//				"ammo_clip_max": 1,
			//				"ammo_reserve": 0,
			//				"state": "holstered"
			//			},
			//			"weapon_3": {
			//				"name": "weapon_hegrenade",
			//				"paintkit": "default",
			//				"type": "Grenade",
			//				"ammo_reserve": 0,
			//				"state": "active"
			//			},
			//			"weapon_4": {
			//				"name": "weapon_flashbang",
			//				"paintkit": "default",
			//				"type": "Grenade",
			//				"ammo_reserve": 1,
			//				"state": "holstered"
			//			},
			//			"weapon_5": {
			//				"name": "weapon_cz75a",
			//				"paintkit": "default",
			//				"type": "Pistol",
			//				"ammo_clip": 12,
			//				"ammo_clip_max": 12,
			//				"ammo_reserve": 12,
			//				"state": "holstered"
			//			},
			//			"weapon_6": {
			//				"name": "weapon_galilar",
			//				"paintkit": "default",
			//				"type": "Rifle",
			//				"ammo_clip": 35,
			//				"ammo_clip_max": 35,
			//				"ammo_reserve": 90,
			//				"state": "holstered"
			//			}
			//		}
			//	},
			//	"provider": {
			//		"name": "Counter-Strike: Global Offensive",
			//		"appid": 730,
			//		"version": 14037,
			//		"steamid": "76561198117545744",
			//		"timestamp": 1729088255
			//	},
			//	"round": {
			//		"phase": "live"
			//	},
			//	"previously": {
			//		"player": {
			//			"weapons": {
			//				"weapon_3": {
			//					"ammo_reserve": 1
			//				}
			//			}
			//		}
			//	}
			//}

			// Weapon name same, their ammo must have changed. Did they shoot, or reload?
			ammoEvent := EventPlayerWeaponAmmoChange_t{
				*changedWeaponData.AmmoClip,
				*newWeaponData.AmmoClip,
			}

			if ammoEvent.NewAmmoAmount >= ammoEvent.OldAmmoAmount {
				return &GameEventDetails{
					EventType: EventPlayerWeaponReloadFinished,
					EventPlayerWeaponReload: &EventPlayerWeaponReloadStarted_t{
						WeaponKey: keyForChangedWeapon,
					},
				}
			} else if ammoEvent.OldAmmoAmount > ammoEvent.NewAmmoAmount {
				return &GameEventDetails{
					EventType:                   EventPlayerWeaponUse,
					EventPlayerWeaponAmmoChange: &ammoEvent,
				}

			} else {
				log.Error().
					Int("OldAmmo", ammoEvent.OldAmmoAmount).
					Int("NewAmmo", ammoEvent.NewAmmoAmount).
					Msg("Unknown Ammo change event")
				return nil
			}
		} else if len(changedWeapons) >= 2 {
			// First lets check if they dropped any weapons
			weaponsDropped := make(WeaponCollection)
			for weaponKey := range gsiEvent.Previous.Player.Weapons {
				if _, exists := gsiEvent.Player.Weapons[weaponKey]; !exists {
					weaponsDropped[weaponKey] = gsiEvent.Previous.Player.Weapons[weaponKey]
				}
			}
			if len(weaponsDropped) > 0 {
				return &GameEventDetails{
					EventType:                       EventPlayerWeaponRemoved,
					EventPlayerWeaponAddedOrRemoved: &EventPlayerWeaponAddedOrRemoved_t{weaponsDropped},
				}
			}

			// No weapons dropped
			var oldWeaponHeldKey string
			var newWeaponHeldKey string
			for weaponKey, weapon := range gsiEvent.Previous.Player.Weapons {
				if slices.Contains([]WeaponState{WeaponStateActive, WeaponStateReloading}, weapon.State) {
					oldWeaponHeldKey = weaponKey
				}
			}
			for weaponKey, weapon := range gsiEvent.Player.Weapons {
				if weapon.State == WeaponStateActive {
					newWeaponHeldKey = weaponKey
				}
			}

			return &GameEventDetails{
				EventType: EventPlayerActiveWeaponSwitched,
				EventPlayerActiveWeaponChange: &EventPlayerActiveWeaponChange_t{
					OldWeaponKey: oldWeaponHeldKey,
					NewWeaponKey: newWeaponHeldKey,
				},
			}
		}
	} else {
		weaponsAdded := make(WeaponCollection)
		for weaponKey := range gsiEvent.Added.Player.Weapons {
			weaponsAdded[weaponKey] = gsiEvent.Player.Weapons[weaponKey]
		}
		return &GameEventDetails{
			EventType:                       EventPlayerWeaponAdded,
			EventPlayerWeaponAddedOrRemoved: &EventPlayerWeaponAddedOrRemoved_t{Weapons: weaponsAdded},
		}
	}

	originalRequest := gsiEvent.GetOriginalRequestFlat()
	log.Warn().
		Str("event", originalRequest).
		Msg("Unknown Weapon change event")
	return nil
}
