package port

import (
	"context"

	"ghorkov32/proletariat-budget-be/openapi"
)

type IngressRepo interface {
	// Ingress operations
	Create(ctx context.Context, ingress openapi.IngressRequest, transactionID string) (string, error)
	GetByID(ctx context.Context, id string) (*openapi.Ingress, error)
	List(ctx context.Context, params openapi.ListIngressesParams) (*openapi.IngressList, error)

	// Recurrence Patterns
	CreateRecurrencePattern(ctx context.Context, recurrencePattern openapi.RecurrencePatternRequest) (string, error)
	UpdateRecurrencePattern(ctx context.Context, id string, recurrencePattern openapi.RecurrencePattern) error
	DeleteRecurrencePattern(ctx context.Context, id string) error
	GetRecurrencePattern(ctx context.Context, id string) (*openapi.RecurrencePattern, error)
}
