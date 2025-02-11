package cs2gsi

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{}

func StartupAndServe(addr string) error {
	loggers := setupLoggers()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("Registering handlers")
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			HandlePOSTRequest(writer, request, loggers)
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

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		log.Info().Msg("Websocket Connection requested")
		conn, err := upgrader.Upgrade(writer, request, nil)
		defer func() {
			err := conn.Close()
			if err != nil {
				log.Err(err).Msg("Failed websocket upgrade")
			}
		}()
		log.Info().Msgf("Websocket Connection upgrade success: %v", request.URL)

		if err != nil {
			log.Err(err)
			return
		}

		HandleWS(writer, request, conn)

		<-shutdown
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
