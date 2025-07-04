package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
)

type CategoryRepoImpl struct {
	db *sql.DB
}

func (c CategoryRepoImpl) Create(ctx context.Context, category openapi.Category, categoryType string) (string, error) {
	queryInsert := `INSERT INTO categories  (name, description, color, background_color, active, category_type) 
						VALUES (?,?,?,?,?,?)`
	result, errInsert := c.db.ExecContext(
		ctx,
		queryInsert,
		category.Name,
		category.Description,
		category.Color,
		category.BackgroundColor,
		category.Active,
		categoryType,
	)
	if errInsert != nil {
		return "", fmt.Errorf("failed to create category: %w", errInsert)
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("failed to get last insert ID: %w", err)
	}
	lastIDStr := fmt.Sprintf("%d", lastID)
	return lastIDStr, nil
}

func (c CategoryRepoImpl) Update(ctx context.Context, id string, category openapi.Category, categoryType string) error {
	queryUpdate := `UPDATE categories SET name=?, description=?, color=?, background_color=?, active=?, category_type=? WHERE id=?`

	result, err := c.db.ExecContext(
		ctx, queryUpdate,
		category.Name,
		category.Description,
		category.Color,
		category.BackgroundColor,
		category.Active,
		categoryType,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	rowsAffected, errRowsAffected := result.RowsAffected()
	if errRowsAffected != nil {
		return fmt.Errorf("failed to update category: %w", errRowsAffected)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

func (c CategoryRepoImpl) Delete(ctx context.Context, id string) error {
	// Only updating "active" field to false, not deleting the record
	queryUpdate := `UPDATE categories SET active = false WHERE id =?`

	result, err := c.db.ExecContext(
		ctx, queryUpdate, id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	rowsAffected, errRowsAffected := result.RowsAffected()
	if errRowsAffected != nil {
		return fmt.Errorf("failed to delete category: %w", errRowsAffected)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

func (c CategoryRepoImpl) GetByID(ctx context.Context, id string) (*openapi.Category, error) {
	query := `SELECT id, name, description, color, background_color, active FROM categories WHERE id=? AND active=true`

	var category openapi.Category
	err := c.db.QueryRowContext(ctx, query, id).Scan(
		&category.Id,
		&category.Name,
		&category.Description,
		&category.Color,
		&category.BackgroundColor,
		&category.Active,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrEntityNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to select category: %w", err)
	}
	return &category, nil
}

func (c CategoryRepoImpl) List(ctx context.Context) ([]openapi.Category, error) {
	query := `SELECT id, name, description, color, background_color, active FROM categories WHERE active=true`

	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to select categories: %w", err)
	}
	defer rows.Close()

	var categories []openapi.Category
	for rows.Next() {
		var category openapi.Category
		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.Description,
			&category.Color,
			&category.BackgroundColor,
			&category.Active,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (c CategoryRepoImpl) FindByType(ctx context.Context, categoryType string) ([]openapi.Category, error) {
	query := `SELECT id, name, description, color, background_color, active FROM categories WHERE category_type=? AND active=true`

	rows, err := c.db.QueryContext(ctx, query, categoryType)
	if err != nil {
		return nil, fmt.Errorf("failed to select categories by type: %w", err)
	}
	defer rows.Close()
	var categories []openapi.Category
	for rows.Next() {
		var category openapi.Category
		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.Description,
			&category.Color,
			&category.BackgroundColor,
			&category.Active,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (c CategoryRepoImpl) FindByIDs(ctx context.Context, ids []string) ([]openapi.Category, error) {
	query := `SELECT id, name, description, color, background_color, active FROM categories WHERE id IN (?) AND active=true`
	rows, err := c.db.QueryContext(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to select categories by IDs: %w", err)
	}
	defer rows.Close()
	var categories []openapi.Category
	for rows.Next() {
		var category openapi.Category
		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.Description,
			&category.Color,
			&category.BackgroundColor,
			&category.Active,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}
