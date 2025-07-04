package port

import (
	"context"
	"ghorkov32/proletariat-budget-be/openapi"
)

type SavingsGoalRepo interface {
	// SavingsGoal operations
	Create(ctx context.Context, savingsGoal openapi.SavingsGoalRequest) (string, error)
	Update(ctx context.Context, id string, savingsGoal openapi.SavingsGoalRequest) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*openapi.SavingsGoal, error)
	List(ctx context.Context) ([]openapi.SavingsGoal, error)
	MarkAsCompleted(ctx context.Context, id string) error
	MarkAsAbandoned(ctx context.Context, id string) error

	// Withdrawal operations
	CreateWithdrawal(ctx context.Context, withdrawal openapi.SavingsWithdrawalRequest) (string, error)
	DeleteWithdrawal(ctx context.Context, id string) error
	GetWithdrawalByID(ctx context.Context, id string) (*openapi.SavingsWithdrawal, error)

	// Contribution operations
	CreateContribution(ctx context.Context, contribution openapi.SavingsContributionRequest) (string, error)
	DeleteContribution(ctx context.Context, id string) error
	GetContributionByID(ctx context.Context, id string) (*openapi.SavingsContribution, error)

	// Transaction operations
	ListSavingsTransactions(ctx context.Context, params openapi.ListSavingsTransactionsParams) (*openapi.SavingsTransactionList, error)
}
