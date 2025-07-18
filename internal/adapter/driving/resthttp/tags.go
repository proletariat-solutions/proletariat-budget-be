package resthttp

import (
	"context"
	"errors"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
	"github.com/rs/zerolog/log"
)

func (c *Controller) ListTags(
	ctx context.Context,
	request openapi.ListTagsRequestObject,
) (
	openapi.ListTagsResponseObject,
	error,
) {
	tags, err := c.useCases.Tags.ListTags(
		ctx,
		nil,
	)
	if err != nil {
		return openapi.ListTags500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Internal Server Error",
			},
		}, nil
	}
	var tagList []openapi.Tag
	for _, tag := range tags {
		tagList = append(
			tagList,
			*ToOAPITag(tag),
		)
	}

	return openapi.ListTags200JSONResponse(tagList), nil
}

func (c *Controller) CreateTag(
	ctx context.Context,
	request openapi.CreateTagRequestObject,
) (
	openapi.CreateTagResponseObject,
	error,
) {
	tag, err := c.useCases.Tags.CreateTag(
		ctx,
		FromOAPITagRequest(
			*request.Body,
		),
	)
	if err != nil {
		if errors.Is(
			err,
			domain.ErrTagAlreadyExists,
		) {
			return openapi.CreateTag409JSONResponse{
				N409JSONResponse: openapi.N409JSONResponse{
					Message: err.Error(),
				},
			}, nil
		} else if errors.Is(
			err,
			domain.ErrTagNameEmpty,
		) || errors.Is(
			err,
			domain.ErrUnknownTagType,
		) {
			return openapi.CreateTag400JSONResponse{
				N400JSONResponse: openapi.N400JSONResponse{
					Message: err.Error(),
				},
			}, nil
		}
		log.Err(err).Msg("Failed to create tag")
		return openapi.CreateTag500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Failed to create tag",
			},
		}, nil
	}

	return openapi.CreateTag201JSONResponse(*ToOAPITag(tag)), nil
}

func (c *Controller) DeleteTag(
	ctx context.Context,
	request openapi.DeleteTagRequestObject,
) (
	openapi.DeleteTagResponseObject,
	error,
) {
	err := c.useCases.Tags.DeleteTag(
		ctx,
		request.Id,
	)
	if err != nil {
		if errors.Is(
			err,
			domain.ErrTagNotFound,
		) {
			return openapi.DeleteTag404JSONResponse{
				N404JSONResponse: openapi.N404JSONResponse{
					Message: err.Error(),
				},
			}, nil
		} else {
			log.Err(err).Msg("Failed to delete tag")
			return openapi.DeleteTag500JSONResponse{
				N500JSONResponse: openapi.N500JSONResponse{
					Message: "Failed to delete tag",
				},
			}, nil
		}
	}
	return openapi.DeleteTag204Response{}, nil
}

func (c *Controller) ListTagsByType(
	ctx context.Context,
	request openapi.ListTagsByTypeRequestObject,
) (
	openapi.ListTagsByTypeResponseObject,
	error,
) {
	tags, err := c.useCases.Tags.ListTags(
		ctx,
		FromOAPITagType(&request.Type),
	)
	if err != nil {
		return openapi.ListTagsByType500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Internal Server Error",
			},
		}, nil
	}
	var tagList []openapi.Tag
	for _, tag := range tags {
		tagList = append(
			tagList,
			*ToOAPITag(tag),
		)
	}
	return openapi.ListTagsByType200JSONResponse(tagList), nil
}
