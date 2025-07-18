package port

import (
	"context"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
)

type TagsRepo interface {
	Create(
		ctx context.Context,
		tag domain.Tag,
	) (
		string,
		error,
	)
	Update(
		ctx context.Context,
		id string,
		tag domain.Tag,
	) error
	Delete(
		ctx context.Context,
		id string,
	) error
	GetByID(
		ctx context.Context,
		id string,
	) (
		*domain.Tag,
		error,
	)
	GetByIDs(
		ctx context.Context,
		ids []string,
	) (
		*[]*domain.Tag,
		error,
	)

	ListByType(
		ctx context.Context,
		tagType domain.TagType,
		ids *[]string,
	) (
		*[]*domain.Tag,
		error,
	)

	LinkTagsToType(
		ctx context.Context,
		foreignID string,
		tags *[]*domain.Tag,
	) error

	List(ctx context.Context) (
		*[]*domain.Tag,
		error,
	)

	GetByNameAndType(
		ctx context.Context,
		name string,
		tagType domain.TagType,
	) (
		*domain.Tag,
		error,
	)
}
