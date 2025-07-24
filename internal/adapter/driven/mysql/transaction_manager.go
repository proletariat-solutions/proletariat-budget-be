package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"ghorkov32/proletariat-budget-be/internal/core/port"
)

type MySQLTransactionManager struct {
	db *sql.DB
}

func NewMySQLTransactionManager(db *sql.DB) port.TransactionManager {
	return &MySQLTransactionManager{db: db}
}

func (m *MySQLTransactionManager) WithTransaction(
	ctx context.Context,
	fn func(ctx context.Context, tx port.TransactionContext) error,
) error {
	// Begin transaction
	sqlTx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create transaction context with repositories that use the transaction
	txCtx := &MySQLTransactionContext{
		tx: sqlTx,
	}

	// Execute the function
	err = fn(ctx, txCtx)
	if err != nil {
		// Rollback on error
		if rollbackErr := sqlTx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w (original error: %v)", rollbackErr, err)
		}
		return err
	}

	// Commit transaction
	if err := sqlTx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

type MySQLTransactionContext struct {
	tx *sql.Tx
}

func (m *MySQLTransactionContext) GetAccountRepo() port.AccountRepo {
	return NewAccountRepoWithTx(m.tx)
}

func (m *MySQLTransactionContext) GetTransactionRepo() port.TransactionRepo {
	return NewTransactionRepoWithTx(m.tx)
}

func (m *MySQLTransactionContext) GetExpenditureRepo() port.ExpenditureRepo {
	return NewExpenditureRepoWithTx(m.tx)
}

func (m *MySQLTransactionContext) GetCategoryRepo() port.CategoryRepo {
	return NewCategoryRepoWithTx(m.tx)
}

func (m *MySQLTransactionContext) GetTagsRepo() port.TagsRepo {
	return NewTagsRepoWithTx(m.tx)
}

func (m *MySQLTransactionContext) GetHouseholdMemberRepo() port.HouseholdMemberRepo {
	return NewHouseholdMemberRepoWithTx(m.tx)
}
