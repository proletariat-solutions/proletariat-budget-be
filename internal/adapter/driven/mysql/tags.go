package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"strings"
)

type TagsRepoImpl struct {
	db *sql.DB
}

func NewTagsRepo(db *sql.DB) port.TagsRepo {
	return &TagsRepoImpl{
		db: db,
	}
}

func (t TagsRepoImpl) Create(
	ctx context.Context,
	tag domain.Tag,
) (
	string,
	error,
) {
	queryInsert := `INSERT INTO tags (name, description, color, background_color, type, created_at) VALUES (?,?,?,?,?, now())`
	result, errInsert := t.db.ExecContext(
		ctx,
		queryInsert,
		tag.Name,
		tag.Description,
		tag.Color,
		tag.BackgroundColor,
		tag.TagType,
	)
	if errInsert != nil {
		return "", fmt.Errorf(
			"failed to create tag: %w",
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
	lastIDStr := fmt.Sprintf(
		"%d",
		lastID,
	)
	return lastIDStr, nil
}

func (t TagsRepoImpl) GetByNameAndType(
	ctx context.Context,
	name string,
	tagType domain.TagType,
) (
	*domain.Tag,
	error,
) {
	querySelect := `SELECT id, name, description, color, background_color, type FROM tags WHERE name=? AND type=?`
	row := t.db.QueryRowContext(
		ctx,
		querySelect,
		name,
		tagType,
	)
	var tag domain.Tag
	err := row.Scan(
		&tag.ID,
		&tag.Name,
		&tag.Description,
		&tag.Color,
		&tag.BackgroundColor,
		&tag.TagType,
	)
	if errors.Is(
		err,
		sql.ErrNoRows,
	) {
		return nil, port.ErrRecordNotFound
	}
	if err != nil {
		return nil, translateError(err)
	}
	return &tag, nil
}
func (t TagsRepoImpl) Update(
	ctx context.Context,
	id string,
	tag domain.Tag,
) error {
	queryUpdate := `UPDATE tags SET name=?, description=?, color=?, background_color=?, type=? WHERE id=?`
	_, err := t.db.ExecContext(
		ctx,
		queryUpdate,
		tag.Name,
		tag.Description,
		tag.Color,
		tag.BackgroundColor,
		tag.TagType,
		id,
	)
	if err != nil {
		return translateError(err)
	}
	return nil
}

func (t TagsRepoImpl) Delete(
	ctx context.Context,
	id string,
) error {
	tag, err := t.GetByID(
		ctx,
		id,
	)
	if err != nil {
		return err
	}

	junctionTable, _, err := t.getJunctionTableByType(tag.TagType)
	if err != nil {
		return translateError(err)
	}
	if junctionTable != nil {
		// Deleting linked records first
		queryDeleteLinked := fmt.Sprintf(
			`DELETE FROM %s WHERE tag_id=?`,
			*junctionTable,
		)
		_, err := t.db.ExecContext(
			ctx,
			queryDeleteLinked,
			id,
		)
		if err != nil {
			return translateError(err)
		}
	}

	queryDelete := `DELETE FROM tags WHERE id=?`
	_, err = t.db.ExecContext(
		ctx,
		queryDelete,
		id,
	)
	if err != nil {
		return translateError(err)
	}
	return nil

}

func (t TagsRepoImpl) GetByID(
	ctx context.Context,
	id string,
) (
	*domain.Tag,
	error,
) {
	query := `SELECT id, name, description, color, background_color, type FROM tags WHERE id=?`
	var tag domain.Tag
	err := t.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&tag.ID,
		&tag.Name,
		&tag.Description,
		&tag.Color,
		&tag.BackgroundColor,
		&tag.TagType,
	)
	if err != nil {
		return nil, translateError(err)
	}
	return &tag, nil
}

func (t TagsRepoImpl) GetByIDs(
	ctx context.Context,
	ids []string,
) (
	*[]*domain.Tag,
	error,
) {
	query := `SELECT id, name, description, color, background_color, type FROM tags WHERE id IN (?)`
	strIDs := strings.Join(
		ids,
		",",
	)
	rows, err := t.db.QueryContext(
		ctx,
		query,
		strIDs,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to select tags: %w",
			err,
		)
	}
	defer rows.Close()
	var tags []*domain.Tag
	for rows.Next() {
		var tag domain.Tag
		err = rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Description,
			&tag.Color,
			&tag.BackgroundColor,
			&tag.TagType,
		)
		if err != nil {
			return nil, translateError(err)
		}
		tags = append(
			tags,
			&tag,
		)
	}
	return &tags, nil
}

func (t TagsRepoImpl) List(
	ctx context.Context,
) (
	*[]*domain.Tag,
	error,
) {
	query := `SELECT id, name, description, color, background_color, type FROM tags`
	var tags []*domain.Tag
	res, err := t.db.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, translateError(err)
	}
	defer res.Close()
	for res.Next() {
		var tag domain.Tag
		err = res.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Description,
			&tag.Color,
			&tag.BackgroundColor,
			&tag.TagType,
		)
		if err != nil {
			return nil, translateError(err)
		}
		tags = append(
			tags,
			&tag,
		)
	}
	return &tags, nil

}

func (t TagsRepoImpl) ListByType(
	ctx context.Context,
	tagType domain.TagType,
	ids *[]string,
) (
	*[]*domain.Tag,
	error,
) {

	// We'll assume that all the tags are the same type
	junctionTable, junctionForeignKey, err := t.getJunctionTableByType(tagType)
	if err != nil {
		return nil, err
	}

	query :=
		`select 
				t.id,
				t.name,
				t.description,
				t.color,
				t.background_color,
				t.type
			from tags t
			where t.type =?`

	if ids != nil {
		query += fmt.Sprintf(
			` AND exists (
				select 1
				from %s jt
				where jt.tag_id = t.id
				)`,
			junctionTable,
		)
		query += fmt.Sprintf(
			` and jt.%s IN (%s)`,
			*junctionForeignKey,
			strings.Join(
				*ids,
				",",
			),
		)
	}
	query += " order by t.id"
	rows, err := t.db.QueryContext(
		ctx,
		query,
		tagType,
	)
	if err != nil {
		return nil, translateError(err)
	}
	defer rows.Close()
	var tags []*domain.Tag
	for rows.Next() {
		var tag domain.Tag
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Description,
			&tag.Color,
			&tag.BackgroundColor,
			&tag.TagType,
		)
		if err != nil {
			return nil, translateError(err)
		}
		tags = append(
			tags,
			&tag,
		)
	}
	return &tags, nil
}

func (t TagsRepoImpl) LinkTagsToType(
	ctx context.Context,
	foreignID string,
	tags *[]*domain.Tag,
) error {
	// We'll assume that all the tags are the same type
	junctionTable, junctionForeignKey, err := t.getJunctionTableByType((*tags)[0].TagType)
	if err != nil {
		return err
	}

	// Clearing records first
	queryDelete := fmt.Sprintf(
		"DELETE FROM %s WHERE %s=?",
		*junctionTable,
		*junctionForeignKey,
	)
	_, err = t.db.ExecContext(
		ctx,
		queryDelete,
		foreignID,
	)
	if err != nil {
		return translateError(err)
	}

	// Recreating the relationships
	queryInsert := fmt.Sprintf(
		"INSERT INTO %s (tag_id, %s) VALUES (?,?)",
		junctionTable,
		junctionForeignKey,
	)
	for _, tag := range *tags {
		_, err = t.db.ExecContext(
			ctx,
			queryInsert,
			(*tag).ID,
			foreignID,
		)
		if err != nil {
			return translateError(err)
		}
	}
	return nil
}

func (t TagsRepoImpl) getJunctionTableByType(tagType domain.TagType) (
	*string,
	*string,
	error,
) {
	var foreignKey string
	var junctionTable string
	switch tagType {
	case "expenditure":
		junctionTable = "expenditure_tags"
		foreignKey = "expenditure_id"
		break
	case "ingress":
		junctionTable = "ingress_tags"
		foreignKey = "ingress_id"
		break
	case "saving_goal":
		junctionTable = "savings_goal_tags"
		foreignKey = "savings_goal_id"
		break
	case "savings_withdrawal":
		junctionTable = "saving_withdrawal_tags"
		foreignKey = "saving_withdrawal_id"
		break
	case "savings_contribution":
		junctionTable = "savings_contribution_tags"
		foreignKey = "savings_contribution_id"
		break
	default:
		return nil, nil, domain.ErrUnknownTagType
	}
	return &junctionTable, &foreignKey, nil
}
