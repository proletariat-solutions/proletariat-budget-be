package port

import (
	"context"
	"ghorkov32/proletariat-budget-be/openapi"
)

type TagsRepo interface {
	Create(ctx context.Context, tag openapi.Tag, tagType string) (string, error)
	Update(ctx context.Context, id string, tag openapi.Tag, tagType string) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*openapi.Tag, error)

	ListByType(ctx context.Context, tagType string, ids *[]string) ([]openapi.Tag, error)

	LinkTagsToType(ctx context.Context, tagType string, tags []openapi.Tag) error
}
