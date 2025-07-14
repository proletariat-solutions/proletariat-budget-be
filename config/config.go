package config

import (
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

type Configs struct {
	App   *App
	MySQL *MySQL
	HTTP  *HTTP
}

func Load() *Configs {
	cfg := Configs{
		App:   &App{},
		MySQL: &MySQL{},
		HTTP:  &HTTP{},
	}
	setupFromLocalFile()
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Err(err).Msg("Unable to load environment variables")
	}

	return &cfg
}

type HTTP struct {
	Timeout  time.Duration `env:"HTTP_CLIENT_TIMEOUT"   envDefault:"30s"`
	RetryMax int           `env:"HTTP_CLIENT_RETRY_MAX" envDefault:"3"`
}

type App struct {
	LogLevel    string        `env:"LOG_LEVEL" envDefault:"debug"`
	ServerPort  int           `env:"SERVER_PORT" envDefault:"9091"`
	ReadTimeout time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"500s"`
}

// Add MySQL configuration
type MySQL struct {
	Host         string `env:"MYSQL_HOST" envDefault:"localhost"`
	Port         string `env:"MYSQL_PORT" envDefault:"3306"`
	Database     string `env:"MYSQL_DATABASE" envDefault:"proletariat_budget"`
	User         string `env:"MYSQL_USER" envDefault:"budgetuser"`
	Password     string `env:"MYSQL_PASSWORD" envDefault:"budgetpassword"`
	MaxOpenConns int    `env:"MYSQL_MAX_OPEN_CONNS" envDefault:"10"`
	MaxIdleConns int    `env:"MYSQL_MAX_IDLE_CONNS" envDefault:"5"`
	ConnMaxLife  int    `env:"MYSQL_CONN_MAX_LIFETIME" envDefault:"300"` // seconds
}

// ParseConfig parses environment variables into the config struct
func setupFromLocalFile() {
	err := godotenv.Load(".env") // DO NOT commit the .env file

	if err != nil {
		if os.IsNotExist(err) {
			log.Warn().Msg(".env file not found, proceeding with original values")
		} else {
			log.Err(err).Msg("failed to read .env file, proceeding with original values")
		}
	} else {
		log.Info().Msg(".env read successfully")
	}
}
