package domain

import (
	"errors"
	"time"
)

type Expenditure struct {
	ID          string       `json:"id"`
	Category    *Category    `json:"category"`
	Declared    bool         `json:"declared"`
	Planned     bool         `json:"planned"`
	Transaction *Transaction `json:"transaction,omitempty"`
	// Making a pointer to each tag since there can be a lot if expenditures are listed
	Tags *[]*Tag   `json:"tags,omitempty"`
	Date time.Time `json:"date"`
}

var (
	ErrExpenditureNotFound = errors.New("expenditure not found")
)
