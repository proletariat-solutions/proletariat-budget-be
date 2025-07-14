package domain

import "errors"

type Category struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     *string `json:"description"`
	Color           string  `json:"color"`
	BackgroundColor string  `json:"background_color"`
	Active          bool    `json:"active"`
}

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryInactive = errors.New("category is inactive")
)
