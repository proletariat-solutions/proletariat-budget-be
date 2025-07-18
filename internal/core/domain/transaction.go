package domain

import (
	"time"
)

type Transaction struct {
	ID              *string            `json:"id"`
	AccountID       string             `json:"account_id"`
	Amount          float32            `json:"amount"`
	Currency        string             `json:"currency"`
	TransactionDate time.Time          `json:"transaction_date"`
	Description     string             `json:"description"`
	TransactionType TransactionType    `json:"transaction_type"`
	BalanceAfter    *float32           `json:"balance_after"`
	Status          *TransactionStatus `json:"status"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}
type TransactionType string

const (
	TransactionTypeExpenditure TransactionType = "expenditure"
	TransactionTypeIngress     TransactionType = "ingress"
	TransactionTypeTransfer    TransactionType = "transfer"
	TransactionTypeRollback    TransactionType = "rollback"
)

// String returns the string representation of the transaction type
func (t TransactionType) String() string {
	return string(t)
}

// TransactionStatus represents the status of a financial transaction
type TransactionStatus string

// IsValid checks if the transaction type is valid
func (t TransactionType) IsValid() bool {
	switch t {
	case TransactionTypeExpenditure, TransactionTypeIngress, TransactionTypeTransfer:
		return true
	}
	return false
}

// Transaction status enum values
const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

// IsValid checks if the transaction status is valid
func (s TransactionStatus) IsValid() bool {
	switch s {
	case TransactionStatusPending, TransactionStatusCompleted,
		TransactionStatusFailed, TransactionStatusCancelled:
		return true
	}
	return false
}

// String returns the string representation of the transaction status
func (s TransactionStatus) String() string {
	return string(s)
}

func RollbackTransaction(t *Transaction) *Transaction {
	statusCompleted := TransactionStatusCompleted
	balanceAfter := *t.BalanceAfter + t.Amount
	return &Transaction{
		ID:              t.ID,
		AccountID:       t.AccountID,
		Amount:          -t.Amount,
		Currency:        t.Currency,
		TransactionDate: time.Now(),
		Description:     "Rollback of transaction " + t.Description,
		TransactionType: TransactionTypeRollback,
		BalanceAfter:    &balanceAfter,
		Status:          &statusCompleted,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}
