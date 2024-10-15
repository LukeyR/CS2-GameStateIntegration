package main

import (
	"fmt"

	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi"
)

func main() {
	cs2gsi.RegisterEventHandler(cs2gsi.EventPlayerWeaponUse, func(gsiEvent cs2gsi.GSIEvent, gameEvent cs2gsi.GameEventDetails) {
		fmt.Printf("%v just shot %v bullets. They have %v bullets left\n", gsiEvent.Player.Name, gameEvent.OldAmmoAmount-gameEvent.NewAmmoAmount, gameEvent.NewAmmoAmount)
	})

	cs2gsi.RegisterEventHandler(cs2gsi.EventPlayerWeaponReload, func(gsiEvent cs2gsi.GSIEvent, gameEvent cs2gsi.GameEventDetails) {
		fmt.Printf("%v just reloaded. They have %v bullets\n", gsiEvent.Player.Name, gameEvent.NewAmmoAmount)
	})

	err := cs2gsi.StartupAndServe(":8000")
	if err != nil {
		return
	}

}
