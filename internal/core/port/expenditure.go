package port

import (
	"context"
	"ghorkov32/proletariat-budget-be/openapi"
)

type ExpenditureRepo interface {
	// Expenditure operations
	Create(ctx context.Context, expenditure openapi.Expenditure, transactionID string) (string, error)
	Update(ctx context.Context, id string, expenditure openapi.Expenditure) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*openapi.Expenditure, error)
	FindExpenditures(ctx context.Context, queryParams openapi.ListExpendituresParams) (*openapi.ExpenditureList, error)
}
