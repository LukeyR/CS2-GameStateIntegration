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
)

func handlePOSTRequest(w http.ResponseWriter, r *http.Request, loggers LoggingHandler) {
	// Log the request body to stdout in Info level
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		loggers.raw.Error().Err(err).Msg("Error reading body")
	}
	requestBodyFlat := &bytes.Buffer{}
	err = json.Compact(requestBodyFlat, requestBody)
	if err != nil {
		loggers.raw.Error().Err(err).Msg("Error flattening")
	}

	//loggers.raw.Debug().Msg(string(requestBody))
	loggers.data.Info().Msg(requestBodyFlat.String())

	event := GSIEvent{}
	err = json.Unmarshal(requestBody, &event)
	if err != nil {
		loggers.raw.Error().Err(err).Msg("Error unmarshalling body")
	}
	fmt.Println(event)

	w.WriteHeader(http.StatusOK)
}

func StartupAndServe(addr string) error {
	loggers := setupLoggers()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	loggers.raw.Info().Msg("Registering handlers")
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			handlePOSTRequest(writer, request, loggers)
		default:
			errMsg := fmt.Sprintf("Unsupported Status code: `%s`", request.Method)
			loggers.raw.Error().Msg(errMsg)
			writer.WriteHeader(http.StatusNotFound)
			_, err := writer.Write([]byte(errMsg))
			if err != nil {
				return
			}
		}
	},
	)

	go func() {
		loggers.raw.Info().Msg("Starting server")
		if err := http.ListenAndServe(addr, nil); err != nil {
			loggers.raw.Panic().Err(err).Msg("Error when serving on HTTP")
		}
	}()

	<-shutdown
	// Any tidy up here
	loggers.raw.Info().Msg("Shutting down server")

	return nil
}
