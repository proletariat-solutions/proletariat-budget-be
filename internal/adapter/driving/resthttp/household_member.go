package resthttp

import (
	"context"
	"ghorkov32/proletariat-budget-be/openapi"
	"github.com/rs/zerolog/log"
)

func (c *Controller) ListHouseholdMembers(ctx context.Context, request openapi.ListHouseholdMembersRequestObject) (openapi.ListHouseholdMembersResponseObject, error) {
	list, err := c.useCases.HouseholdMember.ListHouseholdMembers(ctx)
	if err != nil {
		return openapi.ListHouseholdMembers500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Internal Server Error",
			},
		}, err
	}
	return openapi.ListHouseholdMembers200JSONResponse(*ToOAPIHouseholdMemberList(list)), nil
}

func (c *Controller) CreateHouseholdMember(ctx context.Context, request openapi.CreateHouseholdMemberRequestObject) (openapi.CreateHouseholdMemberResponseObject, error) {
	id, err := c.useCases.HouseholdMember.CreateHouseholdMember(ctx, *FromOAPIHouseholdMemberRequest(request.Body))
	if err != nil {
		return openapi.CreateHouseholdMember500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Internal Server Error",
			},
		}, err
	}

	householdMember, err := c.useCases.HouseholdMember.GetHouseholdMemberByID(ctx, id)
	if err != nil {
		return openapi.CreateHouseholdMember500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Internal Server Error",
			},
		}, err
	}
	return openapi.CreateHouseholdMember201JSONResponse(*ToOAPIHouseholdMember(householdMember)), nil
}

func (c *Controller) DeleteHouseholdMember(ctx context.Context, request openapi.DeleteHouseholdMemberRequestObject) (openapi.DeleteHouseholdMemberResponseObject, error) {
	err := c.useCases.HouseholdMember.DeleteHouseholdMember(ctx, request.Id)
	if err != nil {
		/*if errors.Is(err, domain.ErrEntityNotFound) {
			return openapi.DeleteHouseholdMember404JSONResponse{
				N404JSONResponse: openapi.N404JSONResponse{
					Message: "Household member not found",
				},
			}, nil
		} else {*/
		log.Err(err).Msg("Failed to delete household member")
		return openapi.DeleteHouseholdMember500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Internal Server Error",
			},
		}, err
		//}
	}
	return openapi.DeleteHouseholdMember204Response{}, nil
}

func (c *Controller) GetHouseholdMember(ctx context.Context, request openapi.GetHouseholdMemberRequestObject) (openapi.GetHouseholdMemberResponseObject, error) {
	member, err := c.useCases.HouseholdMember.GetHouseholdMemberByID(ctx, request.Id)
	if err != nil {
		/*if errors.Is(err, domain.ErrEntityNotFound) {
			return openapi.GetHouseholdMember404JSONResponse{
				N404JSONResponse: openapi.N404JSONResponse{
					Message: "Household member not found",
				},
			}, nil
		} else {*/
		log.Err(err).Msg("Failed to get household member")
		return openapi.GetHouseholdMember500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Internal Server Error",
			},
		}, err
		//}
	}
	return openapi.GetHouseholdMember200JSONResponse(*ToOAPIHouseholdMember(member)), nil
}

func (c *Controller) UpdateHouseholdMember(ctx context.Context, request openapi.UpdateHouseholdMemberRequestObject) (openapi.UpdateHouseholdMemberResponseObject, error) {
	err := c.useCases.HouseholdMember.UpdateHouseholdMember(ctx, request.Id, *FromOAPIHouseholdMemberRequest(request.Body))
	if err != nil {
		/*if errors.Is(err, domain.ErrEntityNotFound) {
			return openapi.UpdateHouseholdMember404JSONResponse{
				N404JSONResponse: openapi.N404JSONResponse{
					Message: "Household member not found",
				},
			}, nil
		} else {*/
		log.Err(err).Msg("Failed to update household member")
		return openapi.UpdateHouseholdMember500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Internal Server Error",
			},
		}, err
		//}
	}
	member, err := c.useCases.HouseholdMember.GetHouseholdMemberByID(ctx, request.Id)
	if err != nil {
		return openapi.UpdateHouseholdMember500JSONResponse{
			N500JSONResponse: openapi.N500JSONResponse{
				Message: "Internal Server Error",
			},
		}, err
	}
	return openapi.UpdateHouseholdMember200JSONResponse(*ToOAPIHouseholdMember(member)), nil
}
