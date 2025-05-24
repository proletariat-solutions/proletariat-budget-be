package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog/log"
)

type Configs struct {
	Auth    *Auth
	App     *App
	MongoDB *Mongodb
	HTTP    *HTTP
}

func Load() *Configs {
	cfg := &Configs{
		Auth:    &Auth{},
		App:     &App{},
		MongoDB: &Mongodb{},
		HTTP:    &HTTP{},
	}
	if err := env.Parse(cfg); err != nil {
		log.Fatal().Err(err).Msg("failed to load configs")
	}

	return cfg
}

type Auth struct {
	KeycloakHost    string `env:"KC_HOST" envDefault:"https://auth.hub-dev.invenco.com"`
	KeycloakRealm   string `env:"KC_REALM" envDefault:"invenco-hub"`
	ClientID        string `env:"AUTH_CLIENT_ID"`
	ClientSecret    string `env:"AUTH_CLIENT_SECRET"`
	ServiceUsername string `env:"AUTH_SERVICE_USERNAME"`
	ServicePassword string `env:"AUTH_SERVICE_PASSWORD"`
	LookupApiHost   string `env:"localhost:8080"`
}

type Mongodb struct {
	Database     string `env:"MONGO_DATABASE" envDefault:"mydb"`
	MyCollection string `env:"MONGO_MY_COLLECTION" envDefault:"mycoll"`
}

type HTTP struct {
	Timeout  time.Duration `env:"HTTP_CLIENT_TIMEOUT"   envDefault:"30s"`
	RetryMax int           `env:"HTTP_CLIENT_RETRY_MAX" envDefault:"3"`
}

type App struct {
	LogLevel    string        `env:"LOG_LEVEL" envDefault:"debug"`
	ServerPort  int           `env:"SERVER_PORT" envDefault:"8080"`
	ReadTimeout time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"500s"`
}
