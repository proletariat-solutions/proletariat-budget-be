package domain

import "time"

type ExpenditureList struct {
	Expenditures []Expenditure `json:"expenditures"`
	Metadata     ListMetadata  `json:"metadata"`
}

type ExpenditureListParams struct {
	CategoryID  *string    `json:"categoryid"`
	StartDate   *time.Time `json:"date_from"`
	EndDate     *time.Time `json:"date_to"`
	Declared    *bool      `json:"declared"`
	Planned     *bool      `json:"planned"`
	Currency    *string    `json:"currency"`
	Description *string    `json:"description"`
	AccountID   *string    `json:"account_id"`
	Tags        *[]string  `json:"tags"`
	Limit       *int       `json:"limit"`
	Offset      *int       `json:"offset"`
}
