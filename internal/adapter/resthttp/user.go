package resthttp

import (
	"context"
	"errors"

	"ghorkov32/proletariat-budget-be/internal/common"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"ghorkov32/proletariat-budget-be/internal/core/usecase"
	"ghorkov32/proletariat-budget-be/openapi"
	"github.com/inv-cloud-platform/hub-com-tools-go/hubmiddlewares"
	"github.com/rs/zerolog/log"
)

// Check if operation satisfies interface.
var _ openapi.StrictServerInterface = (*UserController)(nil)

type UserController struct {
	userSvc       *usecase.CreateUser
	accessChecker port.AccessCheck
}

func NewUserController(userSvc *usecase.CreateUser, accessChecker port.AccessCheck) *UserController {
	return &UserController{
		userSvc:       userSvc,
		accessChecker: accessChecker,
	}
}

func (c *UserController) CreateUser(ctx context.Context, req openapi.CreateUserRequestObject) (openapi.CreateUserResponseObject, error) {
	requestID := hubmiddlewares.GetRequestId(ctx)

	if _, errCheck := c.accessChecker.Check(ctx, common.Create); errCheck != nil {
		return openapi.CreateUser403JSONResponse{
			N403JSONResponse: openapi.N403JSONResponse{
				Code:    openapi.ErrCodeForbidden,
				Message: "not allowed to create user",
			},
		}, nil
	}

	err := c.userSvc.Create(ctx, domain.User{
		Name: req.Body.Name,
	})

	switch {
	case errors.Is(err, domain.ErrUserAlreadyExists):
		log.Warn().Err(err).Str(common.RequestID, requestID).Msg("user already exists")

		return openapi.CreateUser409JSONResponse{
			N409JSONResponse: openapi.N409JSONResponse{
				Code:    openapi.ErrCodeConflict,
				Message: "resource already exists",
			}}, err
	case err != nil:
		log.Error().Err(err).Str(common.RequestID, requestID).Msg("unable to create user")

		return openapi.CreateUser500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Code:    openapi.ErrorCodeInternalServerError,
				Message: "unexpected error occurred",
			}}, err
	default:
		return openapi.CreateUser200JSONResponse{}, nil
	}
}
