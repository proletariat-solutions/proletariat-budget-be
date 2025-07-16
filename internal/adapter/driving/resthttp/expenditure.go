package resthttp

import (
	"context"
	"errors"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
	"github.com/rs/zerolog/log"
)

func (c *Controller) ListExpenditures(
	ctx context.Context,
	request openapi.ListExpendituresRequestObject,
) (
	openapi.ListExpendituresResponseObject,
	error,
) {
	list, err := c.useCases.Expenditure.List(
		ctx,
		*FromOAPIExpenditureListParams(&request.Params),
	)
	if err != nil {
		log.Err(err).Msg("Failed to list expenditures")
		return openapi.ListExpenditures500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Failed to list expenditures",
			},
		}, nil
	}
	return openapi.ListExpenditures200JSONResponse(*ToOAPIExpenditureList(list)), nil
}

func (c *Controller) CreateExpenditure(
	ctx context.Context,
	request openapi.CreateExpenditureRequestObject,
) (
	openapi.CreateExpenditureResponseObject,
	error,
) {
	expenditure, err := c.useCases.Expenditure.Create(
		ctx,
		*FromOAPIExpenditureRequest(request.Body),
	)
	if err != nil {
		if errors.Is(
			err,
			domain.ErrCategoryInactive,
		) || errors.Is(
			err,
			domain.ErrAccountInactive,
		) || errors.Is(
			err,
			domain.ErrInsufficientBalance,
		) {
			return openapi.CreateExpenditure409JSONResponse{
				N409JSONResponse: openapi.N409JSONResponse{
					Message: err.Error(),
				},
			}, nil
		}
		if errors.Is(
			err,
			domain.ErrAccountNotFound,
		) || errors.Is(
			err,
			domain.ErrCategoryNotFound,
		) || errors.Is(
			err,
			domain.ErrTagNotFound,
		) {
			return openapi.CreateExpenditure400JSONResponse{
				N400JSONResponse: openapi.N400JSONResponse{
					Message: err.Error(),
				},
			}, nil
		} else {
			log.Err(err).Msg("Failed to create expenditure")
			return openapi.CreateExpenditure500JSONResponse{
				N500JSONResponse: openapi.N500JSONResponse{
					Message: "Failed to create expenditure",
				},
			}, nil
		}
	}
	return openapi.CreateExpenditure201JSONResponse(*ToOAPIExpenditure(expenditure)), nil
}

func (c *Controller) GetExpenditure(
	ctx context.Context,
	request openapi.GetExpenditureRequestObject,
) (
	openapi.GetExpenditureResponseObject,
	error,
) {
	expenditure, err := c.useCases.Expenditure.Get(
		ctx,
		request.Id,
	)
	if err != nil {
		if errors.Is(
			err,
			domain.ErrExpenditureNotFound,
		) {
			return openapi.GetExpenditure404JSONResponse{
				N404JSONResponse: openapi.N404JSONResponse{
					Message: err.Error(),
				},
			}, nil
		} else {
			log.Err(err).Msg("Failed to get expenditure")
			return openapi.GetExpenditure500JSONResponse{
				N500JSONResponse: openapi.N500JSONResponse{
					Message: "Failed to get expenditure",
				},
			}, nil
		}
	}
	return openapi.GetExpenditure200JSONResponse(*ToOAPIExpenditure(expenditure)), nil
}
