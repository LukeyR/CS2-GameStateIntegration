package events

import (
	"CS2-GameStateIntegration/pkg/cs2gsi/structs"
)

type AmmoChangeEventDetails struct {
	OldAmmoAmount int
	NewAmmoAmount int
}

type WeaponReloadStartedEventDetails struct {
	WeaponKey string
}

type WeaponChangeEventDetails struct {
	OldWeaponKey  string
	NewWeaponKey  string
	OldWeaponName string
	NewWeaponName string
}

type ActiveWeaponChangeEventDetails struct {
	OldWeaponKey string
	NewWeaponKey string
}

type WeaponAddedOrRemovedEventDetails struct {
	Weapons structs.WeaponCollection
}

type GameEventDetails struct {
	EventType                       GameEvent
	EventPlayerWeaponAmmoChange     *AmmoChangeEventDetails
	EventPlayerWeaponChange         *WeaponChangeEventDetails
	EventPlayerWeaponAddedOrRemoved *WeaponAddedOrRemovedEventDetails
	EventPlayerActiveWeaponChange   *ActiveWeaponChangeEventDetails
	EventPlayerWeaponReload         *WeaponReloadStartedEventDetails
}
