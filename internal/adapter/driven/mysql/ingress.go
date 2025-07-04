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
	"strings"
)

type IngressRepo struct {
	db       *sql.DB
	tagsRepo *port.TagsRepo
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

func (i IngressRepo) Create(ctx context.Context, ingress openapi.IngressRequest, transactionID string) (string, error) {
	queryInsert := `insert into ingresses
						(category_id, source, is_recurring, transaction_id, created_at, updated_at) 
					VALUES (?,?,?,?,?,now())`
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
	query := `select i.id,
					   i.source,
					   i.is_recurring,
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
					   GROUP_CONCAT(it.tag_id ORDER BY it.tag_id SEPARATOR ',') as tags
				from ingresses i
						 inner JOIN categories c ON i.category_id = c.id
						 inner join transactions t ON i.transaction_id = t.id
						 left join proletariat_budget.ingress_recurrence_patterns irp on i.id = irp.ingress_id
						 left join proletariat_budget.ingress_tags it ON i.id = it.ingress_id
				WHERE i.id = ?
				GROUP BY i.id,
						 i.source,
						 i.is_recurring,
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
						 irp.amount`
	var ingress openapi.Ingress
	ingressRecurrencePattern := openapi.RecurrencePattern{}
	var tags string
	err := i.db.QueryRowContext(ctx, query, id).Scan(
		&ingress.Id,
		&ingress.Source,
		&ingress.IsRecurring,
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
		&tags,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to select ingress: %w", err)
	}
	if *ingress.IsRecurring {
		ingress.RecurrencePattern = &ingressRecurrencePattern
	}
	ingress.Tags, err = (*i.tagsRepo).GetByIDs(ctx, strings.Split(tags, ","))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}
	return &ingress, nil
}

func (i IngressRepo) List(ctx context.Context, params openapi.ListIngressesParams) (*openapi.IngressList, error) {
	querySelect := `select i.id,
						   i.source,
						   i.is_recurring,
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
						   irp.id                                                   as recurrence_pattern_id,
						   irp.frequency,
						   irp.end_date,
						   irp.amount,
						   GROUP_CONCAT(it.tag_id ORDER BY it.tag_id SEPARATOR ',') as tags
					from ingresses i
							 inner JOIN categories c ON i.category_id = c.id
							 inner join transactions t ON i.transaction_id = t.id
							 left join proletariat_budget.ingress_recurrence_patterns irp on i.id = irp.ingress_id
							 left join proletariat_budget.ingress_tags it ON i.id = it.ingress_id
					GROUP BY i.id,
							 i.source,
							 i.is_recurring,
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
							 irp.amount;`

	queryCount := `SELECT COUNT(*) 
					FROM ingresses i
					inner JOIN categories c ON i.category_id = c.id
					inner join transactions t ON i.transaction_id = t.id
					left join proletariat_budget.ingress_recurrence_patterns irp on i.id = irp.ingress_id
					left join proletariat_budget.ingress_tags it ON i.id = it.ingress_id`

	whereClause := make([]string, 0)
	args := make([]any, 0)

	if params.Category != nil {
		whereClause = append(whereClause, "i.category_id =?")
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
				FROM ingress_tags it2
				WHERE it2.ingress_id = i.id AND it2.tag_id IN ?
			) `)
		args = append(args, tagList)
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
	if params.Currency != nil {
		whereClause = append(whereClause, "t.currency =?")
		args = append(args, params.Currency)
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
	tagsByID := make(map[string][]string)
	var ids []string
	for rows.Next() {
		var ingress openapi.Ingress
		var recurrencePattern openapi.RecurrencePattern
		var tags string
		err := rows.Scan(
			&ingress.Id,
			&ingress.Source,
			&ingress.IsRecurring,
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
			&tags,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if *ingress.IsRecurring {
			ingress.RecurrencePattern = &recurrencePattern
		}
		ingresses = append(ingresses, ingress)

	}
	tags, err := (*i.tagsRepo).ListByType(ctx, "ingress", &ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	for _, ingress := range ingresses {
		ingressTags := make([]openapi.Tag, 0, len(tagsByID[ingress.Id]))
		for _, tagID := range tagsByID[ingress.Id] {
			idx := sort.Search(len(tags), func(i int) bool { return tags[i].Id == tagID })
			if idx >= 0 {
				ingressTags = append(ingressTags, tags[idx])
			}
		}
		ingress.Tags = &ingressTags
	}
	return &openapi.IngressList{Total: &count, Incomes: &ingresses}, nil

}
