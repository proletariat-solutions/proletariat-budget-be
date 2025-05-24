package domain

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrUnableToFetchScopes = errors.New("unable to fetch scopes from lookup api")
	ErrForbiddenAction     = errors.New("user is not allowed to execute action")
	ErrForbiddenResource   = errors.New("user is not allowed to access resource")
)
