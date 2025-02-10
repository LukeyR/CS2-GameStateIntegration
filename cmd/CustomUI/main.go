package main

import (
	"fmt"

	"CS2-GameStateIntegration/pkg/cs2gsi"
	"CS2-GameStateIntegration/pkg/cs2gsi/events"
	"CS2-GameStateIntegration/pkg/cs2gsi/structs"
)

func main() {
	cs2gsi.RegisterEventHandler(events.EventPlayerWeaponUse, func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		fmt.Printf("%v just shot %v bullets. They have %v bullets left\n", gsiEvent.Player.Name, gameEvent.EventPlayerWeaponAmmoChange.OldAmmoAmount-gameEvent.EventPlayerWeaponAmmoChange.NewAmmoAmount, gameEvent.EventPlayerWeaponAmmoChange.NewAmmoAmount)
	})

	cs2gsi.RegisterEventHandler(events.EventPlayerWeaponReloadStarted, func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		fmt.Printf("%v just started reloading their %v\n", gsiEvent.Player.Name, gsiEvent.Player.Weapons[gameEvent.EventPlayerWeaponReload.WeaponKey].Name)
	})

	cs2gsi.RegisterEventHandler(events.EventPlayerWeaponReloadFinished, func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		fmt.Printf("%v just reloaded. They have %v bullets\n", gsiEvent.Player.Name, *gsiEvent.Player.Weapons[gameEvent.EventPlayerWeaponReload.WeaponKey].AmmoClip)
	})

	cs2gsi.RegisterEventHandler(events.EventPlayerWeaponChanged, func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		fmt.Printf(
			"%v just changed from %v to %v. They have %v bullets\n",
			gsiEvent.Player.Name,
			gameEvent.EventPlayerWeaponChange.OldWeaponName,
			gameEvent.EventPlayerWeaponChange.NewWeaponName,
			*gsiEvent.Player.Weapons[gameEvent.EventPlayerWeaponChange.NewWeaponKey].AmmoClip,
		)
	})

	cs2gsi.RegisterEventHandler(events.EventPlayerWeaponAdded, func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		for _, weapon := range gameEvent.EventPlayerWeaponAddedOrRemoved.Weapons {
			fmt.Printf("%v just picked up %v.\n", gsiEvent.Player.Name, weapon.Name)
		}
	})

	cs2gsi.RegisterEventHandler(events.EventPlayerWeaponRemoved, func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		for _, weapon := range gameEvent.EventPlayerWeaponAddedOrRemoved.Weapons {
			fmt.Printf("%v just dropped their %v.\n", gsiEvent.Player.Name, weapon.Name)
		}
	})

	cs2gsi.RegisterEventHandler(events.EventPlayerActiveWeaponSwitched, func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		fmt.Printf(
			"%v just holstered their %v and drawn their %v.\n",
			gsiEvent.Player.Name,
			gsiEvent.Player.Weapons[gameEvent.EventPlayerActiveWeaponChange.OldWeaponKey].Name,
			gsiEvent.Player.Weapons[gameEvent.EventPlayerActiveWeaponChange.NewWeaponKey].Name,
		)
	})

	err := cs2gsi.StartupAndServe(":8000")
	if err != nil {
		return
	}

}
