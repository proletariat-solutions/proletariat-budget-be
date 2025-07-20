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
	baseSelectQuery := `select e.id,
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

	baseCountQuery := `select COUNT(*)
                        from expenditures e
                             inner join categories c ON e.category_id = c.id
                             inner join transactions t ON e.transaction_id = t.id
                             left join expenditure_tags et on e.id = et.expenditure_id`

	whereClause, args := r.buildWhereClause(queryParams)

	count, err := r.getExpendituresCount(ctx, baseCountQuery, whereClause, args)
	if err != nil {
		return nil, err
	}

	expenditures, tagsByID, err := r.getExpenditures(ctx, baseSelectQuery, whereClause, args, queryParams)
	if err != nil {
		return nil, err
	}

	err = r.attachTagsToExpenditures(ctx, expenditures, tagsByID)
	if err != nil {
		return nil, err
	}

	return &domain.ExpenditureList{
		Metadata: domain.ListMetadata{
			Total:  count,
			Limit:  *queryParams.Limit,
			Offset: *queryParams.Offset,
		},
		Expenditures: expenditures,
	}, nil
}

func (r *ExpenditureRepo) buildWhereClause(queryParams domain.ExpenditureListParams) (string, []interface{}) {
	var args []interface{}
	var whereConditions []string

	if queryParams.CategoryID != nil {
		whereConditions = append(whereConditions, "e.category_id = ?")
		args = append(args, *queryParams.CategoryID)
	}

	dateCondition, dateArgs := r.buildDateCondition(queryParams)
	if dateCondition != "" {
		whereConditions = append(whereConditions, dateCondition)
		args = append(args, dateArgs...)
	}

	if queryParams.Declared != nil {
		whereConditions = append(whereConditions, "e.declared = ?")
		args = append(args, *queryParams.Declared)
	}

	if queryParams.Planned != nil {
		whereConditions = append(whereConditions, "e.planned = ?")
		args = append(args, *queryParams.Planned)
	}

	if queryParams.Currency != nil {
		whereConditions = append(whereConditions, "t.currency = ?")
		args = append(args, *queryParams.Currency)
	}

	if queryParams.Description != nil {
		whereConditions = append(whereConditions, "t.description LIKE ?")
		args = append(args, "%"+*queryParams.Description+"%")
	}

	if queryParams.Tags != nil {
		tagCondition := r.buildTagsCondition(*queryParams.Tags)
		whereConditions = append(whereConditions, tagCondition)
	}

	if queryParams.AccountID != nil {
		whereConditions = append(whereConditions, "t.account_id = ?")
		args = append(args, *queryParams.AccountID)
	}

	if len(whereConditions) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(whereConditions, " AND "), args
}

func (r *ExpenditureRepo) buildDateCondition(queryParams domain.ExpenditureListParams) (string, []interface{}) {
	var args []interface{}

	if queryParams.StartDate != nil && queryParams.EndDate != nil {
		args = append(args, *queryParams.StartDate, *queryParams.EndDate)
		return "t.transaction_date BETWEEN ? AND ?", args
	}

	if queryParams.StartDate != nil {
		args = append(args, *queryParams.StartDate)
		return "t.transaction_date >= ?", args
	}

	if queryParams.EndDate != nil {
		args = append(args, *queryParams.EndDate)
		return "t.transaction_date <= ?", args
	}

	return "", args
}

func (r *ExpenditureRepo) buildTagsCondition(tags []string) string {
	tagList := make([]string, 0, len(tags))
	for _, tag := range tags {
		tagList = append(tagList, fmt.Sprintf("'%s'", tag))
	}

	return fmt.Sprintf(`EXISTS (
        SELECT 1
        FROM expenditures_expenditure_tags et2
        WHERE et2.expenditure_id = e.id AND et2.tag_id IN (%s)
    )`, strings.Join(tagList, ", "))
}

func (r *ExpenditureRepo) getExpendituresCount(ctx context.Context, baseQuery, whereClause string, args []interface{}) (int, error) {
	query := baseQuery + whereClause

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare count statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count expenditures: %w", err)
	}
	defer result.Close()

	var count int
	if result.Next() {
		if err := result.Scan(&count); err != nil {
			return 0, fmt.Errorf("failed to scan count: %w", err)
		}
	}

	return count, nil
}

func (r *ExpenditureRepo) getExpenditures(ctx context.Context, baseQuery, whereClause string, args []interface{}, queryParams domain.ExpenditureListParams) ([]domain.Expenditure, map[string][]string, error) {
	query := baseQuery + whereClause +
		` GROUP BY e.id, e.declared, e.planned, e.transaction_id, e.created_at, t.account_id, t.amount, t.currency,
          t.transaction_date, t.description, c.id, c.name, c.description, c.color, c.background_color, c.active, c.category_type` +
		" ORDER BY created_at DESC" +
		fmt.Sprintf(" LIMIT %d OFFSET %d", queryParams.Limit, queryParams.Offset)

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to prepare select statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to select expenditures: %w", err)
	}
	defer rows.Close()

	var expenditures []domain.Expenditure
	tagsByID := make(map[string][]string)

	for rows.Next() {
		expenditure, tags, err := r.scanExpenditureRow(rows)
		if err != nil {
			return nil, nil, err
		}

		expenditures = append(expenditures, expenditure)
		tagsByID[expenditure.ID] = strings.Split(tags, ",")
	}

	return expenditures, tagsByID, nil
}

func (r *ExpenditureRepo) scanExpenditureRow(rows *sql.Rows) (domain.Expenditure, string, error) {
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
		return expenditure, "", fmt.Errorf("failed to scan row: %w", err)
	}

	expenditure.Transaction = &transaction
	expenditure.Category = &category

	return expenditure, tags, nil
}

func (r *ExpenditureRepo) attachTagsToExpenditures(ctx context.Context, expenditures []domain.Expenditure, tagsByID map[string][]string) error {
	ids := make([]string, 0, len(expenditures))
	for _, expenditure := range expenditures {
		ids = append(ids, expenditure.ID)
	}

	tags, err := (*r.tagsRepo).ListByType(ctx, "expenditure", &ids)
	if err != nil {
		return fmt.Errorf("failed to get tags: %w", err)
	}

	for i := range expenditures {
		expenditureTags := make([]*domain.Tag, 0, len(tagsByID[expenditures[i].ID]))
		for _, tagID := range tagsByID[expenditures[i].ID] {
			idx := sort.Search(len(*tags), func(j int) bool {
				return (*tags)[j].ID == tagID
			})
			if idx < len(*tags) && (*tags)[idx].ID == tagID {
				expenditureTags = append(expenditureTags, (*tags)[idx])
			}
		}
		expenditures[i].Tags = &expenditureTags
	}

	return nil
}
