package resthttp

import (
	"context"

	"ghorkov32/proletariat-budget-be/openapi"
)

func (c *Controller) ListTransfers(ctx context.Context, request openapi.ListTransfersRequestObject) (openapi.ListTransfersResponseObject, error) {
	panic("implement me")
}

func (c *Controller) CreateTransfer(ctx context.Context, request openapi.CreateTransferRequestObject) (openapi.CreateTransferResponseObject, error) {
	panic("implement me")
}

func (c *Controller) GetTransfer(ctx context.Context, request openapi.GetTransferRequestObject) (openapi.GetTransferResponseObject, error) {
	panic("implement me")
}

func (c *Controller) ListTransactions(ctx context.Context, request openapi.ListTransactionsRequestObject) (openapi.ListTransactionsResponseObject, error) {
	panic("implement me")
}
