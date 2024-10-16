package main

import (
	"fmt"

	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi"
)

func main() {
	cs2gsi.RegisterEventHandler(cs2gsi.EventPlayerWeaponUse, func(gsiEvent *cs2gsi.GSIEvent, gameEvent cs2gsi.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		fmt.Printf("%v just shot %v bullets. They have %v bullets left\n", gsiEvent.Player.Name, gameEvent.EventPlayerWeaponAmmoChange.OldAmmoAmount-gameEvent.EventPlayerWeaponAmmoChange.NewAmmoAmount, gameEvent.EventPlayerWeaponAmmoChange.NewAmmoAmount)
	})

	cs2gsi.RegisterEventHandler(cs2gsi.EventPlayerWeaponReloadStarted, func(gsiEvent *cs2gsi.GSIEvent, gameEvent cs2gsi.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		fmt.Printf("%v just started reloading their %v\n", gsiEvent.Player.Name, gsiEvent.Player.Weapons[gameEvent.EventPlayerWeaponReload.WeaponKey].Name)
	})

	cs2gsi.RegisterEventHandler(cs2gsi.EventPlayerWeaponReloadFinished, func(gsiEvent *cs2gsi.GSIEvent, gameEvent cs2gsi.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		fmt.Printf("%v just reloaded. They have %v bullets\n", gsiEvent.Player.Name, *gsiEvent.Player.Weapons[gameEvent.EventPlayerWeaponReload.WeaponKey].AmmoClip)
	})

	cs2gsi.RegisterEventHandler(cs2gsi.EventPlayerWeaponChanged, func(gsiEvent *cs2gsi.GSIEvent, gameEvent cs2gsi.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		fmt.Printf(
			"%v just changed from %v to %v. They have %v bullets\n",
			gsiEvent.Player.Name,
			gameEvent.EventPlayerWeaponChange.OldWeaponName,
			gameEvent.EventPlayerWeaponChange.NewWeaponName,
			*gsiEvent.Player.Weapons[gameEvent.EventPlayerWeaponChange.NewWeaponKey].AmmoClip,
		)
	})

	cs2gsi.RegisterEventHandler(cs2gsi.EventPlayerWeaponAdded, func(gsiEvent *cs2gsi.GSIEvent, gameEvent cs2gsi.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		for _, weapon := range gameEvent.EventPlayerWeaponAddedOrRemoved.Weapons {
			fmt.Printf("%v just picked up %v.\n", gsiEvent.Player.Name, weapon.Name)
		}
	})

	cs2gsi.RegisterEventHandler(cs2gsi.EventPlayerWeaponRemoved, func(gsiEvent *cs2gsi.GSIEvent, gameEvent cs2gsi.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		for _, weapon := range gameEvent.EventPlayerWeaponAddedOrRemoved.Weapons {
			fmt.Printf("%v just dropped their %v.\n", gsiEvent.Player.Name, weapon.Name)
		}
	})

	cs2gsi.RegisterEventHandler(cs2gsi.EventPlayerActiveWeaponSwitched, func(gsiEvent *cs2gsi.GSIEvent, gameEvent cs2gsi.GameEventDetails) {
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
