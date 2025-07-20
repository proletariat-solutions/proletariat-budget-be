package resthttp

import (
	"context"

	"ghorkov32/proletariat-budget-be/openapi"
)

func (c *Controller) RollbackExpenditure(ctx context.Context, request openapi.RollbackExpenditureRequestObject) (openapi.RollbackExpenditureResponseObject, error) {
	panic("implement me")
}

func (c *Controller) RollbackIngress(ctx context.Context, request openapi.RollbackIngressRequestObject) (openapi.RollbackIngressResponseObject, error) {
	panic("implement me")
}

func (c *Controller) RollbackTransfer(ctx context.Context, request openapi.RollbackTransferRequestObject) (openapi.RollbackTransferResponseObject, error) {
	panic("implement me")
}
