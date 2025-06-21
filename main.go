package main

import (
	"context"
	"database/sql"
	"fmt"
	mysql2 "ghorkov32/proletariat-budget-be/internal/adapter/driven/mysql"
	resthttp2 "ghorkov32/proletariat-budget-be/internal/adapter/driving/resthttp"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ghorkov32/proletariat-budget-be/config"
	"ghorkov32/proletariat-budget-be/internal/adapter/resthttp"
	"ghorkov32/proletariat-budget-be/internal/common"
	"ghorkov32/proletariat-budget-be/internal/core/usecase"
	"ghorkov32/proletariat-budget-be/openapi"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("application starting ...")
	configs := config.Load()
	common.SetupLogger(configs.App.LogLevel)

	appCtx := context.Background()

	accessChecker := usecase.NewAccessChecker(scopeLookup)
	userSvc := usecase.NewCreateUser(userRepo)
	userController := resthttp.NewUserController(userSvc, accessChecker)
	userHandler := openapi.NewStrictHandler(userController, nil)

	// Run database migrations
	migrationConfig := mysql2.MigrationConfig{
		MigrationsPath: "./migrations/mysql",
		DBName:         configs.MySQL.Database,
		DBUser:         configs.MySQL.User,
		DBPassword:     configs.MySQL.Password,
		DBHost:         configs.MySQL.Host,
		DBPort:         configs.MySQL.Port,
	}

	if err := mysql2.RunMigrations(migrationConfig); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	// Initialize MySQL connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		configs.MySQL.User, configs.MySQL.Password, configs.MySQL.Host, configs.MySQL.Port, configs.MySQL.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MySQL")
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(configs.MySQL.MaxOpenConns)
	db.SetMaxIdleConns(configs.MySQL.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(configs.MySQL.ConnMaxLife) * time.Second)

	// Initialize repositories
	accountRepo := mysql2.NewAccountRepo(db)

	oapiSpec, errSw := openapi.GetSwagger()
	if errSw != nil {
		log.Warn().Err(errSw).Msg("unable to fetch openapi spec")
	}

	httpServer := resthttp2.NewHTTPServer(
		configs.App,
		openapi.Handler(userHandler),
		resthttp2.MetricsCollector,
		// TODO:  Middlewares
	)
	go httpServer.Start()
	defer func(ctx context.Context) {
		errS := httpServer.Shutdown(ctx)
		if errS != nil {
			log.Err(errS).Msg("http server shutdown error")
		}
	}(appCtx)

	ctx, stop := signal.NotifyContext(appCtx, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	log.Info().Msg("waiting for app exiting conditions")

	<-ctx.Done()
	log.Info().Msg("application stopping ...")
}
