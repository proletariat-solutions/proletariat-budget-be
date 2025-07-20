package port

import (
	"context"

	"ghorkov32/proletariat-budget-be/internal/core/domain"
)

type ExpenditureRepo interface {
	// Expenditure operations
	Create(ctx context.Context, expenditure domain.Expenditure) (string, error)
	GetByID(ctx context.Context, id string) (*domain.Expenditure, error)
	FindExpenditures(ctx context.Context, queryParams domain.ExpenditureListParams) (*domain.ExpenditureList, error)
}
