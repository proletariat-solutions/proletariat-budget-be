package resthttp

import (
	"context"
	"errors"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
	"github.com/rs/zerolog/log"
)

func (c *Controller) ListCategories(
	ctx context.Context,
	request openapi.ListCategoriesRequestObject,
) (
	openapi.ListCategoriesResponseObject,
	error,
) {
	categories, err := c.useCases.Category.ListCategories(
		ctx,
		FromOAPICategoryType(request.Params.Type),
	)
	if err != nil {
		log.Err(err).Msg("Failed to list categories")
		return openapi.ListCategories500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Failed to list categories",
			},
		}, nil
	}
	total := len(categories)
	return openapi.ListCategories200JSONResponse{
		Categories: ToOAPICategoryList(&categories),
		Total:      &total,
	}, nil
}

func (c *Controller) CreateCategory(
	ctx context.Context,
	request openapi.CreateCategoryRequestObject,
) (
	openapi.CreateCategoryResponseObject,
	error,
) {
	category, err := c.useCases.Category.CreateCategory(
		ctx,
		*FromOAPICategoryRequest(
			request.Body,
			nil,
		),
	)
	if err != nil {
		log.Err(err).Msg("Failed to create category")
		return openapi.CreateCategory500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Failed to create category",
			},
		}, nil
	}
	return openapi.CreateCategory201JSONResponse(*ToOAPICategory(category)), nil
}

func (c *Controller) DeleteCategory(
	ctx context.Context,
	request openapi.DeleteCategoryRequestObject,
) (
	openapi.DeleteCategoryResponseObject,
	error,
) {
	err := c.useCases.Category.DeleteCategory(
		ctx,
		request.Id,
	)
	if err != nil {
		if errors.Is(
			err,
			domain.ErrCategoryNotFound,
		) {
			return openapi.DeleteCategory404JSONResponse{
				N404JSONResponse: openapi.N404JSONResponse{
					Message: err.Error(),
				},
			}, nil
		} else if errors.Is(
			err,
			domain.ErrCategoryUsedInExpenditure,
		) || errors.Is(
			err,
			domain.ErrCategoryUsedInTransfer,
		) || errors.Is(
			err,
			domain.ErrCategoryUsedInSavingGoal,
		) || errors.Is(
			err,
			domain.ErrCategoryUsedInIngress,
		) || errors.Is(
			err,
			domain.ErrCategoryUsedInEntity,
		) {
			return openapi.DeleteCategory400JSONResponse{
				N400JSONResponse: openapi.N400JSONResponse{
					Message: err.Error(),
				},
			}, nil
		}
		log.Err(err).Msg("Failed to delete category")
		return openapi.DeleteCategory500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Internal Server Error",
			},
		}, nil
	}
	return openapi.DeleteCategory204Response{}, nil
}

func (c *Controller) UpdateCategory(
	ctx context.Context,
	request openapi.UpdateCategoryRequestObject,
) (
	openapi.UpdateCategoryResponseObject,
	error,
) {
	category, err := c.useCases.Category.UpdateCategory(
		ctx,
		*FromOAPICategoryRequest(
			request.Body,
			&request.Id,
		),
	)
	if err != nil {
		if errors.Is(
			err,
			domain.ErrCategoryNotFound,
		) {
			return openapi.UpdateCategory404JSONResponse{
				N404JSONResponse: openapi.N404JSONResponse{
					Message: err.Error(),
				},
			}, nil
		}
		log.Err(err).Msg("Failed to update category")
		return openapi.UpdateCategory500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Internal Server Error",
			},
		}, nil
	}
	return openapi.UpdateCategory200JSONResponse(*ToOAPICategory(category)), nil
}

func (c *Controller) ActivateCategory(
	ctx context.Context,
	request openapi.ActivateCategoryRequestObject,
) (
	openapi.ActivateCategoryResponseObject,
	error,
) {
	err := c.useCases.Category.Activate(
		ctx,
		request.Id,
	)
	if err != nil {
		if errors.Is(
			err,
			domain.ErrCategoryNotFound,
		) {
			return openapi.ActivateCategory404JSONResponse{
				N404JSONResponse: openapi.N404JSONResponse{
					Message: err.Error(),
				},
			}, nil
		}
		if errors.Is(
			err,
			domain.ErrCategoryAlreadyActive,
		) {
			return openapi.ActivateCategory400JSONResponse{
				N400JSONResponse: openapi.N400JSONResponse{
					Message: err.Error(),
				},
			}, nil
		}
		log.Err(err).Msg("Failed to deactivate category")
		return openapi.ActivateCategory500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Internal Server Error",
			},
		}, nil
	}
	return openapi.ActivateCategory204Response{}, nil
}

func (c *Controller) DeactivateCategory(
	ctx context.Context,
	request openapi.DeactivateCategoryRequestObject,
) (
	openapi.DeactivateCategoryResponseObject,
	error,
) {
	err := c.useCases.Category.Deactivate(
		ctx,
		request.Id,
	)
	if err != nil {
		if errors.Is(
			err,
			domain.ErrCategoryNotFound,
		) {
			return openapi.DeactivateCategory404JSONResponse{
				N404JSONResponse: openapi.N404JSONResponse{
					Message: err.Error(),
				},
			}, nil
		}
		if errors.Is(
			err,
			domain.ErrCategoryAlreadyInactive,
		) {
			return openapi.DeactivateCategory400JSONResponse{
				N400JSONResponse: openapi.N400JSONResponse{
					Message: err.Error(),
				},
			}, nil
		}
		log.Err(err).Msg("Failed to deactivate category")
		return openapi.DeactivateCategory500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Internal Server Error",
			},
		}, nil
	}
	return openapi.DeactivateCategory204Response{}, nil
}
