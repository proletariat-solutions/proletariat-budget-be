package resthttp

import (
	"ghorkov32/proletariat-budget-be/internal/core/usecase"
	"ghorkov32/proletariat-budget-be/openapi"
)

type Controller struct {
	useCases usecase.UseCases
}

func NewController(useCases usecase.UseCases) openapi.StrictServerInterface {
	return &Controller{useCases: useCases}
}
