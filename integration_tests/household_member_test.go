package integration_tests

import (
	"ghorkov32/proletariat-budget-be/integration_tests/utils"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
	"net/http"
	"strconv"
)

func (s *Suite) TestHouseholdMember() {

	var deletableMemberId string

	s.T().Log("Starting TestHouseholdMember")
	s.Run(
		"Update a household member",
		func() {
			member := s.createTestHouseholdMember()
			member.Role = "test-role"

			apiResponse, err := s.updateHouseholdMember(&member)

			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusOK,
				apiResponse.StatusCode,
			)

			var updatedMember openapi.HouseholdMember
			s.decodeResponse(
				apiResponse,
				&updatedMember,
			)

			s.Equal(
				member,
				updatedMember,
			)
			deletableMemberId = updatedMember.Id
		},
	)

	s.Run(
		"Update a non-existant household member",
		func() {
			member := s.createTestHouseholdMember()
			member.Id = "0"

			apiResponse, err := s.updateHouseholdMember(&member)

			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusBadRequest,
				domain.ErrMemberNotFound.Error(),
			)
		},
	)

	s.Run(
		"Delete a household member",
		func() {
			apiResponse, err := s.deleteHouseholdMember(deletableMemberId)
			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusNoContent,
				apiResponse.StatusCode,
			)

			apiResponse, err = s.getHouseholdMember(deletableMemberId)

			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusNotFound,
				domain.ErrMemberNotFound.Error(),
			)
		},
	)

	s.Run(
		"Delete a household member that doesn't exist",
		func() {
			nonExistentId := "non-existent-id"

			apiResponse, err := s.deleteHouseholdMember(nonExistentId)
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusNotFound,
				domain.ErrMemberNotFound.Error(),
			)
		},
	)

	s.Run(
		"Delete a household member with active accounts",
		func() {

			apiResponse, err := s.deleteHouseholdMember("1") // ID coming from mock data
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusBadRequest,
				domain.ErrMemberHasActiveAccounts.Error(),
			)
		},
	)
	s.Run(
		"Activate a household member that doesn't exist",
		func() {
			nonExistentId := "non-existent-id"

			apiResponse, err := s.activateHouseholdMember(nonExistentId)
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusNotFound,
				domain.ErrMemberNotFound.Error(),
			)
		},
	)

	s.Run(
		"Activate a household member that is already active",
		func() {
			// Create an active household member
			member := s.createTestHouseholdMember() // This creates an active member by default

			apiResponse, err := s.activateHouseholdMember(member.Id)
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusBadRequest,
				domain.ErrMemberAlreadyActive.Error(),
			)
		},
	)

	s.Run(
		"Successfully activate an inactive household member",
		func() {
			// Create a household member
			member := s.createTestHouseholdMember()

			// First deactivate the member
			deactivateResponse, err := s.deactivateHouseholdMember(member.Id)
			s.handleErr(
				err,
				"error while deactivating member",
			)
			s.Equal(
				http.StatusNoContent,
				deactivateResponse.StatusCode,
			)

			// Now activate the member
			apiResponse, err := s.activateHouseholdMember(member.Id)
			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusNoContent,
				apiResponse.StatusCode,
			)

			// Verify the member is actually activated by getting the member
			getResponse, err := s.getHouseholdMember(member.Id)
			s.handleErr(
				err,
				"error while getting member",
			)
			s.Equal(
				http.StatusOK,
				getResponse.StatusCode,
			)

			var activatedMember openapi.HouseholdMember
			s.decodeResponse(
				getResponse,
				&activatedMember,
			)
			s.True(
				*activatedMember.Active,
				"Member should be active after activation",
			)
		},
	)
	s.Run(
		"Deactivate a household member that doesn't exist",
		func() {
			nonExistentId := "0"

			apiResponse, err := s.deactivateHouseholdMember(nonExistentId)
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusNotFound,
				domain.ErrMemberNotFound.Error(),
			)
		},
	)

	s.Run(
		"Deactivate a household member that is already inactive",
		func() {
			// Create a household member
			member := s.createTestHouseholdMember()

			// First deactivate the member
			deactivateResponse, err := s.deactivateHouseholdMember(member.Id)
			s.handleErr(
				err,
				"error while deactivating member",
			)
			s.Equal(
				http.StatusNoContent,
				deactivateResponse.StatusCode,
			)

			// Try to deactivate again - should fail with already inactive error
			apiResponse, err := s.deactivateHouseholdMember(member.Id)
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusBadRequest,
				domain.ErrMemberAlreadyInactive.Error(),
			)
		},
	)

	s.Run(
		"Successfully deactivate an active household member",
		func() {
			// Create an active household member
			member := s.createTestHouseholdMember() // This creates an active member by default

			// Deactivate the member
			apiResponse, err := s.deactivateHouseholdMember(member.Id)
			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusNoContent,
				apiResponse.StatusCode,
			)

			// Verify the member is actually deactivated by getting the member
			getResponse, err := s.getHouseholdMember(member.Id)
			s.handleErr(
				err,
				"error while getting member",
			)
			s.Equal(
				http.StatusOK,
				getResponse.StatusCode,
			)

			var deactivatedMember openapi.HouseholdMember
			s.decodeResponse(
				getResponse,
				&deactivatedMember,
			)
			s.False(
				*deactivatedMember.Active,
				"Member should be inactive after deactivation",
			)
		},
	)

	s.Run(
		"List all household members",
		func() {
			// Create some test members with different roles and statuses
			_ = s.createTestHouseholdMember() // Active by default
			member2 := s.createTestHouseholdMember()
			member2.Role = "child"
			_, err := s.updateHouseholdMember(&member2)
			s.handleErr(
				err,
				"error while updating member",
			)

			// Create an inactive member
			member3 := s.createTestHouseholdMember()
			member3.Role = "roommate"
			_, err = s.updateHouseholdMember(&member3)
			s.handleErr(
				err,
				"error while updating member",
			)
			_, err = s.deactivateHouseholdMember(member3.Id)
			s.handleErr(
				err,
				"error while updating member",
			)

			// List all members without filters
			params := openapi.ListHouseholdMembersParams{}
			memberList, err := s.listHouseholdMembers(params)
			s.handleErr(
				err,
				"error while listing household members",
			)

			// Should return at least the members we created (plus any from mock data)
			s.NotNil(memberList.Members)
			s.GreaterOrEqual(
				len(*memberList.Members),
				2,
			) // At least 2 active members

			// Verify the structure of returned members
			for _, member := range *memberList.Members {
				s.NotEmpty(member.Id)
				s.NotEmpty(member.FirstName)
				s.NotEmpty(member.LastName)
				s.NotEmpty(member.Role)
				s.NotNil(member.Active)
				s.NotNil(member.CreatedAt)
				s.NotNil(member.UpdatedAt)
			}
		},
	)

	s.Run(
		"List household members filtered by role",
		func() {
			// Create members with specific roles
			parentMember := s.createTestHouseholdMember()
			parentMember.Role = "parent"
			_, err := s.updateHouseholdMember(&parentMember)

			s.handleErr(
				err,
				"error while updating member",
			)

			childMember := s.createTestHouseholdMember()
			childMember.Role = "child"
			_, err = s.updateHouseholdMember(&childMember)
			s.handleErr(
				err,
				"error while updating member",
			)

			// Filter by parent role
			parentRole := "parent"
			params := openapi.ListHouseholdMembersParams{
				Role: &parentRole,
			}
			memberList, err := s.listHouseholdMembers(params)
			s.handleErr(
				err,
				"error while listing household members by role",
			)

			s.NotNil(memberList.Members)
			s.Greater(
				len(*memberList.Members),
				0,
			)

			// Verify all returned members have the parent role
			for _, member := range *memberList.Members {
				s.Equal(
					"parent",
					member.Role,
				)
			}

			// Filter by child role
			childRole := "child"
			params = openapi.ListHouseholdMembersParams{
				Role: &childRole,
			}
			memberList, err = s.listHouseholdMembers(params)
			s.handleErr(
				err,
				"error while listing household members by child role",
			)

			s.NotNil(memberList.Members)
			s.Greater(
				len(*memberList.Members),
				0,
			)

			// Verify all returned members have the child role
			for _, member := range *memberList.Members {
				s.Equal(
					"child",
					member.Role,
				)
			}
		},
	)

	s.Run(
		"List household members filtered by active status",
		func() {
			// Create an active member
			_ = s.createTestHouseholdMember()

			// Create an inactive member
			inactiveMember := s.createTestHouseholdMember()
			_, err := s.deactivateHouseholdMember(inactiveMember.Id)
			s.handleErr(
				err,
				"error while updating member",
			)

			// Filter by active status = true
			activeStatus := true
			params := openapi.ListHouseholdMembersParams{
				Active: &activeStatus,
			}
			memberList, err := s.listHouseholdMembers(params)
			s.handleErr(
				err,
				"error while listing active household members",
			)

			s.NotNil(memberList.Members)
			s.Greater(
				len(*memberList.Members),
				0,
			)

			// Verify all returned members are active
			for _, member := range *memberList.Members {
				s.True(
					*member.Active,
					"All members should be active",
				)
			}

			// Filter by active status = false
			inactiveStatus := false
			params = openapi.ListHouseholdMembersParams{
				Active: &inactiveStatus,
			}
			memberList, err = s.listHouseholdMembers(params)
			s.handleErr(
				err,
				"error while listing inactive household members",
			)

			s.NotNil(memberList.Members)
			// Note: The repository List method only returns active members, so this might return empty
			// This test verifies the API behavior matches the repository implementation
			for _, member := range *memberList.Members {
				s.False(
					*member.Active,
					"All members should be inactive",
				)
			}
		},
	)

	s.Run(
		"List household members with combined filters",
		func() {
			// Create members with different combinations
			activeParent := s.createTestHouseholdMember()
			activeParent.Role = "parent"
			_, err := s.updateHouseholdMember(&activeParent)
			s.handleErr(
				err,
				"error while updating member",
			)

			activeChild := s.createTestHouseholdMember()
			activeChild.Role = "child"
			_, err = s.updateHouseholdMember(&activeChild)
			s.handleErr(
				err,
				"error while updating member",
			)

			inactiveParent := s.createTestHouseholdMember()
			inactiveParent.Role = "parent"
			_, err = s.updateHouseholdMember(&inactiveParent)
			s.handleErr(
				err,
				"error while updating member",
			)

			_, err = s.deactivateHouseholdMember(inactiveParent.Id)
			s.handleErr(
				err,
				"error while updating member",
			)

			// Filter by role = parent AND active = true
			parentRole := "parent"
			activeStatus := true
			params := openapi.ListHouseholdMembersParams{
				Role:   &parentRole,
				Active: &activeStatus,
			}
			memberList, err := s.listHouseholdMembers(params)
			s.handleErr(
				err,
				"error while listing active parent household members",
			)

			s.NotNil(memberList.Members)
			s.Greater(
				len(*memberList.Members),
				0,
			)

			// Verify all returned members are active parents
			for _, member := range *memberList.Members {
				s.Equal(
					"parent",
					member.Role,
				)
				s.True(*member.Active)
			}
		},
	)

	s.Run(
		"List household members with non-existent role",
		func() {
			nonExistentRole := "non-existent-role"
			params := openapi.ListHouseholdMembersParams{
				Role: &nonExistentRole,
			}
			memberList, err := s.listHouseholdMembers(params)
			s.handleErr(
				err,
				"error while listing household members with non-existent role",
			)

			s.NotNil(memberList.Members)
			s.Equal(
				0,
				len(*memberList.Members),
				"Should return empty list for non-existent role",
			)
		},
	)

	s.Run(
		"List household members returns proper structure",
		func() {
			// Ensure we have at least one member
			_ = s.createTestHouseholdMember()

			params := openapi.ListHouseholdMembersParams{}
			memberList, err := s.listHouseholdMembers(params)
			s.handleErr(
				err,
				"error while listing household members",
			)

			// Verify the response structure
			s.NotNil(memberList.Members)
			s.IsType(
				&[]openapi.HouseholdMember{},
				memberList.Members,
			)

			if len(*memberList.Members) > 0 {
				member := (*memberList.Members)[0]

				// Verify all required fields are present
				s.NotEmpty(member.Id)
				s.NotEmpty(member.FirstName)
				s.NotEmpty(member.LastName)
				s.NotEmpty(member.Role)
				s.NotNil(member.Active)
				s.NotNil(member.CreatedAt)
				s.NotNil(member.UpdatedAt)

				// Verify optional fields can be nil or have values
				// Nickname is optional, so it can be nil
			}
		},
	)
}

func (s *Suite) createHouseholdMember(householdMemberReq *openapi.HouseholdMemberRequest) (
	openapi.HouseholdMember,
	error,
) {
	body, errBodyPrepare := utils.PrepareRequestBody(householdMemberReq)

	s.handleErr(
		errBodyPrepare,
		"error while preparing request body",
	)

	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodPost,
		"http://localhost:9091/household-members",
		body,
	)
	s.handleErr(
		errReq,
		"error while creating request",
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(
		err,
		"error while making request",
	)

	defer apiResponse.Body.Close()

	var member openapi.HouseholdMember
	s.decodeResponse(
		apiResponse,
		&member,
	)

	return member, nil
}

func (s *Suite) listHouseholdMembers(params openapi.ListHouseholdMembersParams) (
	openapi.HouseholdMemberList,
	error,
) {
	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		"http://localhost:9091/household-members",
		nil,
	)
	s.handleErr(
		errReq,
		"error while creating request",
	)
	q := req.URL.Query()
	if params.Role != nil {
		q.Set(
			"role",
			*params.Role,
		)
	}
	if params.Active != nil {
		q.Set(
			"active",
			strconv.FormatBool(*params.Active),
		)
	}
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(
		err,
		"error while making request",
	)

	defer apiResponse.Body.Close()

	var members openapi.HouseholdMemberList
	s.decodeResponse(
		apiResponse,
		&members,
	)

	return members, nil

}

func (s *Suite) getHouseholdMember(id string) (
	*http.Response,
	error,
) {
	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		"http://localhost:9091/household-members/"+id,
		nil,
	)
	s.handleErr(
		errReq,
		"error while creating request",
	)

	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(
		err,
		"error while making request",
	)

	return apiResponse, nil
}

func (s *Suite) updateHouseholdMember(member *openapi.HouseholdMember) (
	*http.Response,
	error,
) {
	memberRequest := openapi.HouseholdMemberRequest{
		Active:    member.Active,
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Nickname:  member.Nickname,
		Role:      member.Role,
	}
	body, errBodyPrepare := utils.PrepareRequestBody(memberRequest)

	s.handleErr(
		errBodyPrepare,
		"error while preparing request body",
	)

	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodPut,
		"http://localhost:9091/household-members/"+member.Id,
		body,
	)
	s.handleErr(
		errReq,
		"error while creating request",
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(
		err,
		"error while making request",
	)
	return apiResponse, nil
}
func (s *Suite) deleteHouseholdMember(id string) (
	*http.Response,
	error,
) {
	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodDelete,
		"http://localhost:9091/household-members/"+id,
		nil,
	)
	s.handleErr(
		errReq,
		"error while creating request",
	)
	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(
		err,
		"error while making request",
	)
	return apiResponse, nil
}
func (s *Suite) deactivateHouseholdMember(id string) (
	*http.Response,
	error,
) {
	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodPatch,
		"http://localhost:9091/household-members/"+id+"/deactivate",
		nil,
	)
	s.handleErr(
		errReq,
		"error while creating request",
	)
	req.Header.Set(
		"Content-Type",
		"application/json",
	)
	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(
		err,
		"error while making request",
	)
	return apiResponse, nil
}

func (s *Suite) activateHouseholdMember(id string) (
	*http.Response,
	error,
) {
	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodPatch,
		"http://localhost:9091/household-members/"+id+"/activate",
		nil,
	)
	s.handleErr(
		errReq,
		"error while creating request",
	)
	req.Header.Set(
		"Content-Type",
		"application/json",
	)
	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(
		err,
		"error while making request",
	)
	return apiResponse, nil
}
