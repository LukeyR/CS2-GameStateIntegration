package main

import (
	"fmt"
	"time"

	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"
)

func main() {
	cs2gsi.RegisterGlobalHandler(func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		fmt.Printf("(%v) %v %v\n",
			time.Now().Format("2006-01-02 15:04:05.000"),
			events.EnumToEventName[gameEvent.EventType],
			gsiEvent.GetOriginalRequestFlat(),
		)
	})
	cs2gsi.RegisterNonEventHandler(func(gsiEvent *structs.GSIEvent) {
		fmt.Printf("(%v) #N/A %v\n",
			time.Now().Format("2006-01-02 15:04:05.000"),
			gsiEvent.GetOriginalRequestFlat(),
		)
	})

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

	cs2gsi.RegisterEventHandler(events.EventBombPlanted, func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		fmt.Println("The Bomb has been planted")
	})

	cs2gsi.RegisterEventHandler(events.EventBombExploded, func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		fmt.Println("The Bomb has exploded")
	})

	cs2gsi.RegisterEventHandler(events.EventBombDefused, func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		fmt.Println("The Bomb has been defused")
	})

	cs2gsi.RegisterEventHandler(events.EventPlayerHealthChanged, func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		oldHealth := gameEvent.EventPlayerHealthChanged.Old
		newHealth := gameEvent.EventPlayerHealthChanged.New

		direction := ""
		if oldHealth < newHealth {
			direction = "gained"
		} else {
			direction = "lost"
		}

		fmt.Printf(
			"%v just %v health. (%v => %v) \n",
			gsiEvent.Player.Name,
			direction,
			oldHealth,
			newHealth,
		)
	})

	cs2gsi.RegisterEventHandler(events.EventPlayerArmourChanged, func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		oldArmour := gameEvent.EventPlayerArmourChanged.Old
		newArmour := gameEvent.EventPlayerArmourChanged.New

		direction := ""
		if oldArmour < newArmour {
			direction = "gained"
		} else {
			direction = "lost"
		}

		fmt.Printf(
			"%v just %v Armour. (%v => %v) \n",
			gsiEvent.Player.Name,
			direction,
			oldArmour,
			newArmour,
		)
	})

	cs2gsi.RegisterEventHandler(events.EventPlayerAlivenessChanged, func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		//fmt.Println(gsiEvent.GetOriginalRequestFlat())
		oldHealth := gameEvent.EventPlayerHealthChanged.Old
		newHealth := gameEvent.EventPlayerHealthChanged.New

		direction := ""
		if oldHealth > newHealth {
			direction = "died"
		} else {
			direction = "respawned"
		}

		fmt.Printf(
			"%v just %v. (%v => %v) \n",
			gsiEvent.Player.Name,
			direction,
			oldHealth,
			newHealth,
		)
	})

	err := cs2gsi.StartupAndServe(":8000")
	if err != nil {
		return
	}

}
