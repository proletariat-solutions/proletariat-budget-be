package domain

// AccountList represents a paginated list of accounts
type AccountList struct {
	Accounts []Account    `json:"accounts"`
	Metadata ListMetadata `json:"metadata"`
}

// IsEmpty returns true if the account list is empty
func (al *AccountList) IsEmpty() bool {
	return len(al.Accounts) == 0
}

// Count returns the number of accounts in the list
func (al *AccountList) Count() int {
	return len(al.Accounts)
}

type AccountListParams struct {
	Type     *string `form:"type,omitempty" json:"type,omitempty"`
	Currency *string `form:"currency,omitempty" json:"currency,omitempty"`
	Active   *bool   `form:"active,omitempty" json:"active,omitempty"`
	Limit    *int    `form:"limit,omitempty" json:"limit,omitempty"`
	Offset   *int    `form:"offset,omitempty" json:"offset,omitempty"`
}
