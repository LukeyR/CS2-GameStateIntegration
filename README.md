# Counter Strike 2 (CS2) Game State Integration (GSI)
A small package to take advantage of [Game State Integration](https://developer.valvesoftware.com/wiki/Counter-Strike:_Global_Offensive_Game_State_Integration) in CS2

## Example Usage
See [cmd](./cmd) for 2 dummy projects.

## Supported Events
Not every unique event in the game is supported yet (although I welcome changes that add more).
The [following ones](./pkg/cs2gsi/events/event_types.go) are supported:
 ```
| *ID* | *Event Name* | *Description* | Notes |
| 0 | EventHeartbeat  | Heartbeat sent from the game to check the server is still accepting |  |
| 1 | EventPlayerPaused | Player paused the game |  |
| 2 | EventPlayerPlaying | Player has started playing again |  |
| 3 | EventPlayerInTextInput | In a text input (seem's to be triggered when entering chat and console) |  |
| 4 | EventPlayerWeaponUse | Player fired their weapon |  |
| 5 | EventPlayerWeaponReloadStarted | Started Reloading their weapon |  |
| 6 | EventPlayerWeaponReloadFinished | Finished Reloading their weapon |  |
| 7 | EventPlayerWeaponChanged | Changing from one type of gun to another (e.g. AK-47 to SG 553) |  |
| 8 | EventPlayerWeaponAdded | Picked up a wespon |  |
| 9 | EventPlayerWeaponRemoved | Dropped a weapon |  |
| 10 | EventPlayerActiveWeaponSwitched | Switched from e.g. primary to secondary |  |
| 11 | EventPlayerHealthChanged | Health increased/decreased |  |
| 12 | EventPlayerArmourChanged | Armour increased/decreased |  |
| 13 | EventPlayerAlivenessChanged | Player died/respawned |  |
| 14 | EventBombPlanted | Bomb was planted | This one doesn't seem instant. From my testing it will always fire circa 1 second late to account for the extra 1 second (white LED) at the end of the bomb countdown sequence. This also relies on the heartbeat packet, so settting your heartbeat packet lower will yield more accurate timing. I have mine set to 0.5 seconds for testing |
| 15 | EventBombDefused | Bomb was defused |  |
| 16 | EventBombExploded | Bomb exploded |  |
```

## Usage
### Using the go-package
You can use the go package in your own application to build more complicated tooling around game events.
`go get github.com/LukeyR/CS2-GameStateIntegration`

### Using a websocket
If yout not a fan of go, you can download the binary .exe from the latest release, and listen to events via a websocket.
Once you are running the websocket, connect to `ws://127.0.0.1:8000/ws`. You can provide the query paramter `Events` as many times as you like to choose which events to subscribe to. You can use the ID's or the event name form the table above.
e.g. `ws://127.0.0.1:8000/ws?Events=11&Events=EventPlayerAlivenessChanged` would get a message every time the player's health changed, or they died/respawned.

## Building
Nothing special here. Your usual `go build` process should work as normal.