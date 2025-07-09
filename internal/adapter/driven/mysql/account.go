package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"ghorkov32/proletariat-budget-be/openapi"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type AccountRepoImpl struct {
	db *sql.DB
}

func NewAccountRepo(db *sql.DB) port.AccountRepo {
	return &AccountRepoImpl{db: db}
}

func (r *AccountRepoImpl) Create(ctx context.Context, account openapi.AccountRequest) (string, error) {
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
		account.Owner.Id,
		account.AccountInformation,
	)

	if err != nil {
		log.Error().Err(err).Msg("Failed to create account")
		return "", err
	}

	lastID, err := result.LastInsertId()

	if err != nil {
		log.Error().Err(err).Msg("Failed to get last insert ID")
		return "", err
	}

	lastIDStr := strconv.FormatInt(lastID, 10)

	return lastIDStr, nil
}

func (r *AccountRepoImpl) GetByID(ctx context.Context, id string) (*openapi.Account, error) {
	query := `SELECT 
				a.id, a.name, type, institution, currency, 
				initial_balance, current_balance, a.active, 
				description, account_number, account_information, 
				a.created_at, a.updated_at, hm.id, hm.name, 
				hm.surname, hm.nickname, hm.role, hm.active, 
				hm.created_at, hm.updated_at 
				FROM accounts a left join household_members hm on a.owner = hm.id WHERE a.active = true AND a.id =?`

	var account openapi.Account
	householdMember := openapi.HouseholdMember{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.Id,
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
		&householdMember.Id,
		&householdMember.FirstName,
		&householdMember.LastName,
		&householdMember.Nickname,
		&householdMember.Role,
		&householdMember.Active,
		&householdMember.CreatedAt,
		&householdMember.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrEntityNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to select account: %w", err)
	}

	return &account, nil

}

func (r *AccountRepoImpl) Update(ctx context.Context, account openapi.Account) error {
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
		account.Id,
		account.Owner.Id,
	)

	if err != nil {
		log.Error().Err(err).Msg("Failed to update account")
		return err
	}
	return nil
}

func (r *AccountRepoImpl) Deactivate(ctx context.Context, id string) error {
	query := `UPDATE accounts SET active = false, updated_at = NOW() WHERE id =?`
	result, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}
	rowsAffected, errRowsAffected := result.RowsAffected()
	if errRowsAffected != nil {
		return fmt.Errorf("failed to delete account: %w", errRowsAffected)
	}
	if rowsAffected == 0 {
		return domain.ErrEntityNotFound
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
		return fmt.Errorf("failed to delete account: %w", err)
	}
	rowsAffected, errRowsAffected := result.RowsAffected()
	if errRowsAffected != nil {
		return fmt.Errorf("failed to delete account: %w", errRowsAffected)
	}
	if rowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func (r *AccountRepoImpl) List(ctx context.Context, params openapi.ListAccountsParams) (*openapi.AccountList, error) {
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
	if len(whereClause) > 0 {
		query += " WHERE " + strings.Join(whereClause, " AND ")
	}
	query += " ORDER BY a.created_at DESC"
	queryCount := "SELECT COUNT(*) FROM accounts a where " + strings.Join(whereClause, " AND ")
	stmtCount, errQueryCountStmt := r.db.PrepareContext(ctx, queryCount)
	if errQueryCountStmt != nil {
		return nil, fmt.Errorf("failed to prepare count statement: %w", errQueryCountStmt)
	}
	var count int
	errCount := stmtCount.QueryRowContext(ctx, args...).Scan(&count)
	if errCount != nil {
		return nil, fmt.Errorf("failed to count rows: %w", errCount)
	}

	var accounts []openapi.Account
	if count == 0 {
		return &openapi.AccountList{
			Metadata: &openapi.ListMetadata{
				Total: 0, Limit: *params.Limit, Offset: *params.Offset,
			}, Accounts: &accounts,
		}, nil
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select accounts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		householdMember := openapi.HouseholdMember{}
		var account openapi.Account
		err := rows.Scan(
			&account.Id,
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
			&householdMember.Id,
			&householdMember.FirstName,
			&householdMember.LastName,
			&householdMember.Nickname,
			&householdMember.Role,
			&householdMember.Active,
			&householdMember.CreatedAt,
			&householdMember.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		accounts = append(accounts, account)
	}
	return &openapi.AccountList{
		Metadata: &openapi.ListMetadata{
			Total: count, Limit: *params.Limit, Offset: *params.Offset,
		}, Accounts: &accounts}, nil
}

func (r *AccountRepoImpl) HasTransactions(ctx context.Context, id string) (bool, error) {
	query := `SELECT COUNT(*) FROM transactions WHERE account_id =?`
	var count int
	err := r.db.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check for transactions: %w", err)
	}
	return count > 0, nil
}
