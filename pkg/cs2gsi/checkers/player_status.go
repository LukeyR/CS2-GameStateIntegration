package checkers

import (
	"errors"

	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"
)

/*
getPlayerHealthData
return old health, new health, error
*/
func getPlayerHealthDelta(gsiEvent *structs.GSIEvent, checkArmour bool) (int, int, error) {
	previousHealthAvailable := !(gsiEvent.Previous == nil ||
		gsiEvent.Previous.Player == nil ||
		(!checkArmour && gsiEvent.Previous.Player.State.Health == nil) ||
		(checkArmour && gsiEvent.Previous.Player.State.Armor == nil))

	if !previousHealthAvailable {
		return 0, 0, errors.New("previous health/armour data not available, player did not take damage")
	}

	if !checkArmour {
		return *gsiEvent.Previous.Player.State.Health, *gsiEvent.Player.State.Health, nil
	} else {
		return *gsiEvent.Previous.Player.State.Armor, *gsiEvent.Player.State.Armor, nil
	}
}

func CheckEventPlayerHealthChanged(gsiEvent *structs.GSIEvent) *events.GameEventDetails {
	oldHealth, newHealth, err := getPlayerHealthDelta(gsiEvent, false)
	if err != nil {
		return nil
	}

	if newHealth == 0 || oldHealth == 0 {
		return nil
	}

	return &events.GameEventDetails{
		EventType: events.EventPlayerHealthChanged,
		EventPlayerHealthChanged: &events.HealthChangedEventDetails{
			Old: oldHealth,
			New: newHealth,
		},
	}
}

func CheckEventPlayerArmourChanged(gsiEvent *structs.GSIEvent) *events.GameEventDetails {
	oldHealth, newHealth, err := getPlayerHealthDelta(gsiEvent, true)
	if err != nil {
		return nil
	}

	if newHealth == 0 || oldHealth == 0 {
		return nil
	}

	return &events.GameEventDetails{
		EventType: events.EventPlayerArmourChanged,
		EventPlayerArmourChanged: &events.ArmourChangedEventDetails{
			Old: oldHealth,
			New: newHealth,
		},
	}
}

func CheckEventPlayerAliveStatusChanged(gsiEvent *structs.GSIEvent) *events.GameEventDetails {
	oldHealth, newHealth, err := getPlayerHealthDelta(gsiEvent, false)
	if err != nil {
		return nil
	}

	if newHealth != 0 && oldHealth != 0 {
		return nil
	}

	return &events.GameEventDetails{
		EventType: events.EventPlayerAlivenessChanged,
		EventPlayerHealthChanged: &events.HealthChangedEventDetails{
			Old: oldHealth,
			New: newHealth,
		},
	}
}
