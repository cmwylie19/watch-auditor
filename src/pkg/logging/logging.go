package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogging(level string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	//log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true}).With().Logger()

	setLevel(level)
}

func setLevel(level string) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Error().Msg("Invalid log level")
		return
	}
	zerolog.SetGlobalLevel(lvl)
}

func Info(message string) {
	log.Info().Msg(message)
}

func Debug(message string) {
	log.Debug().Msg(message)
}

func Error(message string) {
	log.Error().Msg(message)
}
