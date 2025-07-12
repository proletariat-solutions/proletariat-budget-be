package domain

import "ghorkov32/proletariat-budget-be/openapi"

// HouseholdMemberList represents a paginated list of household members
type HouseholdMemberList struct {
	HouseholdMembers []HouseholdMember `json:"household_members"`
}

// ToOAPIHouseholdMemberList converts domain HouseholdMemberList to OpenAPI HouseholdMemberList
func (hml *HouseholdMemberList) ToOAPIHouseholdMemberList() *openapi.HouseholdMemberList {
	oapiMembers := make([]openapi.HouseholdMember, len(hml.HouseholdMembers))

	for i, member := range hml.HouseholdMembers {
		oapiMembers[i] = *member.ToOAPIHouseholdMember()
	}

	return &openapi.HouseholdMemberList{
		Members: &oapiMembers,
	}
}

// IsEmpty returns true if the household member list is empty
func (hml *HouseholdMemberList) IsEmpty() bool {
	return len(hml.HouseholdMembers) == 0
}

// Count returns the number of household members in the list
func (hml *HouseholdMemberList) Count() int {
	return len(hml.HouseholdMembers)
}

// GetActiveMembers returns only active household members
func (hml *HouseholdMemberList) GetActiveMembers() []HouseholdMember {
	var activeMembers []HouseholdMember
	for _, member := range hml.HouseholdMembers {
		if member.IsActive() {
			activeMembers = append(activeMembers, member)
		}
	}
	return activeMembers
}

// GetMemberByID returns a household member by ID, or nil if not found
func (hml *HouseholdMemberList) GetMemberByID(id string) *HouseholdMember {
	for _, member := range hml.HouseholdMembers {
		if member.ID == id {
			return &member
		}
	}
	return nil
}

// GetMembersByRole returns household members with the specified role
func (hml *HouseholdMemberList) GetMembersByRole(role string) []HouseholdMember {
	var members []HouseholdMember
	for _, member := range hml.HouseholdMembers {
		if member.Role == role {
			members = append(members, member)
		}
	}
	return members
}
