package integration_tests

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ghorkov32/proletariat-budget-be/config"
	"ghorkov32/proletariat-budget-be/integration_tests/containers"
	"ghorkov32/proletariat-budget-be/integration_tests/utils"
	"ghorkov32/proletariat-budget-be/internal/adapter/driven/mysql"
	"ghorkov32/proletariat-budget-be/internal/adapter/driving/middleware"
	"ghorkov32/proletariat-budget-be/internal/adapter/driving/resthttp"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"ghorkov32/proletariat-budget-be/internal/core/usecase"
	"ghorkov32/proletariat-budget-be/openapi"
	mysql2 "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	config      *config.Configs
	dbContainer *containers.MysqlContainer
	db          *sql.DB
	server      *resthttp.App
	ctx         context.Context
}

func TestIntegrationSuite(t *testing.T) {
	// ignore integration tests if flag skip is used
	if testing.Short() {
		t.Skip("skipping integration tests...")
	}

	// run all integration tests with IntegrationSuite as receiver
	suite.Run(t, &Suite{})
}
func (s *Suite) SetupSuite() {
	var err error

	s.config = config.Load()
	if false { // change this to true to spin up a container
		s.dbContainer = containers.NewMysqlContainer()
		s.config.MySQL, err = s.dbContainer.InitContainer(s.config.MySQL)

	}
	s.handleErr(err, "failed to initialize MySQL container")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		s.config.MySQL.User, s.config.MySQL.Password, s.config.MySQL.Host, s.config.MySQL.Port, s.config.MySQL.Database)

	db, err := sql.Open("mysql", dsn)
	s.handleErr(err, "failed to connect to MySQL")
	s.db = db

	s.ctx = context.Background()

	err = s.runMigrations()
	s.handleErr(err, "failed to run migrations")

	ports := instantiatePorts(db)

	useCases := instantiateUseCases(ports)

	controller := resthttp.NewController(useCases)

	handler := openapi.NewStrictHandler(controller, nil)

	oapiSpecs, errSw := openapi.GetSwagger()
	if errSw != nil {
		s.FailNow("failed to get OpenAPI specs")
	}

	s.server = resthttp.NewHTTPServer(
		s.config.App,
		openapi.Handler(handler),
		middleware.DetailedRequestLogger,
		middleware.OpenAPIValidationMiddleware(oapiSpecs),
	)
	go s.server.Start()

	// wait for the server to start
	time.Sleep(1 * time.Second)
}

func (s *Suite) handleErr(err error, msg string) {
	if err != nil {
		s.FailNow(fmt.Sprintf("%v: %v", msg, err))
	}
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
	account := usecase.NewAccountUseCase(ports.Account, ports.HouseholdMembers)
	auth := usecase.NewAuthUseCase(*ports.Auth)
	householdMember := usecase.NewHouseholdMemberUseCase(*ports.HouseholdMembers)

	return &usecase.UseCases{
		Account:         account,
		Auth:            auth,
		HouseholdMember: householdMember,
		// TODO:  Instantiate other use cases
	}
}

func (s *Suite) runMigrations() error {
	// Run database migrations
	migrationConfig := mysql.MigrationConfig{
		MigrationsPath: "../migrations/mysql",
		DBName:         s.config.MySQL.Database,
		DBUser:         s.config.MySQL.User,
		DBPassword:     s.config.MySQL.Password,
		DBHost:         s.config.MySQL.Host,
		DBPort:         s.config.MySQL.Port,
	}

	if err := mysql.RunMigrations(migrationConfig); err != nil {
		return fmt.Errorf("failed to run database migrations: %w", err)
	}

	err := utils.ExecuteSQLFile(s.ctx, s.db, "./mock_data/mock-data.sql")
	if err != nil {
		return err
	}
	return nil
}

func (s *Suite) TearDownSuite() {
	s.ClearTables()
	err := s.server.Shutdown(s.ctx)
	if err != nil {
		s.handleErr(err, "failed to shutdown http server")
	}
	err = s.db.Close()
	if err != nil {
		s.handleErr(err, "failed to close database connection")
	}
	err = s.server.Shutdown(context.Background())
	if err != nil {
		s.handleErr(err, "failed to shutdown http server")
	}
}

func (s *Suite) ClearTables() {
	// Execute each statement separately for better error handling
	statements := []string{
		"SET FOREIGN_KEY_CHECKS = 0",
		"TRUNCATE TABLE proletariat_budget.accounts",
		"TRUNCATE TABLE proletariat_budget.categories",
		"TRUNCATE TABLE proletariat_budget.exchange_rates",
		"TRUNCATE TABLE proletariat_budget.expenditure_tags",
		"TRUNCATE TABLE proletariat_budget.expenditures",
		"TRUNCATE TABLE proletariat_budget.household_members",
		"TRUNCATE TABLE proletariat_budget.ingress_recurrence_patterns",
		"TRUNCATE TABLE proletariat_budget.ingress_tags",
		"TRUNCATE TABLE proletariat_budget.ingresses",
		"TRUNCATE TABLE proletariat_budget.roles",
		"TRUNCATE TABLE proletariat_budget.savings_contribution_tags",
		"TRUNCATE TABLE proletariat_budget.savings_contributions",
		"TRUNCATE TABLE proletariat_budget.savings_goal_tags",
		"TRUNCATE TABLE proletariat_budget.savings_goals",
		"TRUNCATE TABLE proletariat_budget.savings_withdrawal_tags",
		"TRUNCATE TABLE proletariat_budget.savings_withdrawals",
		"TRUNCATE TABLE proletariat_budget.tags",
		"TRUNCATE TABLE proletariat_budget.transaction_rollbacks",
		"TRUNCATE TABLE proletariat_budget.transactions",
		"TRUNCATE TABLE proletariat_budget.transfers",
		"TRUNCATE TABLE proletariat_budget.user_roles",
		"TRUNCATE TABLE proletariat_budget.users",
		"SET FOREIGN_KEY_CHECKS = 1",
	}

	for _, stmt := range statements {
		_, err := s.db.Exec(stmt)
		if err != nil {
			mysqlError := &mysql2.MySQLError{}
			if errors.As(err, &mysqlError) {
				if mysqlError.Number != 1146 {
					s.handleErr(err, fmt.Sprintf("failed to execute: %s", stmt))
				}
			} else {
				s.handleErr(err, fmt.Sprintf("failed to execute: %s", stmt))
			}
		}
	}
}
