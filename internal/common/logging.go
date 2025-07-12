package common

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"path/filepath"
	"strconv"
)

func SetupLogger(level string) {
	l, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse zerlog level")
	}
	zerolog.SetGlobalLevel(l)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
		NoColor:    false, // Enable colors
	}).With().Caller().Str("app", APPNAME).Logger()
}
