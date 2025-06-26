package port

import (
	"context"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
)

type IngressRepo interface {
	// Ingress operations
	Create(ctx context.Context, ingress openapi.IngressRequest) (string, error)
	Update(ctx context.Context, id string, ingress openapi.IngressRequest) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*openapi.Ingress, error)
	List(ctx context.Context, params openapi.ListIngressesParams) ([]openapi.Ingress, error)

	// Recurrence Patterns
	CreateRecurrencePattern(ctx context.Context, recurrencePattern domain.RecurrencePattern) (string, error)
	UpdateRecurrencePattern(ctx context.Context, id string, recurrencePattern domain.RecurrencePattern) error
	DeleteRecurrencePattern(ctx context.Context, id string) error
	GetRecurrencePattern(ctx context.Context, id string) (*domain.RecurrencePattern, error)

	// Category operations
	ListCategories(ctx context.Context) ([]openapi.IngressCategory, error)
	GetCategory(ctx context.Context, id string) (*openapi.IngressCategory, error)
	CreateCategory(ctx context.Context, category openapi.IngressCategoryRequest) (string, error)
	UpdateCategory(ctx context.Context, id string, category openapi.IngressCategoryRequest) error
	DeleteCategory(ctx context.Context, id string) error

	// Ingress tags operations
	FindOrCreateTags(ctx context.Context, tags []string) ([]string, error)
	LinkTagsToIngress(ctx context.Context, tags []string, ingressId string) error
}
