package domain

type Transaction struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	Amount   float64 `json:"amount"`
	Currency string `json:"currency"`
	TransactionDate string `json:"transaction_date"`
	Description string `json:"description"`
	TransactionType TransactionType `json:"transaction_type"`
	BalanceAfter float64 `json:"balance_after"`
	Status TransactionStatus `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
type TransactionType string

const (
	TransactionTypeExpenditure         TransactionType = "expenditure"
	TransactionTypeIngress             TransactionType = "ingress"
	TransactionTypeTransfer            TransactionType = "transfer"
	TransactionTypeSavingsContribution TransactionType = "savings_contribution"
	TransactionTypeSavingsWithdrawal   TransactionType = "savings_withdrawal"
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
	case TransactionTypeExpenditure, TransactionTypeIngress, TransactionTypeTransfer,
		TransactionTypeSavingsContribution, TransactionTypeSavingsWithdrawal:
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