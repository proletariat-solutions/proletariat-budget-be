package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"strconv"
	"strings"
	"time"
)

type AccountRepoImpl struct {
	db *sql.DB
}

func NewAccountRepo(db *sql.DB) port.AccountRepo {
	return &AccountRepoImpl{db: db}
}

func (r *AccountRepoImpl) Create(ctx context.Context, account domain.Account) (*string, error) {
	query := `
        INSERT INTO accounts (
            name, type, institution, currency, initial_balance, 
            current_balance, active, description, account_number,  owner,
            account_information, created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), now())
    `

	result, err := r.db.ExecContext(
		ctx,
		query,
		account.Name,
		account.Type,
		account.Institution,
		account.Currency,
		account.InitialBalance,
		account.InitialBalance,
		account.Active,
		account.Description,
		account.AccountNumber,
		account.OwnerID,
		account.AccountInformation,
	)

	if err != nil {
		return nil, translateError(err)
	}

	lastID, err := result.LastInsertId()

	if err != nil {
		return nil, translateError(err)
	}

	lastIDStr := strconv.FormatInt(lastID, 10)

	return &lastIDStr, nil
}

func (r *AccountRepoImpl) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	query := `SELECT 
				a.id, a.name, type, institution, currency, 
				initial_balance, current_balance, a.active, 
				description, account_number, account_information, 
				a.created_at, a.updated_at, a.owner, hm.id, hm.name, hm.surname, hm.nickname, hm.role, hm.active, hm.created_at, hm.updated_at
				FROM accounts a left join proletariat_budget.household_members hm on a.owner = hm.id  WHERE a.id =?`

	account := &domain.Account{
		Owner: &domain.HouseholdMember{},
	}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.Name,
		&account.Type,
		&account.Institution,
		&account.Currency,
		&account.InitialBalance,
		&account.CurrentBalance,
		&account.Active,
		&account.Description,
		&account.AccountNumber,
		&account.AccountInformation,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.OwnerID,
		&account.Owner.ID,
		&account.Owner.FirstName,
		&account.Owner.LastName,
		&account.Owner.Nickname,
		&account.Owner.Role,
		&account.Owner.Active,
		&account.Owner.CreatedAt,
		&account.Owner.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, port.ErrRecordNotFound
	} else if err != nil {
		return nil, translateError(err)
	}

	return account, nil

}

func (r *AccountRepoImpl) Update(ctx context.Context, account domain.Account) error {
	query := `
        UPDATE accounts SET 
            name =?, type =?, institution =?, currency =?, initial_balance =?, 
            current_balance =?, active =?, description =?, account_number =?, 
            account_information =?, updated_at =?, owner =? 
        WHERE id =?
    `

	now := time.Now()
	account.UpdatedAt = now

	_, err := r.db.ExecContext(
		ctx,
		query,
		account.Name,
		account.Type,
		account.Institution,
		account.Currency,
		account.InitialBalance,
		account.CurrentBalance,
		account.Active,
		account.Description,
		account.AccountNumber,
		account.AccountInformation,
		account.UpdatedAt,
		account.OwnerID,
		account.ID,
	)

	if err != nil {
		return translateError(err)
	}
	return nil
}

func (r *AccountRepoImpl) SetActive(ctx context.Context, id string, active bool) error {
	query := `UPDATE accounts SET active = ?, updated_at = NOW() WHERE id =?`
	result, err := r.db.ExecContext(
		ctx,
		query,
		active,
		id,
	)
	if err != nil {
		return translateError(err)
	}
	rowsAffected, errRowsAffected := result.RowsAffected()
	if errRowsAffected != nil {
		return fmt.Errorf("failed to delete account: %w", errRowsAffected)
	}
	if rowsAffected == 0 {
		return port.ErrRecordNotFound
	}
	return nil
}

func (r *AccountRepoImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM accounts WHERE id =?`
	result, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)
	if err != nil {
		return translateError(err)
	}
	rowsAffected, errRowsAffected := result.RowsAffected()
	if errRowsAffected != nil {
		return translateError(errRowsAffected)
	}
	if rowsAffected == 0 {
		return port.ErrRecordNotFound
	}
	return nil
}

func (r *AccountRepoImpl) List(ctx context.Context, params domain.AccountListParams) (*domain.AccountList, error) {
	query := `SELECT a.id,
					   a.name,
					   type,
					   institution,
					   currency,
					   initial_balance,
					   current_balance,
					   a.active,
					   description,
					   account_number,
					   account_information,
					   a.created_at,
					   a.updated_at,
					   hm.id,
					   hm.name,
					   hm.surname,
					   hm.nickname,
					   hm.role,
					   hm.active,
					   hm.created_at,
					   hm.updated_at
				FROM accounts a
						 left join household_members hm on a.owner = hm.id`

	whereClause := make([]string, 0)
	args := make([]any, 0)
	if params.Currency != nil {
		whereClause = append(whereClause, "a.currency =?")
		args = append(args, *params.Currency)
	}
	if params.Type != nil {
		whereClause = append(whereClause, "a.type =?")
		args = append(args, *params.Type)
	}
	if params.Active != nil {
		whereClause = append(whereClause, "a.active =?")
		args = append(args, *params.Active)
	}

	queryCount := "SELECT COUNT(*) FROM accounts a"
	if len(whereClause) > 0 {
		query += " WHERE " + strings.Join(whereClause, " AND ")
		queryCount += " WHERE " + strings.Join(whereClause, " AND ")

	}
	query += " ORDER BY a.created_at DESC"
	stmtCount, errQueryCountStmt := r.db.PrepareContext(ctx, queryCount)
	if errQueryCountStmt != nil {
		return nil, translateError(errQueryCountStmt)
	}
	var count int
	errCount := stmtCount.QueryRowContext(ctx, args...).Scan(&count)
	if errCount != nil {
		return nil, translateError(errCount)
	}

	var accounts []domain.Account
	if count == 0 {
		return &domain.AccountList{
			Metadata: domain.ListMetadata{
				Total:  0,
				Limit:  *params.Limit,
				Offset: *params.Offset,
			},
			Accounts: accounts,
		}, nil
	}
	query += " LIMIT? OFFSET?"
	args = append(args, params.Limit, params.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, translateError(err)
	}
	defer rows.Close()

	for rows.Next() {
		account := domain.Account{
			Owner: &domain.HouseholdMember{},
		}
		err := rows.Scan(
			&account.ID,
			&account.Name,
			&account.Type,
			&account.Institution,
			&account.Currency,
			&account.InitialBalance,
			&account.CurrentBalance,
			&account.Active,
			&account.Description,
			&account.AccountNumber,
			&account.AccountInformation,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.Owner.ID,
			&account.Owner.FirstName,
			&account.Owner.LastName,
			&account.Owner.Nickname,
			&account.Owner.Role,
			&account.Owner.Active,
			&account.Owner.CreatedAt,
			&account.Owner.UpdatedAt,
		)
		if err != nil {
			return nil, translateError(err)
		}
		accounts = append(accounts, account)
	}
	return &domain.AccountList{
			Metadata: domain.ListMetadata{
				Total:  count,
				Limit:  *params.Limit,
				Offset: *params.Offset,
			},
			Accounts: accounts,
		},
		nil
}

func (r *AccountRepoImpl) HasTransactions(ctx context.Context, id string) (bool, error) {
	query := `SELECT COUNT(*) FROM transactions WHERE account_id =?`
	var count int
	err := r.db.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return false, translateError(err)
	}
	return count > 0, nil
}
