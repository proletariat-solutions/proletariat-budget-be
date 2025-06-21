package resthttp

import (
	"ghorkov32/proletariat-budget-be/internal/core/usecase"
	"ghorkov32/proletariat-budget-be/openapi"
)

var _ openapi.StrictServerInterface = (*Controller)(nil)

type Controller struct {
	authUseCase *usecase.AuthUseCase
}
