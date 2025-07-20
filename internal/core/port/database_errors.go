package port

import (
	"errors"
	"fmt"
)

var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrDuplicateKey        = errors.New("duplicate key constraint violation")
	ErrForeignKeyViolation = errors.New("foreign key constraint violation")
	ErrConnectionFailed    = errors.New("database connection failed")
	ErrSyntaxError         = errors.New("syntax error")
	ErrUnknownError        = errors.New("unknown error")
	ErrNotNullViolation    = errors.New("not null constraint violation")
	ErrDataOutOfRange      = errors.New("data value out of range")
	ErrDataTooLong         = errors.New("data too long for column")
	ErrInvalidDataFormat   = errors.New("invalid data format")
	ErrDataTruncated       = errors.New("data truncated")
	ErrForeignKeyNotFound  = errors.New("foreign key not found")
	ErrDependingForeignKey = errors.New("dependent foreign key constraint")
	ErrConstraintViolation = errors.New("constraint violation")
)

type EntityRelationshipError struct {
	Type                string
	Message             string
	OffendingConstraint string
	Cause               error
}

func (e EntityRelationshipError) Error() string {
	return fmt.Sprintf(
		"%s: %s (offending constraint: %s)",
		e.Type,
		e.Message,
		e.OffendingConstraint,
	)
}

type InfrastructureError struct {
	Type    string
	Message string
	Cause   error
}

func (e *InfrastructureError) Error() string {
	return e.Message
}

func (e *InfrastructureError) Unwrap() error {
	return e.Cause
}

func IsInfrastructureError(err error) bool {
	var infraErr *InfrastructureError
	return err != nil &&
		(errors.As(
			err,
			&infraErr,
		) ||
			errors.Is(
				err,
				ErrConnectionFailed,
			) ||
			errors.Is(
				err,
				ErrSyntaxError,
			) ||
			errors.Is(
				err,
				ErrUnknownError,
			))
}
