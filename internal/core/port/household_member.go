package port

import (
	"context"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
)

type HouseholdMembersRepo interface {
	Create(ctx context.Context, householdMember domain.HouseholdMember) (string, error)
	Update(ctx context.Context, id string, householdMember domain.HouseholdMember) error
	Delete(ctx context.Context, id string) error
	Deactivate(ctx context.Context, id string) error
	Activate(ctx context.Context, id string) error
	CanDelete(ctx context.Context, id string) (bool, error)
	GetByID(ctx context.Context, id string) (*domain.HouseholdMember, error)
	List(ctx context.Context, params *domain.HouseholdMemberListParams) (*domain.HouseholdMemberList, error)
}
