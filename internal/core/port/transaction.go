package port

import (
	"context"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
)

type TransactionRepo interface {
	//
	Create(ctx context.Context, transaction domain.Transaction) (string, error)
	Update(ctx context.Context, id string, transaction domain.Transaction) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*domain.Transaction, error)
	Rollback(ctx context.Context, id string) error
	List(ctx context.Context, params openapi.ListTransactionsParams) (openapi.TransactionList, error)

	/*
		Linking operations
	*/
	LinkWithEntity(ctx context.Context, entityID, entityType, transactionID string) error
}
