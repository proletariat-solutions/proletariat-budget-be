package mysql

import (
	"database/sql"
	"errors"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
	"regexp"
)

// errorMapping maps MySQL error codes to infrastructure errors
var errorMapping = map[uint16]error{
	1048: port.ErrNotNullViolation,
	1062: port.ErrDuplicateKey,
	1064: port.ErrSyntaxError,
	1264: port.ErrDataOutOfRange,
	1265: port.ErrDataTruncated,
	1292: port.ErrInvalidDataFormat,
	1366: port.ErrInvalidDataFormat,
	1369: port.ErrConstraintViolation,
	1406: port.ErrDataTooLong,
	1411: port.ErrInvalidDataFormat,
	1416: port.ErrInvalidDataFormat,
	1451: port.ErrDependingForeignKey,
	1452: port.ErrForeignKeyNotFound,
	1525: port.ErrInvalidDataFormat,
	1644: port.ErrConstraintViolation,
	1690: port.ErrDataOutOfRange,
	3819: port.ErrConstraintViolation,
}

// translateError converts MySQL-specific errors to infrastructure errors
func translateError(err error) error {
	log.Error().Err(err).Msg("Receiving MySQL error")

	if errors.Is(
		err,
		sql.ErrNoRows,
	) {
		return port.ErrRecordNotFound
	}

	var mysqlErr *mysql.MySQLError

	// Try to parse MySQL error
	if errors.As(
		err,
		&mysqlErr,
	) {
		return handleMySQLError(mysqlErr)
	} else {
		return &port.InfrastructureError{
			Type:    "unknown_error",
			Message: err.Error(),
			Cause:   err,
		}
	}

}

// handleMySQLError processes MySQL-specific errors
func handleMySQLError(mysqlErr *mysql.MySQLError) error {
	if isForeignKeyError(mysqlErr.Number) {
		return ForeignKeyErrorMap[ForeignKeyConstraint(extractConstraintName(mysqlErr.Message))][mysqlErr.Number]
	}

	// Use mapping for standard errors
	if mappedErr, exists := errorMapping[mysqlErr.Number]; exists {
		return mappedErr
	}

	return port.ErrUnknownError
}

// isForeignKeyError checks if the error code represents a foreign key violation
func isForeignKeyError(code uint16) bool {
	return code == 1451 || code == 1452
}

// extractConstraintName extracts the constraint name from MySQL error message
func extractConstraintName(message string) string {
	re := regexp.MustCompile(`CONSTRAINT\s+` + "`" + `([^` + "`" + `]+)` + "`")
	matches := re.FindStringSubmatch(message)
	if len(matches) > 1 {
		return matches[1]
	}
	return "unknown_constraint"
}
