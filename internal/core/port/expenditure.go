package port

import (
	"context"
	"ghorkov32/proletariat-budget-be/openapi"
)

type ExpenditureRepo interface {
	// Expenditure operations
	Create(ctx context.Context, expenditure openapi.ExpenditureRequest, transactionID string) (string, error)
	GetByID(ctx context.Context, id string) (*openapi.Expenditure, error)
	FindExpenditures(ctx context.Context, queryParams openapi.ListExpendituresParams) (*openapi.ExpenditureList, error)
}
