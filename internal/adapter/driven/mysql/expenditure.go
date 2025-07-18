package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"sort"
	"strconv"
	"strings"
)

type ExpenditureRepo struct {
	db       *sql.DB
	tagsRepo *port.TagsRepo
}

func NewExpenditureRepo(
	db *sql.DB,
	tagsRepo *port.TagsRepo,
) port.ExpenditureRepo {
	return &ExpenditureRepo{db: db, tagsRepo: tagsRepo}
}

func (r *ExpenditureRepo) Create(
	ctx context.Context,
	expenditure domain.Expenditure,
) (
	string,
	error,
) {
	queryInsert := `insert into expenditures
						(category_id, declared, planned, transaction_id, created_at)
					VALUES (?, ?, ?, ?, ?)`
	result, errInsert := r.db.ExecContext(
		ctx,
		queryInsert,
		expenditure.Category.ID,
		expenditure.Declared,
		expenditure.Planned,
		expenditure.Transaction.ID,
		expenditure.Date,
	)
	if errInsert != nil {
		return "", fmt.Errorf(
			"failed to insert expenditure: %w",
			errInsert,
		)
	}
	expenditureId, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf(
			"failed to get last insert ID: %w",
			err,
		)
	}
	lastIDStr := strconv.FormatInt(
		expenditureId,
		10,
	)
	return lastIDStr, nil
}

func (r *ExpenditureRepo) GetByID(
	ctx context.Context,
	id string,
) (
	*domain.Expenditure,
	error,
) {
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
					where e.id =?
					group by e.id, e.declared, e.planned, e.transaction_id, e.created_at, t.account_id, t.amount, t.currency,
							 t.transaction_date, t.description, c.id, c.name, c.description, c.color, c.background_color, c.active,
							 c.category_type`

	var expenditure domain.Expenditure
	transaction := domain.Transaction{}
	category := domain.Category{}
	var tagsList *string
	err := r.db.QueryRowContext(
		ctx,
		querySelect,
		id,
	).Scan(
		&expenditure.ID,
		&expenditure.Declared,
		&expenditure.Planned,
		&transaction.CreatedAt,
		&transaction.AccountID,
		&transaction.Amount,
		&transaction.Currency,
		&expenditure.Date,
		&transaction.Description,
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Color,
		&category.BackgroundColor,
		&category.Active,
		&tagsList,
	)
	if errors.Is(
		err,
		sql.ErrNoRows,
	) {
		return nil, err // TODO return correct error from domain
	} else if err != nil {
		return nil, fmt.Errorf(
			"failed to select expenditure: %w",
			err,
		)
	}
	expenditure.Transaction = &transaction
	expenditure.Category = &category

	if tagsList != nil && *tagsList != "" {
		tags, err := (*r.tagsRepo).GetByIDs(
			ctx,
			strings.Split(
				*tagsList,
				",",
			),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to fetch tags: %w",
				err,
			)
		}
		expenditure.Tags = tags
	} else {
		tags := make(
			[]*domain.Tag,
			0,
		)
		expenditure.Tags = &tags
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
func (r *ExpenditureRepo) FindExpenditures(
	ctx context.Context,
	queryParams domain.ExpenditureListParams,
) (
	*domain.ExpenditureList,
	error,
) {
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

	if queryParams.CategoryID != nil {
		whereClause = append(
			whereClause,
			"e.category_id =?",
		)
		args = append(
			args,
			*queryParams.CategoryID,
		)
	}
	if queryParams.StartDate != nil && queryParams.EndDate != nil {
		whereClause = append(
			whereClause,
			"t.transaction_date BETWEEN? AND?",
		)
	} else if queryParams.StartDate != nil {
		whereClause = append(
			whereClause,
			"t.transaction_date >=?",
		)
		args = append(
			args,
			*queryParams.StartDate,
		)
	} else if queryParams.EndDate != nil {
		whereClause = append(
			whereClause,
			"t.transaction_date <=?",
		)
		args = append(
			args,
			*queryParams.EndDate,
		)
	}
	if queryParams.Declared != nil {
		whereClause = append(
			whereClause,
			"e.declared =?",
		)
	}
	if queryParams.Planned != nil {
		whereClause = append(
			whereClause,
			"e.planned =?",
		)
	}
	if queryParams.Currency != nil {
		whereClause = append(
			whereClause,
			"t.currency =?",
		)
	}
	if queryParams.Description != nil {
		whereClause = append(
			whereClause,
			"t.description LIKE ?",
		)
		args = append(
			args,
			"%"+*queryParams.Description+"%",
		)
	}
	if queryParams.Tags != nil {
		tagList := make(
			[]string,
			0,
			len(*queryParams.Tags),
		)
		for _, tag := range *queryParams.Tags {
			tagList = append(
				tagList,
				fmt.Sprintf(
					"'%s'",
					tag,
				),
			)
		}
		whereClause = append(
			whereClause,
			fmt.Sprintf(
				` EXISTS (
				SELECT 1
				FROM expenditures_expenditure_tags et2
				WHERE et2.expenditure_id = e.id AND et2.tag_id IN %s
			) `,
				strings.Join(
					tagList,
					", ",
				),
			),
		)
	}
	if queryParams.AccountID != nil {
		// inner select to transaction_expenditures and then to transactions table to get account_id
		whereClause = append(
			whereClause,
			"t.account_id =?",
		)
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
	querySelect += fmt.Sprintf(
		" LIMIT %d OFFSET %d",
		queryParams.Limit,
		queryParams.Offset,
	)
	stmtCount, errQueryCountStmt := r.db.PrepareContext(
		ctx,
		queryCount,
	)
	if errQueryCountStmt != nil {
		return nil, fmt.Errorf(
			"failed to prepare count statement: %w",
			errQueryCountStmt,
		)
	}
	defer stmtCount.Close()
	result, errQueryCount := stmtCount.QueryContext(
		ctx,
		args...,
	)
	if errQueryCount != nil {
		return nil, fmt.Errorf(
			"failed to count expenditures: %w",
			errQueryCount,
		)
	}
	defer result.Close()
	var count int
	result.Next()
	errScan := result.Scan(&count)
	if errScan != nil {
		return nil, errScan
	}

	stmntQuery, errStmtQuerySelect := r.db.PrepareContext(
		ctx,
		querySelect,
	)
	if errStmtQuerySelect != nil {
		return nil, fmt.Errorf(
			"failed to prepare select statement: %w",
			errStmtQuerySelect,
		)
	}

	defer stmntQuery.Close()
	rows, errQuery := stmntQuery.QueryContext(
		ctx,
		args...,
	)
	if errQuery != nil {
		return nil, fmt.Errorf(
			"failed to select expenditures: %w",
			errQuery,
		)
	}

	defer rows.Close()
	var expenditures []domain.Expenditure
	tagsByID := make(map[string][]string)
	var ids []string
	for rows.Next() {
		var expenditure domain.Expenditure
		transaction := domain.Transaction{}
		category := domain.Category{}
		var tags string
		err := rows.Scan(
			&expenditure.ID,
			&expenditure.Declared,
			&expenditure.Planned,
			&transaction.CreatedAt,
			&transaction.AccountID,
			&transaction.Amount,
			&transaction.Currency,
			&expenditure.Date,
			&transaction.Description,
			&category.ID,
			&category.Name,
			&category.Description,
			&category.Color,
			&category.BackgroundColor,
			&category.Active,
			&tags,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to scan row: %w",
				err,
			)
		}
		expenditure.Transaction = &transaction
		expenditure.Category = &category

		expenditures = append(
			expenditures,
			expenditure,
		)
		tagsByID[expenditure.ID] = strings.Split(
			tags,
			",",
		)
		ids = append(
			ids,
			expenditure.ID,
		)
	}

	tags, err := (*r.tagsRepo).ListByType(
		ctx,
		"expenditure",
		&ids,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get tags: %w",
			err,
		)
	}

	for _, expenditure := range expenditures {
		expenditureTags := make(
			[]*domain.Tag,
			0,
			len(tagsByID[expenditure.ID]),
		)
		for _, tagID := range tagsByID[expenditure.ID] {
			idx := sort.Search(
				len(*tags),
				func(i int) bool { return (*tags)[i].ID == tagID },
			)
			if idx >= 0 {
				expenditureTags = append(
					expenditureTags,
					(*tags)[idx],
				)
			}
		}
		expenditure.Tags = &expenditureTags
	}

	expendituresList := &domain.ExpenditureList{
		Metadata: domain.ListMetadata{
			Total:  count,
			Limit:  *queryParams.Limit,
			Offset: *queryParams.Offset,
		},
		Expenditures: expenditures,
	}

	return expendituresList, nil
}
