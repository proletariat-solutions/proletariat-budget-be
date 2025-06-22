package port

import (
	"context"
	"ghorkov32/proletariat-budget-be/openapi"
)

type IngressRepo interface {
	// Ingress operations
	Create(ctx context.Context, ingress openapi.IngressRequest) (string, error)
	Update(ctx context.Context, id string, ingress openapi.IngressRequest) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*openapi.Ingress, error)
	List(ctx context.Context, params openapi.ListIngressesParams) ([]openapi.Ingress, error)

	// Category operations
	ListCategories(ctx context.Context) ([]string, error)
	GetCategory(ctx context.Context, id string) (*openapi.IngressCategory, error)
	CreateCategory(ctx context.Context, category openapi.IngressCategoryRequest) (string, error)
	UpdateCategory(ctx context.Context, id string, category openapi.IngressCategoryRequest) error
	DeleteCategory(ctx context.Context, id string) error

	// Ingress tags operations
	FindOrCreateTags(ctx context.Context, tags []string) ([]string, error)
	LinkTagsToIngress(ctx context.Context, tags []string, ingressId string) error

	// Ingress transactions operations
	ListTransactions(ctx context.Context, params openapi.ListTransactionsParams) (*openapi.TransactionList, error)


}
