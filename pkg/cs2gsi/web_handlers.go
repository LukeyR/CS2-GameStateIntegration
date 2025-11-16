package cs2gsi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func handlePOSTRequest(w http.ResponseWriter, r *http.Request) {
	gsiEvent, err := extractGSIEventFromRequest(r)
	if err != nil {
		return
	}
	gameEvents := findEvents(gsiEvent)
	if len(gameEvents) > 0 {
		for _, event := range gameEvents {
			for _, eventHandler := range gameEventHandlers[event.EventType] {
				eventHandler(gsiEvent, event)
			}
		}
	} else {
		for _, eventHandler := range gameNonEventHandlers {
			eventHandler(gsiEvent)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func handleWS(_ http.ResponseWriter, request *http.Request, conn *websocket.Conn) {

	queryParams := request.URL.Query()

	for _, e := range queryParams["Events"] {

		var eGameEvent events.GameEvent
		var eventExists bool

		eInt, err := strconv.Atoi(e)
		if err != nil {
			eGameEvent, eventExists = events.EventNameToEnum[e]
		} else {
			eGameEvent = events.GameEvent(eInt)
			_, eventExists = events.EnumToEventName[eGameEvent]
		}

		if !eventExists {
			log.Warn().Msgf("%v was requested from %v. Not a valid enum or EventName", e, conn)
			continue
		}

		RegisterEventHandler(eGameEvent, func(gsiEvent *structs.GSIEvent, gameEventDetails events.GameEventDetails) {
			err := conn.WriteJSON(gameEventDetails)
			if err != nil {
				return
			}
		})
	}
}

func extractGSIEventFromRequest(r *http.Request) (*structs.GSIEvent, error) {
	// Log the request body to stdout in Info level
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading body")
	}
	requestBodyFlat := &bytes.Buffer{}
	err = json.Compact(requestBodyFlat, requestBody)
	if err != nil {
		log.Error().Err(err).Msg("Error flattening")
	}

	event, err := structs.NewGSIEvent(string(requestBody))
	if err != nil {
		log.Error().Err(err).Str("original request", string(requestBody)).Msg("Error unmarshalling body")
	}
	return event, err
}
