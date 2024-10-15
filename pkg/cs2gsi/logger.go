package cs2gsi

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const IsoTimestampFileUsable = "2006-01-02T15-04-05"

type LoggingHandler struct {
	raw  zerolog.Logger
	data zerolog.Logger
}

func setupLoggers() LoggingHandler {

	LogsFilePath := "./logs"
	DataLogsFilePath := filepath.Join(LogsFilePath, "data")
	RawLogsFilePath := filepath.Join(LogsFilePath, "raw")

	currDt := time.Now().Format(IsoTimestampFileUsable)
	rawFilePath := filepath.Join(RawLogsFilePath, currDt+".log")
	rawFile, err := os.OpenFile(rawFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal().Err(err).Str("File path", rawFilePath).Msg("Error creating or opening raw log file")
	}

	dataFilePath := filepath.Join(DataLogsFilePath, currDt+".log")
	dataFile, err := os.OpenFile(filepath.Join(DataLogsFilePath, currDt+".log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal().Err(err).Str("File path", dataFilePath).Msg("Error creating or opening data log file")
	}

	rawLogger := zerolog.New(zerolog.MultiLevelWriter(rawFile, os.Stdout)).With().Timestamp().Logger()
	dataLogger := zerolog.New(zerolog.ConsoleWriter{
		Out:        zerolog.MultiLevelWriter(dataFile, os.Stdout),
		NoColor:    true,
		PartsOrder: []string{zerolog.MessageFieldName},
	},
	)

	return LoggingHandler{rawLogger, dataLogger}
}
