package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"ghorkov32/proletariat-budget-be/openapi"
	"sort"
	"strconv"
	"strings"
)

type ExpenditureRepo struct {
	db       *sql.DB
	tagsRepo *port.TagsRepo
}

func NewExpenditureRepo(db *sql.DB, tagsRepo *port.TagsRepo) *ExpenditureRepo {
	return &ExpenditureRepo{
		db:       db,
		tagsRepo: tagsRepo,
	}
}

func (r *ExpenditureRepo) Create(ctx context.Context, expenditure openapi.Expenditure, transactionID string) (string, error) {
	queryInsert := `insert into expenditures
						(category_id, declared, planned, transaction_id, created_at, updated_at)
					VALUES (?, ?, ?, ?, ?, NOW())`
	result, errInsert := r.db.ExecContext(
		ctx,
		queryInsert,
		expenditure.Category.Id,
		expenditure.Declared,
		expenditure.Planned,
		transactionID,
		expenditure.Date,
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
	querySelect := `select e.id,
						   e.declared,
						   e.planned,
						   e.created_at,
						   t.account_id,
						   t.amount,
						   t.currency,
						   t.transaction_date,
						   t.description,
						   c.id,
						   c.name,
						   c.description,
						   c.color,
						   c.background_color,
						   c.active,
						   GROUP_CONCAT(et.tag_id ORDER BY et.tag_id SEPARATOR ',') as tags
					from expenditures e
							 inner join categories c ON e.category_id = c.id
							 inner join transactions t ON e.transaction_id = t.id
							 left join expenditure_tags et on e.id = et.expenditure_id
					group by e.id, e.declared, e.planned, e.transaction_id, e.created_at, t.account_id, t.amount, t.currency,
							 t.transaction_date, t.description, c.id, c.name, c.description, c.color, c.background_color, c.active,
							 c.category_type`

	var expenditure openapi.Expenditure
	var tagsList string
	err := r.db.QueryRowContext(ctx, querySelect, id).Scan(
		&expenditure.Id,
		&expenditure.Declared,
		&expenditure.Planned,
		&expenditure.CreatedAt,
		&expenditure.AccountId,
		&expenditure.Amount,
		&expenditure.Currency,
		&expenditure.Date,
		&expenditure.Description,
		&expenditure.Category.Id,
		&expenditure.Category.Name,
		&expenditure.Category.Description,
		&expenditure.Category.Color,
		&expenditure.Category.BackgroundColor,
		&expenditure.Category.Active,
		&tagsList,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrEntityNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to select expenditure: %w", err)
	}

	expenditure.Tags, err = (*r.tagsRepo).GetByIDs(ctx, strings.Split(tagsList, ","))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
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
	querySelect := `select e.id,
						   e.declared,
						   e.planned,
						   e.created_at,
						   t.account_id,
						   t.amount,
						   t.currency,
						   t.transaction_date,
						   t.description,
						   c.id,
						   c.name,
						   c.description,
						   c.color,
						   c.background_color,
						   c.active,
						   GROUP_CONCAT(et.tag_id ORDER BY et.tag_id SEPARATOR ',') as tags
					from expenditures e
							 inner join categories c ON e.category_id = c.id
							 inner join transactions t ON e.transaction_id = t.id
							 left join expenditure_tags et on e.id = et.expenditure_id`
	queryCount := `select COUNT(*)
						from expenditures e
							 inner join categories c ON e.category_id = c.id
							 inner join transactions t ON e.transaction_id = t.id
							 left join expenditure_tags et on e.id = et.expenditure_id`

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
	querySelect += ` group by e.id, e.declared, e.planned, e.transaction_id, e.created_at, t.account_id, t.amount, t.currency,
							 t.transaction_date, t.description, c.id, c.name, c.description, c.color, c.background_color, c.active,
							 c.category_type`
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
	var count int
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
	tagsByID := make(map[string][]string)
	var ids []string
	for rows.Next() {
		var expenditure openapi.Expenditure
		var tags string
		err := rows.Scan(
			&expenditure.Id,
			&expenditure.Declared,
			&expenditure.Planned,
			&expenditure.CreatedAt,
			&expenditure.AccountId,
			&expenditure.Amount,
			&expenditure.Currency,
			&expenditure.Date,
			&expenditure.Description,
			&expenditure.Category.Id,
			&expenditure.Category.Name,
			&expenditure.Category.Description,
			&expenditure.Category.Color,
			&expenditure.Category.BackgroundColor,
			&expenditure.Category.Active,
			&tags,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		expenditures = append(expenditures, expenditure)
		tagsByID[expenditure.Id] = strings.Split(tags, ",")
		ids = append(ids, expenditure.Id)
	}

	tags, err := (*r.tagsRepo).ListByType(ctx, "expenditure", &ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	for _, expenditure := range expenditures {
		expenditureTags := make([]openapi.Tag, 0, len(tagsByID[expenditure.Id]))
		for _, tagID := range tagsByID[expenditure.Id] {
			idx := sort.Search(len(tags), func(i int) bool { return tags[i].Id == tagID })
			if idx >= 0 {
				expenditureTags = append(expenditureTags, tags[idx])
			}
		}
		expenditure.Tags = &expenditureTags
	}

	expendituresList := &openapi.ExpenditureList{
		Metadata:     &openapi.ListMetadata{Total: count, Limit: *queryParams.Limit, Offset: *queryParams.Offset},
		Expenditures: &expenditures,
	}

	return expendituresList, nil
}
