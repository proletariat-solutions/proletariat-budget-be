package port

import (
	"context"

	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
)

type TransactionRepo interface {
	Create(ctx context.Context, transaction domain.Transaction) (string, error)
	GetByID(ctx context.Context, id string) (*domain.Transaction, error)
	List(ctx context.Context, params openapi.ListTransactionsParams) (*openapi.TransactionList, error)
}
