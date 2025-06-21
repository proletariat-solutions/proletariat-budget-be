package mysql

import (
	"context"
	"database/sql"
	"ghorkov32/proletariat-budget-be/openapi"
)

type IngressRepo struct {
	db *sql.DB
}

func (i IngressRepo) Create(ctx context.Context, ingress openapi.IngressRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (i IngressRepo) Update(ctx context.Context, id string, ingress openapi.IngressRequest) error {
	//TODO implement me
	panic("implement me")
}

func (i IngressRepo) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (i IngressRepo) GetByID(ctx context.Context, id string) (*openapi.Ingress, error) {
	//TODO implement me
	panic("implement me")
}

func (i IngressRepo) List(ctx context.Context, params openapi.ListIngressesParams) ([]openapi.Ingress, error) {
	//TODO implement me
	panic("implement me")
}

func (i IngressRepo) ListCategories(ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (i IngressRepo) GetCategory(ctx context.Context, id string) (*openapi.IngressCategory, error) {
	//TODO implement me
	panic("implement me")
}

func (i IngressRepo) CreateCategory(ctx context.Context, category openapi.IngressCategoryRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (i IngressRepo) UpdateCategory(ctx context.Context, id string, category openapi.IngressCategoryRequest) error {
	//TODO implement me
	panic("implement me")
}

func (i IngressRepo) DeleteCategory(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (i IngressRepo) FindOrCreateTags(ctx context.Context, tags []string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (i IngressRepo) LinkTagsToIngress(ctx context.Context, tags []string, ingressId string) error {
	//TODO implement me
	panic("implement me")
}

func (i IngressRepo) ListTransactions(ctx context.Context, params openapi.ListTransactionsParams) (*openapi.TransactionList, error) {
	//TODO implement me
	panic("implement me")
}
