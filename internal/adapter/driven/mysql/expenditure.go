package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/common"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
	"strconv"
	"strings"
)

type ExpenditureRepo struct {
	db *sql.DB
}

func NewExpenditureRepo(db *sql.DB) *ExpenditureRepo {
	return &ExpenditureRepo{db: db}
}

func (r *ExpenditureRepo) Create(ctx context.Context, expenditure openapi.Expenditure) (string, error) {
	queryInsert := `INSERT INTO expenditures (category_id, date, declared, planned, created_at, updated_at) VALUES (?,?,?,NOW(),NOW())`
	result, errInsert := r.db.ExecContext(
		ctx,
		queryInsert,
		expenditure.CategoryId,
		expenditure.Date,
		expenditure.Declared,
		expenditure.Planned,
	)
	if errInsert != nil {
		return "", fmt.Errorf("failed to insert expenditure: %w", errInsert)
	}
	expenditureId, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("failed to get last insert ID: %w", err)
	}
	lastIDStr := strconv.FormatInt(expenditureId, 10)
	return lastIDStr, nil
}

func (r *ExpenditureRepo) Update(ctx context.Context, id string, expenditure openapi.Expenditure) error {
	queryUpdate := `UPDATE expenditures SET category_id =?, date =?, declared =?, planned =?, updated_at = NOW() WHERE id =?`

	_, err := r.db.ExecContext(
		ctx,
		queryUpdate,
		expenditure.CategoryId,
		expenditure.Date,
		expenditure.Declared,
		expenditure.Planned,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update expenditure: %w", err)
	}

	return nil

}

func (r *ExpenditureRepo) Delete(ctx context.Context, id string) error {
	queryDelete := `DELETE FROM expenditures WHERE id =?`

	_, err := r.db.ExecContext(ctx, queryDelete, id)
	if err != nil {
		return fmt.Errorf("failed to delete expenditure: %w", err)
	}

	clearTagsQuery := `DELETE FROM expenditures_expenditure_tags WHERE expenditure_id =?`
	_, errClearTags := r.db.Exec(clearTagsQuery, id)
	if errClearTags != nil {
		return fmt.Errorf("failed to clear expenditure tags: %w", errClearTags)
	}

	return nil

}

func (r *ExpenditureRepo) GetByID(ctx context.Context, id string) (*openapi.Expenditure, error) {
	querySelect := `SELECT id, category_id, date, declared, planned, created_at, updated_at FROM expenditures WHERE id =?`

	var expenditure openapi.Expenditure
	err := r.db.QueryRowContext(ctx, querySelect, id).Scan(
		&expenditure.Id,
		&expenditure.CategoryId,
		&expenditure.Date,
		&expenditure.Declared,
		&expenditure.Planned,
		&expenditure.CreatedAt,
		&expenditure.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrEntityNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to select expenditure: %w", err)
	}

	return &expenditure, nil
}

//FindExpenditures
/*
 * FindExpenditures returns a list of expenditures based on the provided query parameters.
 *
 * @param queryParams openapi.ListExpendituresParams The query parameters for filtering and sorting expenditures.
 * @return *openapi.ExpenditureList The list of expenditures.
 * @return error If an error occurs while fetching the expenditures.

 * Notes:
	- sort is default by creation of the expenditure.
    - filter is default by not filtering any expenditures.
    - query is built concatenating conditions with AND.
*/
func (r *ExpenditureRepo) FindExpenditures(ctx context.Context, queryParams openapi.ListExpendituresParams) (*openapi.ExpenditureList, error) {
	querySelect := `SELECT id, category_id, date, declared, planned, created_at, updated_at FROM expenditures WHERE %s`
	queryCount := `SELECT COUNT(*) FROM expenditures WHERE %s`

	var args []interface{}
	var whereClause []string

	if queryParams.CategoryId != nil {
		whereClause = append(whereClause, "category_id =?")
		args = append(args, *queryParams.CategoryId)
	}
	if queryParams.StartDate != nil && queryParams.EndDate != nil {
		whereClause = append(whereClause, "date BETWEEN? AND?")
	} else if queryParams.StartDate != nil {
		whereClause = append(whereClause, "date >=?")
		args = append(args, *queryParams.StartDate)
	} else if queryParams.EndDate != nil {
		whereClause = append(whereClause, "date <=?")
		args = append(args, *queryParams.EndDate)
	}
	if queryParams.Declared != nil {
		whereClause = append(whereClause, "declared =?")
	}
	if queryParams.Planned != nil {
		whereClause = append(whereClause, "planned =?")
	}
	if queryParams.Currency != nil {
		whereClause = append(whereClause, "currency =?")
	}
	if queryParams.Description != nil {
		whereClause = append(whereClause, "description LIKE ?")
		args = append(args, "%"+*queryParams.Description+"%")
	}
	if queryParams.Tags != nil {
		tagList := make([]string, 0, len(*queryParams.Tags))
		for _, tag := range *queryParams.Tags {
			tagList = append(tagList, fmt.Sprintf("'%s'", tag))
		}
		whereClause = append(whereClause, fmt.Sprintf("id IN (SELECT expenditure_id FROM expenditures_expenditure_tags WHERE tag_id IN (%s))", strings.Join(tagList, ", ")))
	}
	if queryParams.AccountId != nil {
		// inner select to transaction_expenditures and then to transactions table to get account_id
		querySelect = fmt.Sprintf("%s INNER JOIN transactions ON expenditures.id = transactions.expenditure_id AND transactions.account_id =?", querySelect)
		queryCount = fmt.Sprintf("%s INNER JOIN transactions ON expenditures.id = transactions.expenditure_id AND transactions.account_id =?", queryCount)
	}
	if len(whereClause) == 0 {
		// Removing where from queries
		querySelect = "SELECT id, category_id, date, declared, planned, created_at, updated_at FROM expenditures"
		queryCount = "SELECT COUNT(*) FROM expenditures"
	}
	for i, clause := range whereClause {
		if i > 0 {
			querySelect += " AND "
			queryCount += " AND "
		}
		querySelect = fmt.Sprintf("%s %s", querySelect, clause)
		queryCount = fmt.Sprintf("%s %s", queryCount, clause)
	}
	querySelect += " ORDER BY created_at DESC"
	querySelect += fmt.Sprintf(" LIMIT %d OFFSET %d", queryParams.Limit, queryParams.Offset)
	stmtCount, errQueryCountStmt := r.db.PrepareContext(ctx, queryCount)
	if errQueryCountStmt != nil {
		return nil, fmt.Errorf("failed to prepare count statement: %w", errQueryCountStmt)
	}
	defer stmtCount.Close()
	result, errQueryCount := stmtCount.QueryContext(ctx, args...)
	if errQueryCount != nil {
		return nil, fmt.Errorf("failed to count expenditures: %w", errQueryCount)
	}
	defer result.Close()
	var count float32
	result.Next()
	errScan := result.Scan(&count)
	if errScan != nil {
		return nil, errScan
	}

	stmntQuery, errStmtQuerySelect := r.db.PrepareContext(ctx, querySelect)
	if errStmtQuerySelect != nil {
		return nil, fmt.Errorf("failed to prepare select statement: %w", errStmtQuerySelect)
	}

	defer stmntQuery.Close()
	rows, errQuery := stmntQuery.QueryContext(ctx, args...)
	if errQuery != nil {
		return nil, fmt.Errorf("failed to select expenditures: %w", errQuery)
	}

	defer rows.Close()
	var expenditures []openapi.Expenditure
	for rows.Next() {
		var expenditure openapi.Expenditure
		err := rows.Scan(&expenditure.Id, &expenditure.CategoryId, &expenditure.Date, &expenditure.Declared, &expenditure.Planned, &expenditure.CreatedAt, &expenditure.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		expenditures = append(expenditures, expenditure)
	}

	expendituresList := &openapi.ExpenditureList{
		Total:        &count,
		Expenditures: &expenditures,
	}

	return expendituresList, nil
}

func (r *ExpenditureRepo) FindOrCreateTags(ctx context.Context, tags []string) ([]string, error) {
	//Generate random color and background color for each tag in hex format

	query := `INSERT INTO expenditure_tags (name, color, background_color, description) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE id=id`

	insertedTagsStr := make([]string, 0, len(tags))
	for _, tag := range tags {
		colors, errColors := common.GetPillColor(tag)
		if errColors != nil {
			return nil, fmt.Errorf("failed to get tag colors: %w", errColors)
		}
		var tagId int64
		errQuery := r.db.QueryRowContext(ctx, query, tag, (*colors)[0], (*colors)[1], tag).Scan(&tagId)
		if errQuery != nil {
			return nil, fmt.Errorf("failed to insert or find tag: %w", errQuery)
		}
		insertedTagsStr = append(insertedTagsStr, strconv.FormatInt(tagId, 10))
	}
	return insertedTagsStr, nil
}

func (r *ExpenditureRepo) LinkTagsToExpenditure(ctx context.Context, tags []string, expenditureId string) error {
	query := `INSERT INTO expenditures_expenditure_tags (expenditure_id, tag_id) VALUES (?,?)`

	for _, tagId := range tags {
		_, err := r.db.ExecContext(ctx, query, expenditureId, tagId)
		if err != nil {
			return fmt.Errorf("failed to link tags to expenditure: %w", err)
		}
	}
	return nil
}

func (r *ExpenditureRepo) ListTags(ctx context.Context) ([]string, error) {
	query := `SELECT name FROM expenditure_tags`

	rows, errQuery := r.db.QueryContext(ctx, query)
	if errQuery != nil {
		return nil, fmt.Errorf("failed to select tags: %w", errQuery)
	}

	defer rows.Close()
	var tags []string
	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *ExpenditureRepo) FindCategory(ctx context.Context, id string) (*openapi.ExpenditureCategory, error) {
	query := `SELECT id, name, description FROM expenditures_categories WHERE id =?`

	var category openapi.ExpenditureCategory
	errQuery := r.db.QueryRowContext(ctx, query, id).Scan(&category.Id, &category.Name, &category.Description)
	if errQuery != nil {
		if errors.Is(errQuery, sql.ErrNoRows) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to select category: %w", errQuery)
	}

	return &category, nil
}

func (r *ExpenditureRepo) ListCategories(ctx context.Context) ([]openapi.ExpenditureCategory, error) {
	query := `SELECT id, name FROM expenditures_categories`

	rows, errQuery := r.db.QueryContext(ctx, query)
	if errQuery != nil {
		return nil, fmt.Errorf("failed to select categories: %w", errQuery)
	}

	defer rows.Close()
	var categories []string
	for rows.Next() {
		var categoryId, categoryName string
		err := rows.Scan(&categoryId, &categoryName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		categories = append(categories, fmt.Sprintf("%s (%s)", categoryId, categoryName))
	}

	return categories, nil
}

func (r *ExpenditureRepo) CreateCategory(ctx context.Context, name, description string) (string, error) {
	query := `INSERT INTO expenditures_categories (name, description) VALUES (?,?) RETURNING id`

	var categoryId int64
	errQuery := r.db.QueryRowContext(ctx, query, name, description).Scan(&categoryId)
	if errQuery != nil {
		return "", fmt.Errorf("failed to create category: %w", errQuery)
	}

	return strconv.FormatInt(categoryId, 10), nil
}

func (r *ExpenditureRepo) UpdateCategory(ctx context.Context, id string, name, description string) error {
	query := `UPDATE expenditures_categories SET name =?, description =? WHERE id =?`

	_, err := r.db.ExecContext(ctx, query, name, description, id)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (r *ExpenditureRepo) DeleteCategory(ctx context.Context, id string) error {
	query := `DELETE FROM expenditures_categories WHERE id =?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}
