package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"ghorkov32/proletariat-budget-be/internal/core/port"
	"ghorkov32/proletariat-budget-be/openapi"
)

type IngressRepo struct {
	db       *sql.DB
	tagsRepo port.TagsRepo
}

func NewIngressRepo(
	db *sql.DB,
	tagsRepo port.TagsRepo,
) port.IngressRepo {
	return &IngressRepo{
		db:       db,
		tagsRepo: tagsRepo,
	}
}

func (i IngressRepo) CreateRecurrencePattern(
	ctx context.Context,
	recurrencePattern openapi.RecurrencePatternRequest,
) (
	string,
	error,
) {
	queryInsert := `insert into ingress_recurrence_patterns 
						(frequency,
						 interval_value,
						 amount,
						 to_account_id,
						 description,
						 end_date)
					VALUES (?,
							?,
							?,
							?,
							?,
							?)`
	result, errInsert := i.db.ExecContext(
		ctx,
		queryInsert,
		recurrencePattern.Frequency,
		recurrencePattern.Interval,
		recurrencePattern.Amount,
		recurrencePattern.ToAccountId,
		recurrencePattern.Description,
		recurrencePattern.EndDate,
	)
	if errInsert != nil {
		return "", fmt.Errorf(
			"failed to create recurrence pattern: %w",
			errInsert,
		)
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

func (i IngressRepo) UpdateRecurrencePattern(
	ctx context.Context,
	id string,
	recurrencePattern openapi.RecurrencePattern,
) error {
	queryUpdate := `UPDATE ingress_recurrence_patterns SET frequency=?, interval_value=?, amount=?, to_account_id=?, description=?, end_date=? WHERE id=?`
	_, err := i.db.ExecContext(
		ctx,
		queryUpdate,
		recurrencePattern.Frequency,
		recurrencePattern.Interval,
		recurrencePattern.Amount,
		recurrencePattern.ToAccountId,
		recurrencePattern.Description,
		recurrencePattern.EndDate,
		id,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to update recurrence pattern: %w",
			err,
		)
	}

	return nil
}

func (i IngressRepo) DeleteRecurrencePattern(
	ctx context.Context,
	id string,
) error {
	queryDelete := `DELETE FROM ingress_recurrence_patterns WHERE id=?`
	_, err := i.db.ExecContext(
		ctx,
		queryDelete,
		id,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to delete recurrence pattern: %w",
			err,
		)
	}

	return nil
}

func (i IngressRepo) GetRecurrencePattern(
	ctx context.Context,
	id string,
) (
	*openapi.RecurrencePattern,
	error,
) {
	query := `SELECT id, 
					 frequency,
					 interval_value,
					 amount,
					 to_account_id,
					 description,
						 end_date FROM ingress_recurrence_patterns WHERE id=?`

	var recurrencePattern openapi.RecurrencePattern
	err := i.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&recurrencePattern.Id,
		&recurrencePattern.Frequency,
		&recurrencePattern.Interval,
		&recurrencePattern.Amount,
		&recurrencePattern.ToAccountId,
		&recurrencePattern.Description,
		&recurrencePattern.EndDate,
	)
	if err != nil {
		return nil, translateError(err)
	}

	return &recurrencePattern, nil
}

func (i IngressRepo) Create(
	ctx context.Context,
	ingress openapi.IngressRequest,
	transactionID string,
) (
	string,
	error,
) {
	queryInsert := `insert into ingresses
						(category_id, source, from_recurrency_pattern_id, transaction_id, created_at) 
					VALUES (?,?,?,?,now())`
	result, errInsert := i.db.ExecContext(
		ctx,
		queryInsert,
		ingress.Category,
		ingress.Source,
		ingress.RecurrencePattern.Id,
		transactionID,
		ingress.CreatedAt,
		ingress.Amount,
		ingress.Currency,
	)
	if errInsert != nil {
		return "", fmt.Errorf(
			"failed to create ingress: %w",
			errInsert,
		)
	}
	ingressID, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf(
			"failed to get last insert ID: %w",
			err,
		)
	}
	lastIDStr := strconv.FormatInt(
		ingressID,
		10,
	)

	return lastIDStr, nil
}

func (i IngressRepo) GetByID(
	ctx context.Context,
	id string,
) (
	*openapi.Ingress,
	error,
) {
	query := `select i.id,
					   i.source,
					   i.created_at,
					   t.amount,
					   t.currency,
					   t.account_id,
					   c.id,
					   c.name,
					   c.description,
					   c.color,
					   c.background_color,
					   c.active,
					   irp.id as recurrence_pattern_id,
					   irp.frequency,
					   irp.end_date,
					   irp.amount,
					   irp.interval_value,
					   irp.to_account_id,
					   irp.description,
					   GROUP_CONCAT(it.tag_id ORDER BY it.tag_id SEPARATOR ',') as tags
				from ingresses i
						 inner JOIN categories c ON i.category_id = c.id
						 inner join transactions t ON i.transaction_id = t.id
						 left join proletariat_budget.ingress_recurrence_patterns irp on i.from_recurrency_pattern_id = irp.id
						 left join proletariat_budget.ingress_tags it ON i.id = it.ingress_id
				WHERE i.id = ?
				GROUP BY i.id,
						 i.source,
						 i.created_at,
						 t.amount,
						 t.currency,
						 t.account_id,
						 c.id,
						 c.name,
						 c.description,
						 c.color,
						 c.background_color,
						 c.active,
						 irp.id,
						 irp.frequency,
						 irp.end_date,
						 irp.amount,
						 irp.interval_value,
						 irp.to_account_id,
						 irp.description`
	var ingress openapi.Ingress
	ingressRecurrencePattern := openapi.RecurrencePattern{}
	var tags string
	err := i.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&ingress.Id,
		&ingress.Source,
		&ingress.CreatedAt,
		&ingress.Amount,
		&ingress.Currency,
		&ingress.AccountId,
		&ingress.Category.Id,
		&ingress.Category.Name,
		&ingress.Category.Description,
		&ingress.Category.Color,
		&ingress.Category.BackgroundColor,
		&ingress.Category.Active,
		&ingressRecurrencePattern.Id,
		&ingressRecurrencePattern.Interval,
		&ingressRecurrencePattern.Frequency,
		&ingressRecurrencePattern.EndDate,
		&ingressRecurrencePattern.Amount,
		&ingressRecurrencePattern.Interval,
		&ingressRecurrencePattern.ToAccountId,
		&ingressRecurrencePattern.Description,
		&tags,
	)
	if err != nil {
		return nil, translateError(err)
	}
	if ingressRecurrencePattern.Id != "" {
		ingress.RecurrencePattern = &ingressRecurrencePattern
	}

	return &ingress, nil
}

func (i IngressRepo) List(
	ctx context.Context,
	params openapi.ListIngressesParams,
) (
	*openapi.IngressList,
	error,
) {
	querySelect := `select i.id,
						   i.source,
						   i.created_at,
						   t.amount,
						   t.currency,
						   t.account_id,
						   c.id,
						   c.name,
						   c.description,
						   c.color,
						   c.background_color,
						   c.active,
						   irp.id as recurrence_pattern_id,
						   irp.frequency,
						   irp.end_date,
						   irp.amount,
						   irp.interval_value,
						   irp.to_account_id,
						   irp.description,
						   GROUP_CONCAT(it.tag_id ORDER BY it.tag_id SEPARATOR ',') as tags
					from ingresses i
							 inner JOIN categories c ON i.category_id = c.id
							 inner join transactions t ON i.transaction_id = t.id
							 left join proletariat_budget.ingress_recurrence_patterns irp on i.from_recurrency_pattern_id = irp.id
							 left join proletariat_budget.ingress_tags it ON i.id = it.ingress_id
					GROUP BY i.id,
							 i.source,
							 i.created_at,
							 t.amount,
							 t.currency,
							 t.account_id,
							 c.id,
							 c.name,
							 c.description,
							 c.color,
							 c.background_color,
							 c.active,
							 irp.id,
							 irp.frequency,
							 irp.end_date,
							 irp.amount,
							 irp.interval_value,
							 irp.to_account_id,
							 irp.description`

	queryCount := `SELECT COUNT(*) 
					FROM ingresses i
					inner JOIN categories c ON i.category_id = c.id
					inner join transactions t ON i.transaction_id = t.id
					left join proletariat_budget.ingress_recurrence_patterns irp on i.from_recurrency_pattern_id = irp.id
					left join proletariat_budget.ingress_tags it ON i.id = it.ingress_id`

	whereClause := make(
		[]string,
		0,
	)
	args := make(
		[]any,
		0,
	)

	if params.Category != nil {
		whereClause = append(
			whereClause,
			"i.category_id =?",
		)
		args = append(
			args,
			*params.Category,
		)
	}
	if params.Source != nil {
		whereClause = append(
			whereClause,
			"i.source like %?%",
		)
		args = append(
			args,
			*params.Source,
		)
	}
	if params.Tags != nil && len(*params.Tags) > 0 {
		tagList := make(
			[]string,
			0,
			len(*params.Tags),
		)
		for _, tag := range *params.Tags {
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
			` EXISTS (
				SELECT 1
				FROM ingress_tags it2
				WHERE it2.ingress_id = i.id AND it2.tag_id IN ?
			) `,
		)
		args = append(
			args,
			tagList,
		)
	}
	if params.StartDate != nil && params.EndDate != nil {
		whereClause = append(
			whereClause,
			"t.transaction_date BETWEEN? AND?",
		)
		args = append(
			args,
			*params.StartDate,
			*params.EndDate,
		)
	} else if params.StartDate != nil {
		whereClause = append(
			whereClause,
			"t.transaction_date >=?",
		)
		args = append(
			args,
			*params.StartDate,
		)
	} else if params.EndDate != nil {
		whereClause = append(
			whereClause,
			"t.transaction_date <=?",
		)
		args = append(
			args,
			*params.EndDate,
		)
	}
	if params.Currency != nil {
		whereClause = append(
			whereClause,
			"t.currency =?",
		)
		args = append(
			args,
			params.Currency,
		)
	}
	if params.IsRecurring != nil {
		var condition string
		if *params.IsRecurring {
			condition = "is not null"
		} else {
			condition = "is null"
		}
		whereClause = append(
			whereClause,
			"i.from_recurrency_pattern_id "+condition,
		)
		args = append(
			args,
			*params.IsRecurring,
		)
	}
	if len(whereClause) > 0 {
		querySelect += " WHERE "
		queryCount += " WHERE "
	}
	for i, clause := range whereClause {
		if i > 0 {
			querySelect += AND_CLAUSE
			queryCount += AND_CLAUSE
		}
		querySelect += clause
		queryCount += clause
	}
	querySelect += " ORDER BY i.created_at DESC"
	querySelect += " LIMIT ?"
	querySelect += " OFFSET ?"
	countArgs := args
	args = append(
		args,
		params.Limit,
		params.Offset,
	)
	var ingresses []openapi.Ingress
	var count int

	err := i.db.QueryRowContext(
		ctx,
		queryCount,
		countArgs...,
	).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to count rows: %w",
			err,
		)
	}

	rows, err := i.db.QueryContext(
		ctx,
		querySelect,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to select ingresses: %w",
			err,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var ingress openapi.Ingress
		var recurrencePattern openapi.RecurrencePattern
		var tags string
		errScan := rows.Scan(
			&ingress.Id,
			&ingress.Source,
			&ingress.CreatedAt,
			&ingress.Amount,
			&ingress.Currency,
			&ingress.AccountId,
			&ingress.Category.Id,
			&ingress.Category.Name,
			&ingress.Category.Description,
			&ingress.Category.Color,
			&ingress.Category.BackgroundColor,
			&ingress.Category.Active,
			&recurrencePattern.Id,
			&recurrencePattern.Interval,
			&recurrencePattern.Frequency,
			&recurrencePattern.EndDate,
			&recurrencePattern.Amount,
			&recurrencePattern.Interval,
			&recurrencePattern.ToAccountId,
			&recurrencePattern.Description,
			&tags,
		)
		if errScan != nil {
			return nil, fmt.Errorf(
				"failed to scan row: %w",
				errScan,
			)
		}
		if recurrencePattern.Id != "" {
			ingress.RecurrencePattern = &recurrencePattern
		}
		ingresses = append(
			ingresses,
			ingress,
		)
	}

	return &openapi.IngressList{
		Metadata: &openapi.ListMetadata{
			Total:  count,
			Offset: *params.Offset,
			Limit:  *params.Limit,
		},
		Incomes: &ingresses,
	}, nil
}
