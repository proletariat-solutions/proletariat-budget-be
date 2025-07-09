package port

import (
	"context"
	"ghorkov32/proletariat-budget-be/openapi"
)

type HouseholdMembersRepo interface {
	Create(ctx context.Context, householdMember openapi.HouseholdMemberRequest) (string, error)
	Update(ctx context.Context, id string, householdMember openapi.HouseholdMemberRequest) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*openapi.HouseholdMember, error)
	List(ctx context.Context) (*openapi.HouseholdMemberList, error)
}
