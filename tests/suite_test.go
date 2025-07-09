package tests

import (
	"database/sql"
	"ghorkov32/proletariat-budget-be/config"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	config *config.App
	ports  *port.Ports
	db     *sql.DB
}
