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

	// Categories
	ListCategories(ctx context.Context) ([]string, error)
	GetCategory(ctx context.Context, id string) (*openapi.SavingsCategory, error)
	CreateCategory(ctx context.Context, category openapi.SavingsCategoryRequest) (string, error)
	UpdateCategory(ctx context.Context, id string, category openapi.SavingsCategoryRequest) error
	DeleteCategory(ctx context.Context, id string) error

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

	// Tags operations
	FindOrCreateTags(ctx context.Context, tags []string) ([]string, error)
	LinkTagsToSavingsGoal(ctx context.Context, tags []string, savingsGoalId string) error
	LinkTagsToWithdrawal(ctx context.Context, tags []string, withdrawalId string) error
	LinkTagsToContribution(ctx context.Context, tags []string, contributionId string) error
	CreateTag(ctx context.Context, tag openapi.TagRequest) (string, error)
}
