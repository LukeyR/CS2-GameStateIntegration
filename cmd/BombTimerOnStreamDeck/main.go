package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"time"

	streamDeckEvents "github.com/LukeyR/CS2-GameStateIntegration/cmd/BombTimerOnStreamDeck/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/events"
	"github.com/LukeyR/CS2-GameStateIntegration/pkg/cs2gsi/structs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setUpLogger() {
	tempDirectory := filepath.Join(os.Getenv("localappdata"), "temp", "BombTimerOnStreamDeck")
	err := os.MkdirAll(tempDirectory, 0666)
	if err != nil {
		log.Fatal().Msg("Failed to create log directory")
	}

	currDt := time.Now().Format(cs2gsi.IsoTimestampFileUsable)
	logFilePath := filepath.Join(tempDirectory, currDt+".log")
	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Msg("Failed to create log file")
	}
	log.Logger = zerolog.New(logFile).Level(zerolog.InfoLevel).With().Timestamp().Logger()
	log.Info().Msg("Logger successfully configured.")
}

type CLIArgs struct {
	port          string
	info          interface{}
	uuid          string
	registerEvent string
}

func setUpCliFlags() *CLIArgs {
	portPtr := flag.String("port", "", "WS Port")
	uuidPrt := flag.String("pluginUUID", "", "Plugin UUID")
	registerEventPrt := flag.String("registerEvent", "", "Registration Event")
	infoPtr := flag.String("info", "{}", "StreamDeck Info")

	flag.Parse()
	var info interface{} // TODO: make this not an interface{}
	err := json.Unmarshal([]byte(*infoPtr), &info)
	if err != nil {
		log.Error().Err(err).Str("rawInfo", *infoPtr).Msg("Failed to parse info from cli args")
	}

	return &CLIArgs{
		port:          *portPtr,
		uuid:          *uuidPrt,
		registerEvent: *registerEventPrt,
		info:          info,
	}
}

func main() {
	setUpLogger()
	flags := setUpCliFlags()
	log.Info().Msgf("args: %+v", flags)

	shutdown := make(chan struct{})
	c := setupWebsocket(
		flags.port,
		&streamDeckEvents.RegisterEvent{Event: flags.registerEvent, UUID: flags.uuid},
	)
	c.AddEventCallback("willAppear",
		func(conn *StreamDeckConnection, e []byte) {
			log.Info().Bytes("Full msg", e).Send()
			var event streamDeckEvents.WillAppearEvent
			err := json.Unmarshal(e, &event)
			if err != nil {
				log.Error().Err(err).Msg("Couldn't unmarshall willAppear event")
			}
			log.Info().Msgf("Got context %s", event.Context)
			conn.Contexts[event.Context] = struct{}{}
			log.Info().Msgf("Saved context %s", event.Context)
			//_ = c.SetImage("images/black")
		})

	defer c.CloseWebsocket()

	cs2gsi.RegisterEventHandler(events.EventBombPlanted, func(gsiEvent *structs.GSIEvent, gameEvent events.GameEventDetails) {
		log.Info().Msg("Received bomb planted")
		const bombDuration = 40
		ticker := time.NewTicker(100 * time.Millisecond)
		done := make(chan bool)

		go func() {
			log.Info().Msg("starting to listen to ticker")
			start := time.Now().UnixNano() / int64(time.Millisecond)
			for {
				select {
				case <-done:
					return
				case _ = <-ticker.C:
					diff := float64(time.Now().UnixNano()/int64(time.Millisecond)-start) / 1000
					if bombDuration-diff > 0.0 {
						_ = c.SetTitle(strconv.FormatFloat(bombDuration-diff, 'f', 1, 64))
					} else {
						_ = c.SetTitle("")
						_ = c.SetState(1)
						time.Sleep(4570 * time.Millisecond) // Explosion gif is 4.57 seconds long
						_ = c.SetState(0)
					}
				}
			}
		}()

		time.Sleep((bombDuration + 1) * time.Second)
		ticker.Stop()
		log.Info().Msg("Stopping ticker")
		done <- true
	})

	log.Info().Msg("Starting up server cs listener")
	err := cs2gsi.StartupAndServe(":33942")
	if err != nil {
		return
	}

	<-shutdown
	log.Info().Msg("Shutting Down")
}
