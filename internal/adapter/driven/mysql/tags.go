package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"ghorkov32/proletariat-budget-be/openapi"
	"strings"
)

type TagsRepoImpl struct {
	db *sql.DB
}

func (t TagsRepoImpl) Create(ctx context.Context, tag openapi.Tag, tagType string) (string, error) {
	queryInsert := `INSERT INTO tags (name, description, color, background_color, type, created_at) VALUES (?,?,?,?,?, now())`
	result, errInsert := t.db.ExecContext(
		ctx,
		queryInsert,
		tag.Name,
		tag.Description,
		tag.Color,
		tag.BackgroundColor,
		tagType,
	)
	if errInsert != nil {
		return "", fmt.Errorf("failed to create tag: %w", errInsert)
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("failed to get last insert ID: %w", err)
	}
	lastIDStr := fmt.Sprintf("%d", lastID)
	return lastIDStr, nil
}

func (t TagsRepoImpl) Update(ctx context.Context, id string, tag openapi.Tag, tagType string) error {
	queryUpdate := `UPDATE tags SET name=?, description=?, color=?, background_color=?, type=? WHERE id=?`
	_, err := t.db.ExecContext(
		ctx,
		queryUpdate,
		tag.Name,
		tag.Description,
		tag.Color,
		tag.BackgroundColor,
		tagType,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}
	return nil
}

func (t TagsRepoImpl) Delete(ctx context.Context, id string) error {
	// Deleting linked records first
	tables := []string{"expenditure_tags", "ingress_tags", "savings_goal_tags", "saving_withdrawal_tags", "savings_contribution_tags"}
	queryDeleteLinked := `DELETE FROM ? WHERE tag_id=?`
	for _, table := range tables {
		_, err := t.db.ExecContext(
			ctx,
			queryDeleteLinked,
			table,
			id,
		)
		if err != nil {
			return fmt.Errorf("failed to delete linked records: %w", err)
		}
	}

	queryDelete := `DELETE FROM tags WHERE id=?`
	_, err := t.db.ExecContext(
		ctx,
		queryDelete,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil

}

func (t TagsRepoImpl) GetByID(ctx context.Context, id string) (*openapi.Tag, error) {
	query := `SELECT id, name, description, color, background_color FROM tags WHERE id=?`
	var tag openapi.Tag
	err := t.db.QueryRowContext(ctx, query, id).Scan(
		&tag.Id,
		&tag.Name,
		&tag.Description,
		&tag.Color,
		&tag.BackgroundColor,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to select tag: %w", err)
	}
	return &tag, nil
}

func (t TagsRepoImpl) GetByIDs(ctx context.Context, ids []string) (*[]openapi.Tag, error) {
	query := `SELECT id, name, description, color, background_color FROM tags WHERE id IN (?)`
	strIDs := strings.Join(ids, ",")
	rows, err := t.db.QueryContext(ctx, query, strIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to select tags: %w", err)
	}
	defer rows.Close()
	var tags []openapi.Tag
	for rows.Next() {
		var tag openapi.Tag
		err := rows.Scan(
			&tag.Id,
			&tag.Name,
			&tag.Description,
			&tag.Color,
			&tag.BackgroundColor,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tags = append(tags, tag)
	}
	return &tags, nil
}

func (t TagsRepoImpl) ListByType(ctx context.Context, tagType string, ids *[]string) ([]openapi.Tag, error) {
	var junctionTable string
	var junctionForeignKey string
	switch tagType {
	case "expenditure":
		junctionTable = "expenditure_tags"
		junctionForeignKey = "expenditure_id"
		break
	case "ingress":
		junctionTable = "ingress_tags"
		junctionForeignKey = "ingress_id"
		break
	case "savings_goal":
		junctionTable = "savings_goal_tags"
		junctionForeignKey = "savings_goal_id"
		break
	case "savings_withdrawal":
		junctionTable = "saving_withdrawal_tags"
		junctionForeignKey = "saving_withdrawal_id"
		break
	case "savings_contribution":
		junctionTable = "savings_contribution_tags"
		junctionForeignKey = "savings_contribution_id"
		break
	default:
		return nil, fmt.Errorf("unknown tag type: %s", tagType)
	}
	query := fmt.Sprintf(
		`select 
				t.id,
				t.name,
				t.description,
				t.color,
				t.background_color,
				t.created_at,
				t.updated_at
			from tags t
			where exists (
				select 1
				from %s jt
				where jt.tag_id = t.id
				)`,
		junctionTable)

	if ids != nil {
		query += fmt.Sprintf(` and jt.%s IN (%s)`, junctionForeignKey, strings.Join(*ids, ","))
	}
	query += " order by t.id"
	rows, err := t.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to select tags: %w", err)
	}
	defer rows.Close()
	var tags []openapi.Tag
	for rows.Next() {
		var tag openapi.Tag
		err := rows.Scan(
			&tag.Id,
			&tag.Name,
			&tag.Description,
			&tag.Color,
			&tag.BackgroundColor,
			&tag.CreatedAt,
			&tag.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (t TagsRepoImpl) LinkTagsToType(ctx context.Context, tagType, foreignID string, tags []openapi.Tag) error {
	var junctionTable string
	var junctionForeignKey string
	switch tagType {
	case "expenditure":
		junctionTable = "expenditure_tags"
		junctionForeignKey = "expenditure_id"
		break
	case "ingress":
		junctionTable = "ingress_tags"
		junctionForeignKey = "ingress_id"
		break
	case "savings_goal":
		junctionTable = "savings_goal_tags"
		junctionForeignKey = "savings_goal_id"
		break
	case "savings_withdrawal":
		junctionTable = "saving_withdrawal_tags"
		junctionForeignKey = "saving_withdrawal_id"
		break
	case "savings_contribution":
		junctionTable = "savings_contribution_tags"
		junctionForeignKey = "savings_contribution_id"
		break
	default:
		return fmt.Errorf("unknown tag type: %s", tagType)
	}

	// Clearing records first
	queryDelete := fmt.Sprintf("DELETE FROM %s WHERE %s=?", junctionTable, junctionForeignKey)
	_, err := t.db.ExecContext(
		ctx,
		queryDelete,
		foreignID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete linked records: %w", err)
	}

	// Recreating the relationships
	queryInsert := fmt.Sprintf("INSERT INTO %s (tag_id, %s) VALUES (?,?)", junctionTable, junctionForeignKey)
	for _, tag := range tags {
		_, err := t.db.ExecContext(
			ctx,
			queryInsert,
			tag.Id,
			foreignID,
		)
		if err != nil {
			return fmt.Errorf("failed to link tag to %s: %w", tagType, err)
		}
	}
	return nil
}
