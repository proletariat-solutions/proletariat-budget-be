package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type AccountRepo struct {
	db *sql.DB
}

func NewAccountRepo(db *sql.DB) *AccountRepo {
	return &AccountRepo{db: db}
}

func (r *AccountRepo) Create(ctx context.Context, account openapi.Account) (string, error) {

	now := time.Now()
	account.CreatedAt = now
	account.UpdatedAt = now

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
		account.CurrentBalance,
		account.Active,
		account.Description,
		account.AccountNumber,
		account.Owner.Id,
		account.AccountInformation,
		account.CreatedAt,
		account.UpdatedAt,
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

func (r *AccountRepo) GetByID(ctx context.Context, id string) (*openapi.Account, error) {
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

func (r *AccountRepo) Update(ctx context.Context, account openapi.Account) error {
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

func (r *AccountRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM accounts WHERE id =?`

	result, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)

	if err != nil {
		log.Error().Err(err).Msg("Failed to delete account")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("account not found")
	}
	return nil
}

func (r *AccountRepo) List(ctx context.Context) ([]openapi.Account, error) {
	query := `SELECT 
                a.id, a.name, type, institution, currency, initial_balance, current_balance, a.active, 
                description, account_number, account_information, a.created_at, a.updated_at,  hm.id, hm.name, hm.surname, hm.nickname, hm.role, hm.active, hm.created_at, hm.updated_at
                FROM accounts a left join household_members hm on a.owner = hm.id WHERE a.active = true`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to select accounts: %w", err)
	}
	defer rows.Close()

	var accounts []openapi.Account
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
	return accounts, nil
}
