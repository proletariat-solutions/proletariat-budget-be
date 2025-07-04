package mysql

import (
	"context"
	"database/sql"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
)

type TransactionRepoImpl struct {
	db *sql.DB
}

func (t TransactionRepoImpl) Create(ctx context.Context, transaction domain.Transaction) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t TransactionRepoImpl) Update(ctx context.Context, id string, transaction domain.Transaction) error {
	//TODO implement me
	panic("implement me")
}

func (t TransactionRepoImpl) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (t TransactionRepoImpl) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (t TransactionRepoImpl) Rollback(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (t TransactionRepoImpl) List(ctx context.Context, params openapi.ListTransactionsParams) (openapi.TransactionList, error) {
	//TODO implement me
	panic("implement me")
}

func (t TransactionRepoImpl) LinkWithEntity(ctx context.Context, entityID, entityType, transactionID string) error {
	//TODO implement me
	panic("implement me")
}
