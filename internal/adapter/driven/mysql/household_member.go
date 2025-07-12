package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"strconv"
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
		return "", fmt.Errorf("failed to create household member: %w", err)
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
		return fmt.Errorf("failed to update household member: %w", err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err // TODO return correct error from domain
		}
		return nil, fmt.Errorf("failed to select household member: %w", err)
	}
	return &householdMember, nil
}

func (h HouseholdMemberRepository) CanDelete(ctx context.Context, id string) (bool, error) {
	query := `SELECT COUNT(*) FROM accounts WHERE owner = ?`
	var count int
	err := h.db.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check household member deletion: %w", err)
	}
	return count == 0, nil
}

func (h HouseholdMemberRepository) List(ctx context.Context) (*domain.HouseholdMemberList, error) {
	query := `SELECT id, name, surname, nickname, role, active, created_at, updated_at FROM household_members WHERE active = true`
	rows, err := h.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to select household members: %w", err)
	}
	defer rows.Close()
	var householdMembers []domain.HouseholdMember
	for rows.Next() {
		var householdMember domain.HouseholdMember
		err := rows.Scan(&householdMember.ID, &householdMember.FirstName, &householdMember.LastName, &householdMember.Nickname, &householdMember.Role, &householdMember.Active, &householdMember.CreatedAt, &householdMember.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		householdMembers = append(householdMembers, householdMember)
	}
	return &domain.HouseholdMemberList{HouseholdMembers: householdMembers}, nil
}
