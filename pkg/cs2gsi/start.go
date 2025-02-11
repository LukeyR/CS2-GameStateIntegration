package cs2gsi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"

	"github.com/rs/zerolog/log"
)

func handlePOSTRequest(w http.ResponseWriter, r *http.Request, loggers ExtraLoggers) {
	gsiEvent, err := extractGSIEventFromRequest(r, loggers)
	if err != nil {
		return
	}
	gameEvents := findEvents(gsiEvent)
	for _, event := range gameEvents {
		for _, eventHandler := range gameEventHandlers[event.EventType] {
			eventHandler(gsiEvent, event)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func extractGSIEventFromRequest(r *http.Request, loggers ExtraLoggers) (*structs.GSIEvent, error) {
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

	//log.Debug().Msg(string(requestBody))
	//fmt.Println(string(requestBody))
	loggers.data.Info().Msg(requestBodyFlat.String())

	event, err := structs.NewGSIEvent(string(requestBody))
	if err != nil {
		log.Error().Err(err).Str("original request", string(requestBody)).Msg("Error unmarshalling body")
	}
	return event, err
}

func StartupAndServe(addr string) error {
	loggers := setupLoggers()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("Registering handlers")
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			handlePOSTRequest(writer, request, loggers)
		default:
			errMsg := fmt.Sprintf("Unsupported Method: `%s`", request.Method)
			log.Error().Msg(errMsg)
			writer.WriteHeader(http.StatusNotFound)
			_, err := writer.Write([]byte(errMsg))
			if err != nil {
				return
			}
		}
	},
	)

	go func() {
		log.Info().Msg("Starting server")
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Panic().Err(err).Msg("Error when serving on HTTP")
		}
	}()

	<-shutdown
	// Any tidy up here
	log.Info().Msg("Shutting down server")

	return nil
}
