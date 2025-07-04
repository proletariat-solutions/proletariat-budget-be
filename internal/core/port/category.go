package port

import (
	"context"
	"ghorkov32/proletariat-budget-be/openapi"
)

type CategoryRepo interface {
	Create(ctx context.Context, category openapi.Category, categoryType string) (string, error)
	Update(ctx context.Context, id string, category openapi.Category, categoryType string) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*openapi.Category, error)
	List(ctx context.Context) ([]openapi.Category, error)
	FindByType(ctx context.Context, categoryType string) ([]openapi.Category, error)
	FindByIDs(ctx context.Context, ids []string) ([]openapi.Category, error)
}
