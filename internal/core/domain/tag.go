package domain

import "errors"

type Tag struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     *string `json:"description"`
	Color           string  `json:"color"`
	BackgroundColor string  `json:"background_color"`
	TagType         string  `json:"type"`
}

var (
	ErrTagNotFound = errors.New("tag not found")
)
