package checkers

import (
	"slices"

	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"

	"github.com/rs/zerolog/log"
)

func CheckEventWeaponsChanged(gsiEvent *structs.GSIEvent) *events.GameEventDetails {
	previousWeaponsPresent := !(gsiEvent.Previous == nil || gsiEvent.Previous.Player == nil || len(gsiEvent.Previous.Player.Weapons) == 0)
	newWeaponsPresent := !(gsiEvent.Added == nil || gsiEvent.Added.Player == nil || len(gsiEvent.Added.Player.Weapons) == 0)
	playerDead := *gsiEvent.Player.State.Health == 0

	if (!previousWeaponsPresent && !newWeaponsPresent) || playerDead {
		return nil
	}

	if previousWeaponsPresent && !newWeaponsPresent {
		changedWeapons := gsiEvent.Previous.Player.Weapons
		if len(changedWeapons) == 1 {
			// Player has reloading, or switched weapons of same position for each
			// other (e.g. swapped ak for awp, or glock for p250, etc. etc.)
			keyForChangedWeapon := ""
			var changedWeaponData structs.Weapon
			for weaponKey, weapon := range gsiEvent.Previous.Player.Weapons {
				keyForChangedWeapon = weaponKey
				changedWeaponData = *weapon
			}
			newWeaponData := gsiEvent.Player.Weapons[keyForChangedWeapon]
			// keyForChangedWeapon is the weapon slot that has changed; lets figure out what has changed

			// Check if the name of the old weapon to the new on has changed - if so, they have changed weapons
			if changedWeaponData.Name != "" &&
				changedWeaponData.Name != newWeaponData.Name {
				return &events.GameEventDetails{
					EventType: events.EventPlayerWeaponChanged,
					EventPlayerWeaponChange: &events.WeaponChangeEventDetails{
						OldWeaponKey:  keyForChangedWeapon,
						NewWeaponKey:  keyForChangedWeapon,
						OldWeaponName: changedWeaponData.Name,
						NewWeaponName: newWeaponData.Name,
					},
				}
			}

			if newWeaponData.State == structs.WeaponStateReloading {
				return &events.GameEventDetails{
					EventType: events.EventPlayerWeaponReloadStarted,
					EventPlayerWeaponReload: &events.WeaponReloadStartedEventDetails{
						WeaponKey: keyForChangedWeapon,
					},
				}
			}

			if changedWeaponData.State == structs.WeaponStateReloading {
				return &events.GameEventDetails{
					EventType: events.EventPlayerWeaponReloadFinished,
					EventPlayerWeaponReload: &events.WeaponReloadStartedEventDetails{
						WeaponKey: keyForChangedWeapon,
					},
				}
			}

			/* Add handling for grenades
			 {
				"map": {
					"round_wins": {
						"1": "ct_win_time",
						"2": "ct_win_time",
						"3": "ct_win_time",
						"4": "ct_win_time",
						"5": "t_win_bomb"
					},
					"mode": "competitive",
					"name": "de_dust2",
					"phase": "live",
					"round": 5,
					"team_ct": {
						"score": 4,
						"consecutive_round_losses": 1,
						"timeouts_remaining": 1,
						"matches_won_this_series": 0
					},
					"team_t": {
						"score": 1,
						"consecutive_round_losses": 3,
						"timeouts_remaining": 1,
						"matches_won_this_series": 0
					},
					"num_matches_to_win_series": 0
				},
				"player": {
					"steamid": "76561198117545744",
					"name": "Phoenix",
					"observer_slot": 0,
					"team": "T",
					"activity": "playing",
					"match_stats": {
						"kills": 0,
						"assists": 0,
						"deaths": 0,
						"mvps": 1,
						"score": 4
					},
					"state": {
						"health": 99,
						"armor": 100,
						"helmet": true,
						"flashed": 0,
						"smoked": 0,
						"burning": 0,
						"money": 27700,
						"round_kills": 0,
						"round_killhs": 0,
						"equip_value": 4700
					},
					"weapons": {
						"weapon_0": {
							"name": "weapon_knife_t",
							"paintkit": "default",
							"type": "Knife",
							"state": "holstered"
						},
						"weapon_1": {
							"name": "weapon_c4",
							"paintkit": "default",
							"type": "C4",
							"state": "holstered"
						},
						"weapon_2": {
							"name": "weapon_taser",
							"paintkit": "default",
							"ammo_clip": 1,
							"ammo_clip_max": 1,
							"ammo_reserve": 0,
							"state": "holstered"
						},
						"weapon_3": {
							"name": "weapon_hegrenade",
							"paintkit": "default",
							"type": "Grenade",
							"ammo_reserve": 0,
							"state": "active"
						},
						"weapon_4": {
							"name": "weapon_flashbang",
							"paintkit": "default",
							"type": "Grenade",
							"ammo_reserve": 1,
							"state": "holstered"
						},
						"weapon_5": {
							"name": "weapon_cz75a",
							"paintkit": "default",
							"type": "Pistol",
							"ammo_clip": 12,
							"ammo_clip_max": 12,
							"ammo_reserve": 12,
							"state": "holstered"
						},
						"weapon_6": {
							"name": "weapon_galilar",
							"paintkit": "default",
							"type": "Rifle",
							"ammo_clip": 35,
							"ammo_clip_max": 35,
							"ammo_reserve": 90,
							"state": "holstered"
						}
					}
				},
				"provider": {
					"name": "Counter-Strike: Global Offensive",
					"appid": 730,
					"version": 14037,
					"steamid": "76561198117545744",
					"timestamp": 1729088255
				},
				"round": {
					"phase": "live"
				},
				"previously": {
					"player": {
						"weapons": {
							"weapon_3": {
								"ammo_reserve": 1
							}
						}
					}
				}
			}
			*/

			// Weapon name same, their ammo must have changed. Did they shoot, or reload?
			ammoEvent := events.AmmoChangeEventDetails{
				OldAmmoAmount: *changedWeaponData.AmmoClip,
				NewAmmoAmount: *newWeaponData.AmmoClip,
			}

			if ammoEvent.NewAmmoAmount >= ammoEvent.OldAmmoAmount {
				return &events.GameEventDetails{
					EventType: events.EventPlayerWeaponReloadFinished,
					EventPlayerWeaponReload: &events.WeaponReloadStartedEventDetails{
						WeaponKey: keyForChangedWeapon,
					},
				}
			} else if ammoEvent.OldAmmoAmount > ammoEvent.NewAmmoAmount {
				return &events.GameEventDetails{
					EventType:                   events.EventPlayerWeaponUse,
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
			weaponsDropped := make(structs.WeaponCollection)
			for weaponKey := range gsiEvent.Previous.Player.Weapons {
				if _, exists := gsiEvent.Player.Weapons[weaponKey]; !exists {
					weaponsDropped[weaponKey] = gsiEvent.Previous.Player.Weapons[weaponKey]
				}
			}
			if len(weaponsDropped) > 0 {
				return &events.GameEventDetails{
					EventType:                       events.EventPlayerWeaponRemoved,
					EventPlayerWeaponAddedOrRemoved: &events.WeaponAddedOrRemovedEventDetails{weaponsDropped},
				}
			}

			// No weapons dropped
			var oldWeaponHeldKey string
			var newWeaponHeldKey string
			for weaponKey, weapon := range gsiEvent.Previous.Player.Weapons {
				if slices.Contains([]structs.WeaponState{structs.WeaponStateActive, structs.WeaponStateReloading}, weapon.State) {
					oldWeaponHeldKey = weaponKey
				}
			}
			for weaponKey, weapon := range gsiEvent.Player.Weapons {
				if weapon.State == structs.WeaponStateActive {
					newWeaponHeldKey = weaponKey
				}
			}

			return &events.GameEventDetails{
				EventType: events.EventPlayerActiveWeaponSwitched,
				EventPlayerActiveWeaponChange: &events.ActiveWeaponChangeEventDetails{
					OldWeaponKey: oldWeaponHeldKey,
					NewWeaponKey: newWeaponHeldKey,
				},
			}
		}
	} else {
		weaponsAdded := make(structs.WeaponCollection)
		for weaponKey := range gsiEvent.Added.Player.Weapons {
			weaponsAdded[weaponKey] = gsiEvent.Player.Weapons[weaponKey]
		}
		return &events.GameEventDetails{
			EventType:                       events.EventPlayerWeaponAdded,
			EventPlayerWeaponAddedOrRemoved: &events.WeaponAddedOrRemovedEventDetails{Weapons: weaponsAdded},
		}
	}

	originalRequest := gsiEvent.GetOriginalRequestFlat()
	log.Warn().
		Str("event", originalRequest).
		Msg("Unknown Weapon change event")
	return nil
}
