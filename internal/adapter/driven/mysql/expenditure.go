package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	queryInsert := `INSERT INTO expenditures (category_id, declared, planned, created_at, updated_at) VALUES (?,?,?,NOW(),NOW())`
	result, errInsert := r.db.ExecContext(
		ctx,
		queryInsert,
		expenditure.Category.Id,
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
	queryUpdate := `UPDATE expenditures SET category_id =?,  declared =?, planned =?, updated_at = NOW() WHERE id =?`

	_, err := r.db.ExecContext(
		ctx,
		queryUpdate,
		expenditure.Category.Id,
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

	return nil

}

func (r *ExpenditureRepo) GetByID(ctx context.Context, id string) (*openapi.Expenditure, error) {
	querySelect := `select t.account_id,
						   t.amount,
						   e.category_id,
						   t.created_at as createdAt,
						   t.transaction_date as date,
						   e.declared,
						   t.description,
						   e.planned,
						   GROUP_CONCAT(et.tag_id ORDER BY et.tag_id SEPARATOR ',') as tags,
						   t.updated_at
					from expenditures e
							 inner join
						 transactions t
						 on e.transaction_id = t.id
							 left join
						 expenditure_tags et
						 on e.id = et.expenditure_id WHERE id =?
						 group by t.account_id, t.amount, e.category_id, t.created_at, t.transaction_date, e.declared, t.description, e.planned, t.updated_at `

	var expenditure openapi.Expenditure
	var categoryId int64
	var tags string
	err := r.db.QueryRowContext(ctx, querySelect, id).Scan(
		&expenditure.AccountId,
		&expenditure.Amount,
		&categoryId,
		&expenditure.CreatedAt,
		&expenditure.Date,
		&expenditure.Declared,
		&expenditure.Description,
		&expenditure.Planned,
		&expenditure.CreatedAt,
		&tags,
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
	querySelect := `select t.account_id,
						   t.amount,
						   e.category_id,
						   t.created_at as createdAt,
						   t.transaction_date as date,
						   e.declared,
						   t.description,
						   e.planned,
						   GROUP_CONCAT(et.tag_id ORDER BY et.tag_id SEPARATOR ',') as tags,
						   t.updated_at
					from expenditures e
							 inner join
						 transactions t
						 on e.transaction_id = t.id
							 left join
						 expenditure_tags et
						 on e.id = et.expenditure_id`
	queryCount := `select COUNT(*)
					from expenditures e
							 inner join
						 transactions t
						 on e.transaction_id = t.id
							 left join
						 expenditure_tags et
						 on e.id = et.expenditure_id
`
	var args []interface{}
	var whereClause []string

	if queryParams.CategoryId != nil {
		whereClause = append(whereClause, "e.category_id =?")
		args = append(args, *queryParams.CategoryId)
	}
	if queryParams.StartDate != nil && queryParams.EndDate != nil {
		whereClause = append(whereClause, "t.transaction_date BETWEEN? AND?")
	} else if queryParams.StartDate != nil {
		whereClause = append(whereClause, "t.transaction_date >=?")
		args = append(args, *queryParams.StartDate)
	} else if queryParams.EndDate != nil {
		whereClause = append(whereClause, "t.transaction_date <=?")
		args = append(args, *queryParams.EndDate)
	}
	if queryParams.Declared != nil {
		whereClause = append(whereClause, "e.declared =?")
	}
	if queryParams.Planned != nil {
		whereClause = append(whereClause, "e.planned =?")
	}
	if queryParams.Currency != nil {
		whereClause = append(whereClause, "t.currency =?")
	}
	if queryParams.Description != nil {
		whereClause = append(whereClause, "t.description LIKE ?")
		args = append(args, "%"+*queryParams.Description+"%")
	}
	if queryParams.Tags != nil {
		tagList := make([]string, 0, len(*queryParams.Tags))
		for _, tag := range *queryParams.Tags {
			tagList = append(tagList, fmt.Sprintf("'%s'", tag))
		}
		whereClause = append(whereClause, fmt.Sprintf(` EXISTS (
				SELECT 1
				FROM expenditures_expenditure_tags et2
				WHERE et2.expenditure_id = e.id AND et2.tag_id IN %s
			) `, strings.Join(tagList, ", ")))
	}
	if queryParams.AccountId != nil {
		// inner select to transaction_expenditures and then to transactions table to get account_id
		whereClause = append(whereClause, "t.account_id =?")
	}
	if len(whereClause) > 0 {
		querySelect += " WHERE "
		queryCount += " WHERE "
	}
	for i, clause := range whereClause {
		if i > 0 {
			querySelect += " AND "
			queryCount += " AND "
		}
		querySelect += clause
		queryCount += clause
	}
	querySelect += " group by t.account_id, t.amount, e.category_id, t.created_at, t.transaction_date, e.declared, t.description, e.id, e.planned, t.updated_at "
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
		var categoryId int64
		var tags string
		err := rows.Scan(
			&expenditure.AccountId,
			&expenditure.Amount,
			&categoryId,
			&expenditure.CreatedAt,
			&expenditure.Date,
			&expenditure.Declared,
			&expenditure.Description,
			&expenditure.Planned,
			&expenditure.CreatedAt,
			&tags,
			&expenditure.UpdatedAt,
		)
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
