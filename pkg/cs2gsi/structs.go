package cs2gsi

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

type Team struct {
	Score                  int `json:"score"`
	ConsecutiveRoundLosses int `json:"consecutive_round_losses"`
	TimeoutsRemaining      int `json:"timeouts_remaining"`
	MatchesWonThisSeries   int `json:"matches_won_this_series"`
}

func (team *Team) String() string {
	return fmt.Sprintf(
		"Team{ \n\t\t\t\t"+
			"Score: %v, \n\t\t\t\t"+
			"ConsecutiveRoundLosses: %v, \n\t\t\t\t"+
			"TimeoutsRemaining: %v,\n\t\t\t\t"+
			"MatchesWonThisSeries: %v, \n\t\t\t"+
			"}",
		team.Score,
		team.ConsecutiveRoundLosses,
		team.TimeoutsRemaining,
		team.MatchesWonThisSeries,
	)
}

type CSMap struct {
	RoundWins             map[string]string `json:"round_wins"`
	Mode                  string            `json:"mode"`
	Name                  string            `json:"name"`
	Phase                 string            `json:"phase"`
	Round                 int               `json:"Round"`
	TeamCt                *Team             `json:"team_ct"`
	TeamT                 *Team             `json:"team_t"`
	NumMatchesToWinSeries int               `json:"num_matches_to_win_series"`
}

func (csmap *CSMap) String() string {
	return fmt.Sprintf(
		"CSMap{ \n\t\t\t"+
			"RoundWins: %v, \n\t\t\t"+
			"Mode: %v, \n\t\t\t"+
			"Name: %v,\n\t\t\t"+
			"Phase: %v, \n\t\t\t"+
			"Round: %v \n\t\t\t"+
			"TeamCt: %v \n\t\t\t"+
			"TeamT: %v, \n\t\t\t"+
			"NumMatchesToWinSeries: %v \n\t"+
			"}",
		csmap.RoundWins,
		csmap.Mode,
		csmap.Name,
		csmap.Phase,
		csmap.Round,
		csmap.TeamCt,
		csmap.TeamT,
		csmap.NumMatchesToWinSeries,
	)
}

type WeaponType string

const (
	WesponTypePistol  WeaponType = "Pistol"
	WeaponTypeGrenade WeaponType = "Grenade"
)

type WeaponState string

const (
	WesponStateHolstered WeaponState = "holstered"
	WeaponStateActive    WeaponState = "active"
	WeaponStateReloading WeaponState = "reloading"
)

type Weapon struct {
	Name        string      `json:"name"`
	Paintkit    string      `json:"paintkit"`
	Type        WeaponType  `json:"type"`
	AmmoClip    *int        `json:"ammo_clip"`
	AmmoClipMax *int        `json:"ammo_clip_max"`
	AmmoReserve *int        `json:"ammo_reserve"`
	State       WeaponState `json:"state"`
}

func (weapon *Weapon) String() string {
	return fmt.Sprintf(
		"Weapon{ \n\t\t\t"+
			"Name: %v, \n\t\t\t"+
			"Paintkit: %v, \n\t\t\t"+
			"Type: %v,\n\t\t\t"+
			"AmmoClip: %v, \n\t\t\t"+
			"AmmoClipMax: %v \n\t\t\t"+
			"AmmoReserve: %v, \n\t\t\t"+
			"State: %v \n\t"+
			"}",
		weapon.Name,
		weapon.Paintkit,
		weapon.Type,
		weapon.AmmoClip,
		weapon.AmmoClipMax,
		weapon.AmmoReserve,
		weapon.State,
	)
}

type WeaponCollection map[string]*Weapon

func (wCollection *WeaponCollection) String() string {
	toReturn := ""
	for weaponKey, weaponData := range *wCollection {
		toReturn += fmt.Sprintf("%v: %v,\n", weaponKey, weaponData)
	}
	return toReturn
}

type PlayerActivity string

const (
	PlayerActivityPaused      PlayerActivity = "menu"
	PlayerActivityPlaying     PlayerActivity = "playing"
	PlayerActivityInTextInput PlayerActivity = "textinput"
)

type Player struct {
	Steamid      string         `json:"steamid"`
	Name         string         `json:"name"`
	ObserverSlot int            `json:"observer_slot"`
	Team         string         `json:"Team"`
	Activity     PlayerActivity `json:"activity"`
	MatchStats   struct {
		Kills   int `json:"kills"`
		Assists int `json:"assists"`
		Deaths  int `json:"deaths"`
		Mvps    int `json:"mvps"`
		Score   int `json:"score"`
	} `json:"match_stats"`

	State struct {
		Health      int  `json:"health"`
		Armor       int  `json:"armor"`
		Helmet      bool `json:"helmet"`
		Flashed     int  `json:"flashed"`
		Smoked      int  `json:"smoked"`
		Burning     int  `json:"burning"`
		Money       int  `json:"money"`
		RoundKills  int  `json:"round_kills"`
		RoundKillHS int  `json:"round_killhs"`
		EquipValue  int  `json:"equip_value"`
	} `json:"state"`
	Weapons WeaponCollection `json:"Weapons"`
}

func (player *Player) String() string {
	return fmt.Sprintf(
		"Player{\n\t\t"+
			"SteamID: %v, \n\t\t"+
			"Name: %v, \n\t\t"+
			"ObserverSlot: %v, \n\t\t"+
			"Team: %v, \n\t\t"+
			"Activity: %v, \n\t\t"+
			"MatchStats: {\n\t\t\t"+
			"Kills: %v, \n\t\t\t"+
			"Assists: %v, \n\t\t\t"+
			"Deaths: %v,\n\t\t\t"+
			"Mvps: %v,\n\t\t\t"+
			"Score: %v \n\t\t"+
			"}\n\t\t"+
			"State: {\n\t\t\t"+
			"Health: %v, \n\t\t\t"+
			"Armor: %v, \n\t\t\t"+
			"Helmet: %v,\n\t\t\t"+
			"Flashed: %v,\n\t\t\t"+
			"Smoked: %v \n\t\t\t"+
			"Burning: %v, \n\t\t\t"+
			"Money: %v, \n\t\t\t"+
			"RoundKills: %v,\n\t\t\t"+
			"RoundKillHS: %v,\n\t\t\t"+
			"EquipValue: %v \n\t\t"+
			"}\n\t"+
			"Weapons(%v): {\n\t\t"+
			"%v\n\t\t"+
			"}\n\t"+
			"}",
		player.Steamid,
		player.Name,
		player.ObserverSlot,
		player.Team,
		player.Activity,
		player.MatchStats.Kills,
		player.MatchStats.Assists,
		player.MatchStats.Deaths,
		player.MatchStats.Mvps,
		player.MatchStats.Score,
		player.State.Health,
		player.State.Armor,
		player.State.Helmet,
		player.State.Flashed,
		player.State.Smoked,
		player.State.Burning,
		player.State.Money,
		player.State.RoundKills,
		player.State.RoundKillHS,
		player.State.EquipValue,
		len(player.Weapons),
		player.Weapons,
	)
}

type Provider struct {
	Name      string `json:"name"`
	Appid     int    `json:"appid"`
	Version   int    `json:"version"`
	SteamID   string `json:"steamid"`
	Timestamp int    `json:"timestamp"`
}

func (provider *Provider) String() string {
	return fmt.Sprintf(
		"Provider{ \n\t\t"+
			"Name: %v, \n\t\t"+
			"Appid: %v, \n\t\t"+
			"Version: %v,\n\t\t"+
			"SteamID: %v, \n\t\t"+
			"Timestamp: %v \n\t"+
			"}",
		provider.Name,
		provider.Appid,
		provider.Version,
		provider.SteamID,
		provider.Timestamp,
	)
}

type Round struct {
	Phase   string `json:"phase"`
	WinTeam string `json:"win_team"`
	Bomb    string `json:"bomb"`
}

func (round *Round) String() string {
	return fmt.Sprintf(
		"Round{ \n\t\t"+
			"Phase: %v, \n\t\t"+
			"WinTeam: %v, \n\t\t"+
			"Bomb: %v \n\t"+
			"}",
		round.Phase,
		round.WinTeam,
		round.Bomb,
	)
}

type GSIEventAdded struct {
	Player *struct {
		Weapons map[string]bool `json:"weapons"`
	} `json:"player"`
}

type GSIEvent struct {
	CSMap    *CSMap         `json:"map"`
	Player   *Player        `json:"Player"`
	Provider *Provider      `json:"Provider"`
	Round    *Round         `json:"Round"`
	Previous *GSIEvent      `json:"Previously"`
	Added    *GSIEventAdded `json:"Added"`

	OriginalData string `json:"-"`
}

func (gsiEvent *GSIEvent) String() string {
	res, err := json.Marshal(gsiEvent)
	if err != nil {
		log.Error().Err(err).Msg("Error Marshalling GSIEvent")
	}
	return string(res)
}

func NewGSIEvent(requestBody string) (*GSIEvent, error) {
	newEvent := &GSIEvent{}
	err := json.Unmarshal([]byte(requestBody), newEvent)
	if err != nil {
		return nil, err
	}
	newEvent.OriginalData = requestBody
	return newEvent, nil
}
func (gsiEvent *GSIEvent) GetOriginalRequestFlat() string {
	returnBytes := &bytes.Buffer{}
	err := json.Compact(returnBytes, []byte(gsiEvent.OriginalData))
	if err != nil {
		log.Warn().Err(err).Msg("Error Compacting GSIEvent")
		return gsiEvent.OriginalData
	}
	return returnBytes.String()
}
