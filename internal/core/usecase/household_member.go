package usecase

import (
	"context"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
)

type HouseholdMemberUseCase struct {
	householdMembersRepo port.HouseholdMembersRepo
}

func NewHouseholdMemberUseCase(householdMembersRepo port.HouseholdMembersRepo) *HouseholdMemberUseCase {
	return &HouseholdMemberUseCase{householdMembersRepo: householdMembersRepo}
}

func (u *HouseholdMemberUseCase) ListHouseholdMembers(ctx context.Context) (*domain.HouseholdMemberList, error) {
	return u.householdMembersRepo.List(ctx)
}

func (u *HouseholdMemberUseCase) CreateHouseholdMember(ctx context.Context, householdMember domain.HouseholdMember) (string, error) {
	return u.householdMembersRepo.Create(ctx, householdMember)
}

func (u *HouseholdMemberUseCase) UpdateHouseholdMember(ctx context.Context, id string, householdMember domain.HouseholdMember) error {
	return u.householdMembersRepo.Update(ctx, id, householdMember)
}

func (u *HouseholdMemberUseCase) DeleteHouseholdMember(ctx context.Context, id string) error {
	return u.householdMembersRepo.Delete(ctx, id)
}

func (u *HouseholdMemberUseCase) GetHouseholdMemberByID(ctx context.Context, id string) (*domain.HouseholdMember, error) {
	return u.householdMembersRepo.GetByID(ctx, id)
}
