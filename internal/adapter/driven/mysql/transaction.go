package mysql

import (
	"context"
	"database/sql"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
)

type TransactionRepo struct {
	db *sql.DB
}

func (t TransactionRepo) Create(ctx context.Context, transaction domain.Transaction) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t TransactionRepo) Update(ctx context.Context, id string, transaction domain.Transaction) error {
	//TODO implement me
	panic("implement me")
}

func (t TransactionRepo) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (t TransactionRepo) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (t TransactionRepo) Rollback(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (t TransactionRepo) List(ctx context.Context, params openapi.ListTransactionsParams) (openapi.TransactionList, error) {
	//TODO implement me
	panic("implement me")
}

func (t TransactionRepo) LinkWithExpenditure(ctx context.Context, expenditureID string, transactionID string) error {
	//TODO implement me
	panic("implement me")
}

func (t TransactionRepo) LinkWithIngress(ctx context.Context, ingressID string, transactionID string) error {
	//TODO implement me
	panic("implement me")
}

func (t TransactionRepo) LinkWithTransfer(ctx context.Context, transferID string, transactionID string) error {
	//TODO implement me
	panic("implement me")
}

func (t TransactionRepo) LinkWithSavingsContribution(ctx context.Context, contributionID string, transactionID string) error {
	//TODO implement me
	panic("implement me")
}

func (t TransactionRepo) LinkWithSavingsWithdrawal(ctx context.Context, withdrawalID string, transactionID string) error {
	//TODO implement me
	panic("implement me")
}
