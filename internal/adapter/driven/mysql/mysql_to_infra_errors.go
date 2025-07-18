package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
	"regexp"
)

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
	if errors.As(
		err,
		&mysqlErr,
	) {
		switch mysqlErr.Number {
		case 1048: // Column cannot be null
			return port.ErrNotNullViolation
		case 1062: // Duplicate entry
			return port.ErrDuplicateKey
		case 1064: // SQL syntax error
			return port.ErrSyntaxError
		case 1264: // Out of range value for column
			return port.ErrDataOutOfRange
		case 1265: // Data truncated for column
			return port.ErrDataTruncated
		case 1292: // Incorrect datetime value
			return port.ErrInvalidDataFormat
		case 1366: // Incorrect string value for column
			return port.ErrInvalidDataFormat
		case 1369: // CHECK constraint failed
			return port.ErrConstraintViolation
		case 1406: // Data too long for column
			return port.ErrDataTooLong
		case 1411: // Incorrect datetime value for function
			return port.ErrInvalidDataFormat
		case 1416: // Cannot get geometry object from data
			return port.ErrInvalidDataFormat
		case 1451:
		case 1452: // Foreign key constraint fails
			return fmt.Errorf(
				"%w: %s",
				port.ErrForeignKeyViolation,
				extractConstraintName(mysqlErr.Message),
			)
		case 1525: // Incorrect DECIMAL value
			return port.ErrInvalidDataFormat
		case 1644: // Unhandled user-defined exception condition (SIGNAL)
			return port.ErrConstraintViolation
		case 1690: // BIGINT value is out of range
			return port.ErrDataOutOfRange
		case 3819: // Check constraint violation
			return port.ErrConstraintViolation
		default:
			return port.ErrUnknownError
		}
	}

	return &port.InfrastructureError{
		Type:    "unknown_error",
		Message: err.Error(),
		Cause:   err,
	}
}

// extractConstraintName extracts the constraint name from MySQL error message
func extractConstraintName(message string) string {
	// MySQL foreign key error messages typically contain the constraint name
	// Example: "Cannot add or update a child row: a foreign key constraint fails (`db`.`table`, CONSTRAINT `constraint_name` FOREIGN KEY ...)"
	re := regexp.MustCompile(`CONSTRAINT\s+` + "`" + `([^` + "`" + `]+)` + "`")
	matches := re.FindStringSubmatch(message)
	if len(matches) > 1 {
		return matches[1]
	}
	return "unknown_constraint"
}
