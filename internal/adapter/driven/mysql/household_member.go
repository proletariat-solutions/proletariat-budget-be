package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
	"strconv"
)

type HouseholdMemberRepository struct {
	db *sql.DB
}

func (h HouseholdMemberRepository) Create(ctx context.Context, householdMember openapi.HouseholdMember) (string, error) {
	query := `INSERT INTO household_members (name, surname, nickname, role, active, created_at, updated_at) VALUES (?,?,?,?,true, now(), NOW())`
	result, err := h.db.ExecContext(
		ctx, query, householdMember.FirstName, householdMember.LastName, householdMember.Nickname, householdMember.Role,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create household member: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("failed to get last inserted id: %w", err)
	}
	return strconv.FormatInt(id, 10), nil
}

func (h HouseholdMemberRepository) Update(ctx context.Context, householdMember openapi.HouseholdMember) error {
	query := `UPDATE household_members SET name =?, surname =?, nickname =?, role =?, updated_at = NOW() WHERE id =?`
	result, err := h.db.ExecContext(
		ctx, query, householdMember.FirstName, householdMember.LastName, householdMember.Nickname, householdMember.Role, householdMember.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to update household member: %w", err)
	}
	rowsAffected, errRowsAffected := result.RowsAffected()
	if errRowsAffected != nil {
		return fmt.Errorf("failed to update household member: %w", errRowsAffected)
	}
	if rowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func (h HouseholdMemberRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE household_members SET active = false, updated_at = NOW() WHERE id =?`
	result, err := h.db.ExecContext(
		ctx, query, id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete household member: %w", err)
	}
	rowsAffected, errRowsAffected := result.RowsAffected()
	if errRowsAffected != nil {
		return fmt.Errorf("failed to delete household member: %w", errRowsAffected)
	}
	if rowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func (h HouseholdMemberRepository) GetByID(ctx context.Context, id string) (*openapi.HouseholdMember, error) {
	query := `SELECT id, name, surname, nickname, role, active, created_at, updated_at FROM household_members WHERE id =?`
	var householdMember openapi.HouseholdMember
	row := h.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&householdMember.Id, &householdMember.FirstName, &householdMember.LastName, &householdMember.Nickname, &householdMember.Role, &householdMember.Active, &householdMember.CreatedAt, &householdMember.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to select household member: %w", err)
	}
	return &householdMember, nil
}

func (h HouseholdMemberRepository) List(ctx context.Context) ([]openapi.HouseholdMember, error) {
	query := `SELECT id, name, surname, nickname, role, active, created_at, updated_at FROM household_members WHERE active = true`
	rows, err := h.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to select household members: %w", err)
	}
	defer rows.Close()
	var householdMembers []openapi.HouseholdMember
	for rows.Next() {
		var householdMember openapi.HouseholdMember
		err := rows.Scan(&householdMember.Id, &householdMember.FirstName, &householdMember.LastName, &householdMember.Nickname, &householdMember.Role, &householdMember.Active, &householdMember.CreatedAt, &householdMember.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		householdMembers = append(householdMembers, householdMember)
	}
	return householdMembers, nil
}
