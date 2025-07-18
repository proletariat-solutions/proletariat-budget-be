package port

import (
	"context"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
)

type CategoryRepo interface {
	Create(
		ctx context.Context,
		category domain.Category,
	) (
		string,
		error,
	)
	Update(
		ctx context.Context,
		category domain.Category,
	) error
	Delete(
		ctx context.Context,
		id string,
	) error
	GetByID(
		ctx context.Context,
		id string,
	) (
		*domain.Category,
		error,
	)
	List(ctx context.Context) (
		[]domain.Category,
		error,
	)
	FindByType(
		ctx context.Context,
		categoryType domain.CategoryType,
	) (
		[]domain.Category,
		error,
	)
	FindByIDs(
		ctx context.Context,
		ids []string,
	) (
		[]domain.Category,
		error,
	)
}
