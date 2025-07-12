package resthttp

import (
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
)

func FromOAPIAccount(a *openapi.Account) *domain.Account {
	return &domain.Account{
		ID:                 &a.Id,
		Name:               a.Name,
		Type:               domain.AccountType(a.Type),
		Currency:           a.Currency,
		InitialBalance:     a.InitialBalance,
		CurrentBalance:     a.CurrentBalance,
		Description:        a.Description,
		Institution:        a.Institution,
		AccountNumber:      a.AccountNumber,
		AccountInformation: a.AccountInformation,
		OwnerID:            &a.Owner.Id,
		Active:             a.Active,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}
}

// FromOAPIAccountRequest converts an OpenAPI AccountRequest to domain Account
func FromOAPIAccountRequest(a *openapi.AccountRequest) *domain.Account {
	var ownerID *string
	if a.Owner != nil {
		ownerID = &a.Owner.Id
	}

	return &domain.Account{
		Name:               a.Name,
		Type:               domain.AccountType(a.Type),
		Currency:           a.Currency,
		InitialBalance:     a.InitialBalance,
		CurrentBalance:     a.InitialBalance, // Set current balance to initial balance for new accounts
		Description:        a.Description,
		Institution:        a.Institution,
		AccountNumber:      a.AccountNumber,
		AccountInformation: a.AccountInformation,
		OwnerID:            ownerID,
		Active:             a.Active,
	}
}

// ToOAPIAccount converts a domain Account to OpenAPI Account
func ToOAPIAccount(account domain.Account) *openapi.Account {
	var id string
	if account.ID != nil {
		id = *account.ID
	}

	return &openapi.Account{
		Id:                 id,
		Name:               account.Name,
		Type:               openapi.AccountType(account.Type),
		Currency:           account.Currency,
		InitialBalance:     account.InitialBalance,
		CurrentBalance:     account.CurrentBalance,
		Description:        account.Description,
		Institution:        account.Institution,
		AccountNumber:      account.AccountNumber,
		AccountInformation: account.AccountInformation,
		Owner:              ToOAPIHouseholdMember(account.Owner),
		Active:             account.Active,
		CreatedAt:          account.CreatedAt,
		UpdatedAt:          account.UpdatedAt,
	}
}

// ToOAPIAccountList converts domain AccountList to OpenAPI AccountList
func ToOAPIAccountList(al *domain.AccountList) *openapi.AccountList {
	oapiAccounts := make([]openapi.Account, 0, len(al.Accounts))
	for _, account := range al.Accounts {
		oapiAccounts = append(oapiAccounts, *ToOAPIAccount(account))
	}
	return &openapi.AccountList{
		Accounts: &oapiAccounts,
		Metadata: &openapi.ListMetadata{
			Total:  al.Metadata.Total,
			Limit:  al.Metadata.Limit,
			Offset: al.Metadata.Offset,
		},
	}
}

func FromOAPIAccountListParams(params *openapi.ListAccountsParams) *domain.AccountListParams {
	return &domain.AccountListParams{
		Type:     params.Type,
		Currency: params.Currency,
		Active:   params.Active,
		Limit:    params.Limit,
		Offset:   params.Offset,
	}
}

func ToOAPIHouseholdMember(member *domain.HouseholdMember) *openapi.HouseholdMember {
	return &openapi.HouseholdMember{
		Id:        member.ID,
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Nickname:  member.Nickname,
		Role:      member.Role,
		Active:    &member.Active,
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}
}

func FromOAPIHouseholdMember(member *openapi.HouseholdMember) *domain.HouseholdMember {
	return &domain.HouseholdMember{
		ID:        member.Id,
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Nickname:  member.Nickname,
		Role:      member.Role,
		Active:    *member.Active,
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}
}

func FromOAPIHouseholdMemberRequest(member *openapi.HouseholdMemberRequest) *domain.HouseholdMember {
	return &domain.HouseholdMember{
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Nickname:  member.Nickname,
		Role:      member.Role,
		Active:    *member.Active,
	}
}

func ToOAPIHouseholdMemberList(members *domain.HouseholdMemberList) *openapi.HouseholdMemberList {
	oapiMembers := make([]openapi.HouseholdMember, 0, len(members.HouseholdMembers))
	for _, member := range members.HouseholdMembers {
		oapiMembers = append(oapiMembers, *ToOAPIHouseholdMember(&member))
	}
	return &openapi.HouseholdMemberList{
		Members: &oapiMembers,
	}
}
