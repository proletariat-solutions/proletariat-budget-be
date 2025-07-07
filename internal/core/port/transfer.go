package port

import (
	"context"
	"ghorkov32/proletariat-budget-be/openapi"
)

type TransferRepo interface {
	Create(ctx context.Context, transfer openapi.Transfer, incomingTxID, outgoingTxID string) (string, error)
	GetByID(ctx context.Context, id string) (*openapi.Transfer, error)
	List(ctx context.Context, params openapi.ListTransfersParams) (*openapi.TransferList, error)
}
