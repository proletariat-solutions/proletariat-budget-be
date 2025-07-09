package domain

import "errors"

var (
	ErrEntityNotFound         = errors.New("entity not found")
	ErrAccountHasTransactions = errors.New("account has transactions")
)
