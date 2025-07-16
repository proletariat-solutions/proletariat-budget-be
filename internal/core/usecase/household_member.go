package usecase

import (
	"context"
	"errors"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
)

type HouseholdMemberUseCase struct {
	householdMembersRepo port.HouseholdMembersRepo
}

func NewHouseholdMemberUseCase(householdMembersRepo port.HouseholdMembersRepo) *HouseholdMemberUseCase {
	return &HouseholdMemberUseCase{householdMembersRepo: householdMembersRepo}
}

func (u *HouseholdMemberUseCase) ListHouseholdMembers(
	ctx context.Context,
	params *domain.HouseholdMemberListParams,
) (
	*domain.HouseholdMemberList,
	error,
) {
	return u.householdMembersRepo.List(
		ctx,
		params,
	)
}

func (u *HouseholdMemberUseCase) CreateHouseholdMember(
	ctx context.Context,
	householdMember domain.HouseholdMember,
) (
	string,
	error,
) {
	return u.householdMembersRepo.Create(
		ctx,
		householdMember,
	)
}

func (u *HouseholdMemberUseCase) UpdateHouseholdMember(
	ctx context.Context,
	id string,
	householdMember domain.HouseholdMember,
) error {
	_, err := u.householdMembersRepo.GetByID(
		ctx,
		id,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return domain.ErrMemberNotFound
		}
		return err
	}
	return u.householdMembersRepo.Update(
		ctx,
		id,
		householdMember,
	)
}

func (u *HouseholdMemberUseCase) DeleteHouseholdMember(
	ctx context.Context,
	id string,
) error {
	_, err := u.householdMembersRepo.GetByID(
		ctx,
		id,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return domain.ErrMemberNotFound
		}
		return err
	}
	canDelete, err := u.householdMembersRepo.CanDelete(
		ctx,
		id,
	)
	if err != nil {
		return err
	}
	if !canDelete {
		return domain.ErrMemberHasActiveAccounts
	}
	return u.householdMembersRepo.Delete(
		ctx,
		id,
	)
}

func (u *HouseholdMemberUseCase) GetHouseholdMemberByID(
	ctx context.Context,
	id string,
) (
	*domain.HouseholdMember,
	error,
) {
	member, err := u.householdMembersRepo.GetByID(
		ctx,
		id,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return nil, domain.ErrMemberNotFound
		}
		return nil, err
	}
	return member, nil
}

func (u *HouseholdMemberUseCase) DeactivateHouseholdMember(
	ctx context.Context,
	id string,
) error {
	member, err := u.householdMembersRepo.GetByID(
		ctx,
		id,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return domain.ErrMemberNotFound
		}
		return err
	}
	if !member.Active {
		return domain.ErrMemberAlreadyInactive
	}
	return u.householdMembersRepo.Deactivate(
		ctx,
		id,
	)
}

func (u *HouseholdMemberUseCase) CanDeleteHouseholdMember(
	ctx context.Context,
	id string,
) (
	bool,
	error,
) {
	_, err := u.householdMembersRepo.GetByID(
		ctx,
		id,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return false, domain.ErrMemberNotFound
		}
		return false, err
	}
	return u.householdMembersRepo.CanDelete(
		ctx,
		id,
	)
}

func (u *HouseholdMemberUseCase) ActivateHouseholdMember(
	ctx context.Context,
	id string,
) error {
	member, err := u.householdMembersRepo.GetByID(
		ctx,
		id,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return domain.ErrMemberNotFound
		}
		return err
	}
	if member.Active {
		return domain.ErrMemberAlreadyActive
	}
	return u.householdMembersRepo.Activate(
		ctx,
		id,
	)
}
