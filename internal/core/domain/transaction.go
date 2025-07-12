package domain

import (
	"ghorkov32/proletariat-budget-be/openapi"
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
	BalanceAfter    *float64           `json:"balance_after"`
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

func FromOAPIExpenditure(e *openapi.Expenditure) *Transaction {
	return &Transaction{
		AccountID:       e.AccountId,
		Amount:          e.Amount,
		Currency:        *e.Currency,
		TransactionDate: e.Date.Time,
		Description:     *e.Description,
		TransactionType: TransactionTypeExpenditure,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

func FromOAPIExpenditureRequest(e *openapi.ExpenditureRequest) *Transaction {
	return &Transaction{
		AccountID:       e.AccountId,
		Amount:          e.Amount,
		Currency:        *e.Currency,
		TransactionDate: e.Date.Time,
		Description:     *e.Description,
		TransactionType: TransactionTypeExpenditure,
	}
}

func FromOAPIIngress(i *openapi.Ingress) *Transaction {
	return &Transaction{
		AccountID:       i.AccountId,
		Amount:          i.Amount,
		Currency:        i.Currency,
		TransactionDate: i.Date.Time,
		Description:     *i.Description,
		TransactionType: TransactionTypeIngress,
		CreatedAt:       *i.CreatedAt,
	}
}

func FromOAPIIngressRequest(i *openapi.IngressRequest) *Transaction {
	return &Transaction{
		AccountID:       i.AccountId,
		Amount:          i.Amount,
		Currency:        i.Currency,
		TransactionDate: i.Date.Time,
		Description:     *i.Description,
		TransactionType: TransactionTypeIngress,
	}
}

func FromOAPITransferDebit(t *openapi.Transfer, sourceAccountCurrency string) *Transaction {
	return &Transaction{
		AccountID:       t.SourceAccountId,
		Amount:          *t.SourceAmount,
		Currency:        sourceAccountCurrency,
		TransactionDate: t.Date.Time,
		Description:     *t.Description,
		TransactionType: TransactionTypeTransfer,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func FromOAPITransferCredit(t *openapi.Transfer, destinationAccountCurrency string) *Transaction {
	return &Transaction{
		AccountID:       t.DestinationAccountId,
		Amount:          *t.DestinationAmount,
		Currency:        destinationAccountCurrency,
		TransactionDate: t.Date.Time,
		Description:     *t.Description,
		TransactionType: TransactionTypeTransfer,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func RollbackTransaction(t *Transaction) *Transaction {
	statusCompleted := TransactionStatusCompleted
	balanceAfter := *t.BalanceAfter + float64(t.Amount)
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
