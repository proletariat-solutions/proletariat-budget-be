package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"strconv"
	"strings"
)

type HouseholdMemberRepository struct {
	db *sql.DB
}

func NewHouseholdMemberRepository(db *sql.DB) port.HouseholdMembersRepo {
	return &HouseholdMemberRepository{db: db}
}

func (h HouseholdMemberRepository) Create(ctx context.Context, householdMember domain.HouseholdMember) (string, error) {
	query := `INSERT INTO household_members (name, surname, nickname, role, active, created_at, updated_at) VALUES (?,?,?,?,true, now(), NOW())`
	result, err := h.db.ExecContext(
		ctx, query, householdMember.FirstName, householdMember.LastName, householdMember.Nickname, householdMember.Role,
	)
	if err != nil {
		return "", translateError(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("failed to get last inserted id: %w", err)
	}
	return strconv.FormatInt(id, 10), nil
}

func (h HouseholdMemberRepository) Update(ctx context.Context, id string, householdMember domain.HouseholdMember) error {
	query := `UPDATE household_members SET name =?, surname =?, nickname =?, role =?, updated_at = NOW() WHERE id =?`
	result, err := h.db.ExecContext(
		ctx, query, householdMember.FirstName, householdMember.LastName, householdMember.Nickname, householdMember.Role, id,
	)
	if err != nil {
		return translateError(err)
	}
	rowsAffected, errRowsAffected := result.RowsAffected()
	if errRowsAffected != nil {
		return fmt.Errorf("failed to update household member: %w", errRowsAffected)
	}
	if rowsAffected == 0 {
		return err // TODO return correct error from domain
	}
	return nil
}

func (h HouseholdMemberRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM household_members WHERE id =?`
	result, err := h.db.ExecContext(
		ctx, query, id,
	)
	if err != nil {
		return translateError(err)
	}
	rowsAffected, errRowsAffected := result.RowsAffected()
	if errRowsAffected != nil {
		return fmt.Errorf("failed to delete household member: %w", errRowsAffected)
	}
	if rowsAffected == 0 {
		return port.ErrRecordNotFound
	}
	return nil
}

func (h HouseholdMemberRepository) Deactivate(ctx context.Context, id string) error {
	query := `UPDATE household_members SET active = false, updated_at = NOW() WHERE id =?`
	result, err := h.db.ExecContext(
		ctx, query, id,
	)
	if err != nil {
		return translateError(err)
	}
	rowsAffected, errRowsAffected := result.RowsAffected()
	if errRowsAffected != nil {
		return fmt.Errorf("failed to delete household member: %w", errRowsAffected)
	}
	if rowsAffected == 0 {
		return err // TODO return correct error from domain
	}
	return nil
}

func (h HouseholdMemberRepository) Activate(ctx context.Context, id string) error {
	query := `UPDATE household_members SET active = true, updated_at = NOW() WHERE id =?`
	result, err := h.db.ExecContext(
		ctx, query, id,
	)
	if err != nil {
		return translateError(err)
	}
	rowsAffected, errRowsAffected := result.RowsAffected()
	if errRowsAffected != nil {
		return fmt.Errorf("failed to delete household member: %w", errRowsAffected)
	}
	if rowsAffected == 0 {
		return err // TODO return correct error from domain
	}
	return nil
}

func (h HouseholdMemberRepository) GetByID(ctx context.Context, id string) (*domain.HouseholdMember, error) {
	query := `SELECT id, name, surname, nickname, role, active, created_at, updated_at FROM household_members WHERE id =?`
	var householdMember domain.HouseholdMember
	row := h.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&householdMember.ID, &householdMember.FirstName, &householdMember.LastName, &householdMember.Nickname, &householdMember.Role, &householdMember.Active, &householdMember.CreatedAt, &householdMember.UpdatedAt)
	if err != nil {
		return nil, translateError(err)
	}
	return &householdMember, nil
}

func (h HouseholdMemberRepository) CanDelete(ctx context.Context, id string) (bool, error) {
	query := `SELECT COUNT(*) FROM accounts WHERE owner = ? AND active = true`
	var count int
	err := h.db.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return false, translateError(err)
	}
	return count == 0, nil
}

func (h HouseholdMemberRepository) List(ctx context.Context, params *domain.HouseholdMemberListParams) (*domain.HouseholdMemberList, error) {
	query := `SELECT id, name, surname, nickname, role, active, created_at, updated_at FROM household_members`

	whereClause := make([]string, 0)
	args := make([]any, 0)
	if params.Active != nil {
		whereClause = append(whereClause, "active =?")
		args = append(args, *params.Active)
	}
	if params.Role != nil {
		whereClause = append(whereClause, "role =?")
		args = append(args, *params.Role)
	}
	if len(whereClause) > 0 {
		query += " WHERE " + strings.Join(whereClause, " AND ")
	}
	rows, err := h.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, translateError(err)
	}
	defer rows.Close()
	var householdMembers []domain.HouseholdMember
	for rows.Next() {
		var householdMember domain.HouseholdMember
		err = rows.Scan(&householdMember.ID, &householdMember.FirstName, &householdMember.LastName, &householdMember.Nickname, &householdMember.Role, &householdMember.Active, &householdMember.CreatedAt, &householdMember.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		householdMembers = append(householdMembers, householdMember)
	}
	return &domain.HouseholdMemberList{HouseholdMembers: householdMembers}, nil
}
