package usecase

import (
	"context"
	"errors"
	"strings"

	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
)

type CategoryUseCase struct {
	categoryRepo port.CategoryRepo
}

func NewCategoryUseCase(categoryRepo port.CategoryRepo) *CategoryUseCase {
	return &CategoryUseCase{categoryRepo: categoryRepo}
}

func (uc *CategoryUseCase) ListCategories(
	ctx context.Context,
	categoryType *domain.CategoryType,
) (
	[]domain.Category,
	error,
) {
	var categories []domain.Category
	var err error
	if categoryType == nil {
		categories, err = uc.categoryRepo.List(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		categories, err = uc.categoryRepo.FindByType(
			ctx,
			*categoryType,
		)
		if err != nil {
			return nil, err
		}
	}

	return categories, nil
}

func (uc *CategoryUseCase) GetCategory(
	ctx context.Context,
	id string,
) (
	*domain.Category,
	error,
) {
	category, err := uc.categoryRepo.GetByID(
		ctx,
		id,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return nil, domain.ErrCategoryNotFound
		}

		return nil, err
	}

	return category, nil
}

func (uc *CategoryUseCase) CreateCategory(
	ctx context.Context,
	category domain.Category,
) (
	*domain.Category,
	error,
) {
	id, err := uc.categoryRepo.Create(
		ctx,
		category,
	)
	if err != nil {
		return nil, err
	}

	createdCategory, err := uc.categoryRepo.GetByID(
		ctx,
		id,
	)
	if err != nil {
		return nil, err
	}

	return createdCategory, nil
}

func (uc *CategoryUseCase) UpdateCategory(
	ctx context.Context,
	category domain.Category,
) (
	*domain.Category,
	error,
) {
	_, err := uc.categoryRepo.GetByID(
		ctx,
		category.ID,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return nil, domain.ErrCategoryNotFound
		}

		return nil, err
	}
	err = uc.categoryRepo.Update(
		ctx,
		category,
	)
	if err != nil {
		return nil, err
	}

	updatedCategory, err := uc.categoryRepo.GetByID(
		ctx,
		category.ID,
	)
	if err != nil {
		return nil, err
	}

	return updatedCategory, nil
}

func (uc *CategoryUseCase) DeleteCategory(
	ctx context.Context,
	id string,
) error {
	err := uc.categoryRepo.Delete(
		ctx,
		id,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return domain.ErrCategoryNotFound
		}
		if port.IsInfrastructureError(err) {
			return uc.getErrorConstraintType(err)
		}
	}

	return nil
}

func (uc *CategoryUseCase) getErrorConstraintType(err error) error {
	if strings.Contains(
		err.Error(),
		"expenditures",
	) {
		return domain.ErrCategoryUsedInExpenditure
	} else if strings.Contains(
		err.Error(),
		"ingress",
	) {
		return domain.ErrCategoryUsedInIngress
	} else if strings.Contains(
		err.Error(),
		"transfer",
	) {
		return domain.ErrCategoryUsedInTransfer
	} else if strings.Contains(
		err.Error(),
		"saving_goal",
	) {
		return domain.ErrCategoryUsedInSavingGoal
	}

	return domain.ErrCategoryUsedInEntity
}

func (uc *CategoryUseCase) Activate(
	ctx context.Context,
	id string,
) error {
	category, err := uc.categoryRepo.GetByID(
		ctx,
		id,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return domain.ErrCategoryNotFound
		}

		return err
	}
	err = category.Activate()
	if err != nil {
		return err
	}
	err = uc.categoryRepo.Update(
		ctx,
		*category,
	)
	if err != nil {
		return err
	}

	return nil
}

func (uc *CategoryUseCase) Deactivate(
	ctx context.Context,
	id string,
) error {
	category, err := uc.categoryRepo.GetByID(
		ctx,
		id,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return domain.ErrCategoryNotFound
		}

		return err
	}
	err = category.Deactivate()
	if err != nil {
		return err
	}
	err = uc.categoryRepo.Update(
		ctx,
		*category,
	)
	if err != nil {
		return err
	}

	return nil
}
