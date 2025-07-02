package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
)

type IngressRepo struct {
	db *sql.DB
}

func (i IngressRepo) CreateRecurrencePattern(ctx context.Context, recurrencePattern domain.RecurrencePattern) (string, error) {
	queryInsert := `INSERT INTO ingress_recurrence_patterns (ingress_id, frequency, interval_value, amount, end_date) VALUES (?,?,?,?,?)`
	result, errInsert := i.db.ExecContext(
		ctx,
		queryInsert,
		recurrencePattern.IngressID,
		recurrencePattern.Frequency,
		recurrencePattern.Interval,
		recurrencePattern.Amount,
		recurrencePattern.EndDate,
	)
	if errInsert != nil {
		return "", fmt.Errorf("failed to create recurrence pattern: %w", errInsert)
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("failed to get last insert ID: %w", err)
	}
	lastIDStr := fmt.Sprintf("%d", lastID)
	return lastIDStr, nil
}

func (i IngressRepo) UpdateRecurrencePattern(ctx context.Context, id string, recurrencePattern domain.RecurrencePattern) error {
	queryUpdate := `UPDATE ingress_recurrence_patterns SET frequency=?, interval_value=?, amount=?, end_date=? WHERE id=?`
	_, err := i.db.ExecContext(
		ctx,
		queryUpdate,
		recurrencePattern.Frequency,
		recurrencePattern.Interval,
		recurrencePattern.Amount,
		recurrencePattern.EndDate,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update recurrence pattern: %w", err)
	}
	return nil
}

func (i IngressRepo) DeleteRecurrencePattern(ctx context.Context, id string) error {

	queryDelete := `DELETE FROM ingress_recurrence_patterns WHERE id=?`
	_, err := i.db.ExecContext(
		ctx,
		queryDelete,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete recurrence pattern: %w", err)
	}
	return nil
}

func (i IngressRepo) GetRecurrencePattern(ctx context.Context, id string) (*domain.RecurrencePattern, error) {
	query := `SELECT id, frequency, interval_value, amount, end_date FROM ingress_recurrence_patterns WHERE id=?`

	var recurrencePattern domain.RecurrencePattern
	err := i.db.QueryRowContext(ctx, query, id).Scan(
		&recurrencePattern.ID,
		&recurrencePattern.Frequency,
		&recurrencePattern.Interval,
		&recurrencePattern.Amount,
		&recurrencePattern.EndDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to select recurrence pattern: %w", err)
	}
	return &recurrencePattern, nil
}

func (i IngressRepo) Create(ctx context.Context, ingress openapi.IngressRequest) (string, error) {
	queryInsert := `INSERT INTO ingresses (category_id, source, is_recurring, created_at, updated_at) VALUES (?,?,?,?,NOW(),NOW())`
	result, errInsert := i.db.ExecContext(
		ctx,
		queryInsert,
		ingress.Category,
		ingress.Source,
		ingress.Description,
		ingress.IsRecurring,
	)
	if errInsert != nil {
		return "", fmt.Errorf("failed to create ingress: %w", errInsert)
	}
	ingressID, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("failed to get last insert ID: %w", err)
	}
	lastIDStr := fmt.Sprintf("%d", ingressID)
	return lastIDStr, nil
}

func (i IngressRepo) Update(ctx context.Context, id string, ingress openapi.IngressRequest) error {
	queryUpdate := `UPDATE ingresses SET category_id=?, source=?, is_recurring=?, updated_at=NOW() WHERE id=?`
	_, err := i.db.ExecContext(
		ctx,
		queryUpdate,
		ingress.Category,
		ingress.Source,
		ingress.IsRecurring,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update ingress: %w", err)
	}
	return nil
}

func (i IngressRepo) Delete(ctx context.Context, id string) error {
	queryDelete := `DELETE FROM ingresses WHERE id=?`
	_, err := i.db.ExecContext(
		ctx,
		queryDelete,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete ingress: %w", err)
	}
	return nil
}

func (i IngressRepo) GetByID(ctx context.Context, id string) (*openapi.Ingress, error) {
	query := `SELECT id, category_id, source, is_recurring, created_at, updated_at FROM ingresses WHERE id=?`
	var ingress openapi.Ingress
	err := i.db.QueryRowContext(ctx, query, id).Scan(
		&ingress.Id,
		&ingress.Category,
		&ingress.Source,
		&ingress.IsRecurring,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to select ingress: %w", err)
	}
	return &ingress, nil
}

func (i IngressRepo) List(ctx context.Context, params openapi.ListIngressesParams) ([]openapi.Ingress, error) {
	querySelect := `select t.account_id,
						   t.amount,
						   i.category_id,
						   t.created_at                                               as createdAt,
						   t.currency,
						   t.transaction_date                                         as date,
						   t.description,
						   i.id,
						   i.is_recurring,
						   i.source,
						   GROUP_CONCAT(itj.tag_id ORDER BY itj.tag_id SEPARATOR ',') as tags,
						   irp.end_date,
						   irp.frequency,
						   irp.interval_value
					from ingresses i
						 inner join
							 transactions t
							 on i.transaction_id = t.id
						 left join
							 ingress_tags itj
							 on i.id = itj.ingress_id
						 left join
							 ingress_recurrence_patterns irp
							 on i.id = irp.ingress_id`

	queryCount := `SELECT COUNT(*) FROM ingresses i
						        inner join
									 transactions t
									 on i.transaction_id = t.id
								 left join
									 ingress_tags itj
									 on i.id = itj.ingress_id
								 left join
									 ingress_recurrence_patterns irp
									 on i.id = irp.ingress_id`

	whereClause := make([]string, 0)
	args := make([]any, 0)

	if params.Category != nil {
		whereClause = append(whereClause, "i.category =?")
		args = append(args, *params.Category)
	}
	if params.Source != nil {
		whereClause = append(whereClause, "i.source like %?%")
		args = append(args, *params.Source)
	}
	if params.Tags != nil && len(*params.Tags) > 0 {
		tagList := make([]string, 0, len(*params.Tags))
		for _, tag := range *params.Tags {
			tagList = append(tagList, fmt.Sprintf("'%s'", tag))
		}
		whereClause = append(whereClause, ` EXISTS (
				SELECT 1
				FROM ingress_tags itj2
				WHERE itj2.ingress_id = i.id AND eet2.tag_id IN ?
			) `)
		args = append(args, tagList)
	}
	if params.Source != nil {
		whereClause = append(whereClause, "i.source =?")
		args = append(args, *params.Source)
	}
	if params.StartDate != nil && params.EndDate != nil {
		whereClause = append(whereClause, "t.transaction_date BETWEEN? AND?")
		args = append(args, *params.StartDate, *params.EndDate)
	} else if params.StartDate != nil {
		whereClause = append(whereClause, "t.transaction_date >=?")
		args = append(args, *params.StartDate)
	} else if params.EndDate != nil {
		whereClause = append(whereClause, "t.transaction_date <=?")
		args = append(args, *params.EndDate)
	}
	if params.IsRecurring != nil {
		whereClause = append(whereClause, "i.is_recurring =?")
		args = append(args, *params.IsRecurring)
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
	querySelect += " ORDER BY i.created_at DESC"
	querySelect += " LIMIT ?"
	querySelect += " OFFSET ?"
	countArgs := args
	args = append(args, params.Limit, params.Offset)
	var ingresses []openapi.Ingress
	var count int

	err := i.db.QueryRowContext(ctx, queryCount, countArgs...).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to count rows: %w", err)
	}

	rows, err := i.db.QueryContext(ctx, querySelect, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select ingresses: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var ingress openapi.Ingress
		var recurrencePattern openapi.RecurrencePattern
		err := rows.Scan(
			&ingress.AccountId,
			&ingress.Amount,
			&ingress.Category,
			&ingress.CreatedAt,
			&ingress.Currency,
			&ingress.Date,
			&ingress.Description,
			&ingress.Id,
			&ingress.IsRecurring,
			&ingress.Source,
			&recurrencePattern.EndDate,
			&recurrencePattern.Frequency,
			&recurrencePattern.Interval,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if *ingress.IsRecurring {
			ingress.RecurrencePattern = &recurrencePattern
		}
		ingresses = append(ingresses, ingress)

	}
	return ingresses, nil

}
