package main

import (
	"context"
	"database/sql"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/adapter/driven/mysql"
	"ghorkov32/proletariat-budget-be/internal/adapter/driving/middleware"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ghorkov32/proletariat-budget-be/config"
	"ghorkov32/proletariat-budget-be/internal/adapter/driving/resthttp"
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

	db := initDB(configs)

	defer db.Close()

	err := runMigrations(configs)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	ports := instantiatePorts(db)

	useCases := instantiateUseCases(ports)

	controller := resthttp.NewController(useCases)

	handler := openapi.NewStrictHandler(controller, nil)

	oapiSpecs, errSw := openapi.GetSwagger()
	if errSw != nil {
		log.Fatal().Err(errSw).Msg("failed to get OpenAPI specs")
	}

	httpServer := resthttp.NewHTTPServer(
		configs.App,
		openapi.Handler(handler),
		resthttp.MetricsCollector,
		middleware.DetailedRequestLogger,
		middleware.OpenAPIValidationMiddleware(oapiSpecs),
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

func runMigrations(configs *config.Configs) error {
	// Run database migrations
	migrationConfig := mysql.MigrationConfig{
		MigrationsPath: "./migrations/mysql",
		DBName:         configs.MySQL.Database,
		DBUser:         configs.MySQL.User,
		DBPassword:     configs.MySQL.Password,
		DBHost:         configs.MySQL.Host,
		DBPort:         configs.MySQL.Port,
	}

	if err := mysql.RunMigrations(migrationConfig); err != nil {
		return fmt.Errorf("failed to run database migrations: %w", err)
	}
	return nil
}

func initDB(configs *config.Configs) *sql.DB {
	// Initialize MySQL connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		configs.MySQL.User, configs.MySQL.Password, configs.MySQL.Host, configs.MySQL.Port, configs.MySQL.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MySQL")
	}

	// Configure connection pool
	db.SetMaxOpenConns(configs.MySQL.MaxOpenConns)
	db.SetMaxIdleConns(configs.MySQL.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(configs.MySQL.ConnMaxLife) * time.Second)
	return db
}

func instantiatePorts(db *sql.DB) *port.Ports {
	accountRepo := mysql.NewAccountRepo(db)
	authRepo := mysql.NewAuthRepository(db, os.Getenv("JWT_SECRET"))
	categoryRepo := mysql.NewCategoryRepo(db)
	tagsRepo := mysql.NewTagsRepo(db)
	expenditureRepo := mysql.NewExpenditureRepo(db, &tagsRepo)
	householdMembersRepo := mysql.NewHouseholdMemberRepository(db)
	ingressRepo := mysql.NewIngressRepo(db, &tagsRepo)
	savingsGoalRepo := mysql.NewSavingGoalRepo(db, &tagsRepo)
	transactionRepo := mysql.NewTransactionRepo(db, tagsRepo)
	return &port.Ports{
		Account:          &accountRepo,
		Auth:             &authRepo,
		Category:         &categoryRepo,
		Expenditure:      &expenditureRepo,
		HouseholdMembers: &householdMembersRepo,
		Ingress:          &ingressRepo,
		SavingGoal:       &savingsGoalRepo,
		Tags:             &tagsRepo,
		Transaction:      &transactionRepo,
	}
}

func instantiateUseCases(ports *port.Ports) *usecase.UseCases {
	account := usecase.NewAccountUseCase(ports.Account)
	auth := usecase.NewAuthUseCase(*ports.Auth)
	householdMember := usecase.NewHouseholdMemberUseCase(*ports.HouseholdMembers)

	return &usecase.UseCases{
		Account:         account,
		Auth:            auth,
		HouseholdMember: householdMember,
		// TODO:  Instantiate other use cases
	}
}
