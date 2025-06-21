package port

import (
	"context"
	"ghorkov32/proletariat-budget-be/openapi"
)

type ExpenditureRepo interface {
	// Expenditure operations
	Create(ctx context.Context, expenditure openapi.Expenditure) (string, error)
	Update(ctx context.Context, id string, expenditure openapi.Expenditure) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*openapi.Expenditure, error)
	FindExpenditures(ctx context.Context, queryParams openapi.ListExpendituresParams) (*openapi.ExpenditureList, error)

	// Tag operations
	FindOrCreateTags(ctx context.Context, tags []string) ([]string, error)
	LinkTagsToExpenditure(ctx context.Context, tags []string, expenditureId string) error
	ListTags(ctx context.Context) ([]string, error)

	// Category operations
	FindCategory(ctx context.Context, id string) (*openapi.ExpenditureCategory, error)
	ListCategories(ctx context.Context) ([]string, error)
	CreateCategory(ctx context.Context, name, description string) (string, error)
	UpdateCategory(ctx context.Context, id string, name, description string) error
	DeleteCategory(ctx context.Context, id string) error
}
