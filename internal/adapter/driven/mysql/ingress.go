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
	queryInsert := `INSERT INTO ingresses (category, source, date, is_recurring, created_at, updated_at) VALUES (?,?,?,?,NOW(),NOW())`
	result, errInsert := i.db.ExecContext(
        ctx,
        queryInsert,
        ingress.Category,
        ingress.Source,
		ingress.Description,
        ingress.Date,
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
	queryUpdate := `UPDATE ingresses SET category=?, source=?, description=?, date=?, is_recurring=?, updated_at=NOW() WHERE id=?`
	_, err := i.db.ExecContext(
        ctx,
        queryUpdate,
        ingress.Category,
        ingress.Source,
		ingress.Description,
        ingress.Date,
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
	query := `SELECT id, category, source, description, date, is_recurring, created_at, updated_at FROM ingresses WHERE id=?`
	var ingress openapi.Ingress
	err := i.db.QueryRowContext(ctx, query, id).Scan(
        &ingress.Id,
        &ingress.Category,
        &ingress.Source,
        &ingress.Description,
		&ingress.Date,
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
	query := `SELECT 
		i.id, 
		category, 
		source, 
		t.description,
		date, 
		is_recurring, 
		group_concat(itj.tag_id SEPARATOR ', ') as tags,
		irp.frequency as frequency,
        irp.interval_value as interval_value,
        irp.amount as amount,
        irp.end_date as end_interval_date
	FROM ingresses i 
    left join ingress_recurrence_patterns irp on i.id = irp.ingress_id
	left join ingress_tags_junction itj ON i.id = itj.ingress_id
	inner join transaction_ingresses ti on i.id = ti.ingress_id
	inner join proletariat_budget.transactions t on ti.transaction_id = t.id`

	countQuery := `SELECT COUNT(*) FROM ingresses i`

	whereClause := make([]string, 0)

	if params.Category != nil {
        whereClause = append(whereClause, "i.category =?")
        countQuery += " WHERE i.category =?"
    }
	if params.Source != nil {
        whereClause = append(whereClause, "i.source =?")
        countQuery += " AND i.source =?"
    }
	if params.Tags!= nil && len(*params.Tags) > 0 {
		whereClause = append(whereClause, "itj.tag_id IN ("+strings.Repeat("?", len(*params.Tags))+")")
        countQuery += " AND itj.tag_id IN ("+strings.Repeat("?", len(*params.Tags))+")"
        for _, tag := range *params.Tags {
            whereClause[len(whereClause)-1] += fmt.Sprintf(" AND itj.tag_id = %d", tag)
            countQuery += fmt.Sprintf(" AND itj.tag_id = %d", tag)
        }
	}
	if params.Source != nil {
		whereClause = append(whereClause, "i.source =?")
        countQuery += " AND i.source =?"
	}
	if params.EndDate != nil {
		whereClause = append(whereClause, "i.date <=?")
        countQuery += " AND i.date <=?"
	}
	if params.StartDate!= nil {
		whereClause = append(whereClause, "i.date >=?")
        countQuery += " AND i.date >=?"
	}
	if params.IsRecurring != nil {
        whereClause = append(whereClause, "i.is_recurring =?")
        countQuery += " AND i.is_recurring =?"
    }
	if len(whereClause) > 0 {
        query += " WHERE " + strings.Join(whereClause, " AND ")
        countQuery += " WHERE " + strings.Join(whereClause, " AND ")
    }
	query += " ORDER BY i.created_at DESC"
	query += " LIMIT?"
	queryCount := countQuery + " LIMIT?"
	query += " OFFSET?"
	var ingresses []openapi.Ingress
	var count int
	err := i.db.QueryRowContext(ctx, queryCount, params.Category, params.Category, params.Source, params.Source, params.Limit, params.Limit, params.Offset).Scan(&count)
	if err != nil {
        return nil, fmt.Errorf("failed to count rows: %w", err)
    }
	query += " LIMIT?"
	rows, err := i.db.QueryContext(ctx, query, params.Category, params.Category, params.Source, params.Source, params.Limit, params.Limit, params.Offset)
	if err != nil {
        return nil, fmt.Errorf("failed to select ingresses: %w", err)
    }
	defer rows.Close()
	for rows.Next() {
		var ingress openapi.Ingress
		var recurrencePattern openapi.RecurrencePattern
		var tags string
        err := rows.Scan(
            &ingress.Id,
            &ingress.Category,
            &ingress.Source,
            &ingress.Description,
            &ingress.Date,
            &ingress.IsRecurring,
			tags,
			&recurrencePattern.Frequency,
			&recurrencePattern.Interval,
            &ingress.Amount,
            &recurrencePattern.EndDate,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan row: %w", err)
        }
		if *ingress.IsRecurring {
            ingress.RecurrencePattern = &recurrencePattern
        }
		tagIDs := strings.Split(tags, ", ")
		if len(tagIDs) > 0 {
			ingress.Tags = &tagIDs
		}
        ingresses = append(ingresses, ingress)

	}
	return ingresses, nil

}

func (i IngressRepo) ListCategories(ctx context.Context) ([]openapi.IngressCategory, error) {
	query := `SELECT id, name, description, color, background_color, active from ingress_categories`
	rows, err := i.db.QueryContext(ctx, query)
	if err != nil {
        return nil, fmt.Errorf("failed to select categories: %w", err)
    }
	defer rows.Close()
	var categories []openapi.IngressCategory
	for rows.Next() {
		var category openapi.IngressCategory
        err := rows.Scan(&category.Id, &category.Name, &category.Description, &category.Color, &category.BackgroundColor, &category.Active)
        if err != nil {
            return nil, fmt.Errorf("failed to scan row: %w", err)
        }
        categories = append(categories, category)
	}
	return categories, nil
}

func (i IngressRepo) GetCategory(ctx context.Context, id string) (*openapi.IngressCategory, error) {
	query := `SELECT id, name, description, color, background_color, active from ingress_categories WHERE id =?`
    row := i.db.QueryRowContext(ctx, query, id)
    var category openapi.IngressCategory
    err := row.Scan(&category.Id, &category.Name, &category.Description, &category.Color, &category.BackgroundColor, &category.Active)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, domain.ErrEntityNotFound
    } else if err != nil {
        return nil, fmt.Errorf("failed to select category: %w", err)
    }
    return &category, nil
}

func (i IngressRepo) CreateCategory(ctx context.Context, category openapi.IngressCategoryRequest) (string, error) {
	query := `INSERT INTO ingress_categories (name, description, color, background_color, active) VALUES (?,?,?,?,true) RETURNING id`
	result, err := i.db.ExecContext(
        ctx, query, category.Name, category.Description, category.Color, category.BackgroundColor,
    )
	if err != nil {
        return "", fmt.Errorf("failed to create category: %w", err)
    }
	id, err := result.LastInsertId()
	if err != nil {
        return "", fmt.Errorf("failed to get last insert id: %w", err)
    }
	return strconv.FormatInt(id, 10), nil
}

func (i IngressRepo) UpdateCategory(ctx context.Context, id string, category openapi.IngressCategoryRequest) error {
	query := `UPDATE ingress_categories SET name =?, description =?, color =?, background_color =? WHERE id =?`
	_, err := i.db.ExecContext(ctx, query, category.Name, category.Description, category.Color, category.BackgroundColor, id)
	if err != nil {
        return fmt.Errorf("failed to update category: %w", err)
    }
	return nil
}

func (i IngressRepo) DeleteCategory(ctx context.Context, id string) error {
	query := `UPDATE ingress_categories SET active = false WHERE id =?`
	_, err := i.db.ExecContext(ctx, query, id)
	if err != nil {
        return fmt.Errorf("failed to delete category: %w", err)
    }
	return nil
}

func (i IngressRepo) FindOrCreateTags(ctx context.Context, tags []string) ([]string, error) {
	query := `INSERT INTO ingress_tags (name, description, color, background_color) 
				VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE name =?, description =?, color =?, background_color =?`
	var insertedTagIDs []string
	for _, tag := range tags {
		colors, errColors := common.GetPillColor(tag)
		if errColors != nil {
            return nil, fmt.Errorf("failed to get tag colors: %w", errColors)
        }
		result, err := i.db.ExecContext(
            ctx, query, tag, tag, (*colors)[0], (*colors)[1],
        )
		if err != nil {
            return nil, fmt.Errorf("failed to insert or update tag: %w", err)
        }
		lastInsertId, err := result.LastInsertId()
		if err != nil {
            return nil, fmt.Errorf("failed to get last insert id: %w", err)
        }
		insertedTagIDs = append(insertedTagIDs, fmt.Sprintf("%d", lastInsertId))
	}
	return insertedTagIDs, nil

}

func (i IngressRepo) LinkTagsToIngress(ctx context.Context, tags []string, ingressId string) error {
	query := `INSERT INTO ingress_tags_junction (ingress_id, tag_id) VALUES (?,?)`
	for _, tagId := range tags {
		_, err := i.db.ExecContext(ctx, query, ingressId, tagId)
        if err != nil {
            return fmt.Errorf("failed to link tag to ingress: %w", err)
        }
	}
	return nil
}
