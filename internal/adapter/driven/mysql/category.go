package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
)

type CategoryRepoImpl struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) port.CategoryRepo {
	return &CategoryRepoImpl{db: db}
}

func (c CategoryRepoImpl) Create(
	ctx context.Context,
	category domain.Category,
) (
	string,
	error,
) {
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
		category.CategoryType,
	)
	if errInsert != nil {
		return "", translateError(errInsert)
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf(
			"failed to get last insert ID: %w",
			err,
		)
	}
	lastIDStr := strconv.FormatInt(
		lastID,
		10,
	)

	return lastIDStr, nil
}

func (c CategoryRepoImpl) Update(
	ctx context.Context,
	category domain.Category,
) error {
	queryUpdate := `UPDATE categories SET name=?, description=?, color=?, background_color=?, active=?, category_type=? WHERE id=?`

	result, err := c.db.ExecContext(
		ctx,
		queryUpdate,
		category.Name,
		category.Description,
		category.Color,
		category.BackgroundColor,
		category.Active,
		category.CategoryType,
		category.ID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to update category: %w",
			err,
		)
	}
	rowsAffected, errRowsAffected := result.RowsAffected()
	if errRowsAffected != nil {
		return translateError(err)
	}
	if rowsAffected == 0 {
		return port.ErrRecordNotFound
	}

	return nil
}

func (c CategoryRepoImpl) Delete(
	ctx context.Context,
	id string,
) error {
	queryUpdate := `delete from categories where id=?`

	result, err := c.db.ExecContext(
		ctx,
		queryUpdate,
		id,
	)
	if err != nil {
		return translateError(err)
	}
	rowsAffected, errRowsAffected := result.RowsAffected()
	if errRowsAffected != nil {
		return translateError(err)
	}
	if rowsAffected == 0 {
		return port.ErrRecordNotFound
	}

	return nil
}

func (c CategoryRepoImpl) GetByID(
	ctx context.Context,
	id string,
) (
	*domain.Category,
	error,
) {
	query := `SELECT id, name, description, color, background_color, active, category_type FROM categories WHERE id=?`

	var category domain.Category
	err := c.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Color,
		&category.BackgroundColor,
		&category.Active,
		&category.CategoryType,
	)
	if errors.Is(
		err,
		sql.ErrNoRows,
	) {
		return nil, port.ErrRecordNotFound
	} else if err != nil {
		return nil, translateError(err)
	}

	return &category, nil
}

func (c CategoryRepoImpl) List(ctx context.Context) (
	[]domain.Category,
	error,
) {
	query := `SELECT id, name, description, color, background_color, active, category_type FROM categories WHERE active=true`

	rows, err := c.db.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to select categories: %w",
			err,
		)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.Color,
			&category.BackgroundColor,
			&category.Active,
			&category.CategoryType,
		)
		if err != nil {
			return nil, translateError(err)
		}
		categories = append(
			categories,
			category,
		)
	}

	return categories, nil
}

func (c CategoryRepoImpl) FindByType(
	ctx context.Context,
	categoryType domain.CategoryType,
) (
	[]domain.Category,
	error,
) {
	query := `SELECT id, name, description, color, background_color, active, category_type FROM categories WHERE category_type=? AND active=true`

	rows, err := c.db.QueryContext(
		ctx,
		query,
		categoryType,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to select categories by type: %w",
			err,
		)
	}
	defer rows.Close()
	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.Color,
			&category.BackgroundColor,
			&category.Active,
			&category.CategoryType,
		)
		if err != nil {
			return nil, translateError(err)
		}
		categories = append(
			categories,
			category,
		)
	}

	return categories, nil
}

func (c CategoryRepoImpl) FindByIDs(
	ctx context.Context,
	ids []string,
) (
	[]domain.Category,
	error,
) {
	query := `SELECT id, name, description, color, background_color, active FROM categories WHERE id IN (?) AND active=true`
	rows, err := c.db.QueryContext(
		ctx,
		query,
		ids,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to select categories by IDs: %w",
			err,
		)
	}
	defer rows.Close()
	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.Color,
			&category.BackgroundColor,
			&category.Active,
		)
		if err != nil {
			return nil, translateError(err)
		}
		categories = append(
			categories,
			category,
		)
	}

	return categories, nil
}
