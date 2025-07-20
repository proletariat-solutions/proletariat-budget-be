package domain

import "errors"

// Category domain errors
var (
	ErrCategoryNotFound              = errors.New("category not found")
	ErrCategoryHasActiveSavingsGoals = errors.New("category has active savings goals and cannot be deleted")
	ErrCategoryInactive              = errors.New("category is inactive")
	ErrCategoryAlreadyActive         = errors.New("category is already active")
	ErrCategoryAlreadyInactive       = errors.New("category is already inactive")
	ErrCategoryUsedInExpenditure     = errors.New("category is used in expenditures")
	ErrCategoryUsedInTransfer        = errors.New("category is used in transfers")
	ErrCategoryUsedInSavingGoal      = errors.New("category is used in saving goals")
	ErrCategoryUsedInIngress         = errors.New("category is used in ingresses")
	ErrCategoryUsedInEntity          = errors.New("category is used in entity")
)

type Category struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Description     *string      `json:"description"`
	Color           string       `json:"color"`
	BackgroundColor string       `json:"background_color"`
	Active          bool         `json:"active"`
	CategoryType    CategoryType `json:"type"`
}

type CategoryType string

const (
	CategoryTypeIngress     CategoryType = "ingress"
	CategoryTypeExpenditure CategoryType = "expenditure"
	CategoryTypeTransfer    CategoryType = "transfer"
	CategoryTypeSavingGoal  CategoryType = "saving_goal"
)

func (c *Category) Activate() error {
	if c.Active {
		return ErrCategoryAlreadyActive
	}
	c.Active = true
	return nil
}

func (c *Category) Deactivate() error {
	if !c.Active {
		return ErrCategoryAlreadyInactive
	}
	c.Active = false
	return nil
}
