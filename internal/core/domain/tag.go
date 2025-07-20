package domain

import "errors"

// Tag domain errors
var (
	ErrTagNotFound      = errors.New("tag not found")
	ErrTagInUse         = errors.New("tag is in use and cannot be deleted")
	ErrUnknownTagType   = errors.New("unknown tag type")
	ErrTagAlreadyExists = errors.New("tag with this name and type already exists")
	ErrTagNameEmpty     = errors.New("tag name cannot be empty")
)

type Tag struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     *string `json:"description"`
	Color           *string `json:"color"`
	BackgroundColor *string `json:"background_color"`
	TagType         TagType `json:"type"`
}

// TagType represents the typeof a tag (e.g., expenditure, Transfer, etc.)
type TagType string

const (
	TagTypeExpenditure         TagType = "expenditure"
	TagTypeTransfer            TagType = "transfer"
	TagTypeIngress             TagType = "ingress"
	TagTypeSavingsContribution TagType = "savings_contribution"
	TagTypeSavingsWithdrawal   TagType = "savings_withdrawal"
	TagTypeSavingsGoal         TagType = "saving_goal"
)

func (t Tag) Validate() error {
	switch t.TagType {
	case TagTypeExpenditure, TagTypeTransfer, TagTypeIngress, TagTypeSavingsContribution, TagTypeSavingsWithdrawal, TagTypeSavingsGoal:
		break
	default:
		return ErrUnknownTagType
	}
	if len(t.Name) == 0 {
		return ErrTagNameEmpty
	}
	return nil
}
