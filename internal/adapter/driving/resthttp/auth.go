package resthttp

import (
	"context"

	"ghorkov32/proletariat-budget-be/openapi"
)

func (c *Controller) Login(
	ctx context.Context,
	request openapi.LoginRequestObject,
) (
	openapi.LoginResponseObject,
	error,
) {
	panic("implement me")
}

func (c *Controller) RefreshToken(
	ctx context.Context,
	request openapi.RefreshTokenRequestObject,
) (
	openapi.RefreshTokenResponseObject,
	error,
) {
	panic("implement me")
}

func (c *Controller) RegisterUser(
	ctx context.Context,
	request openapi.RegisterUserRequestObject,
) (
	openapi.RegisterUserResponseObject,
	error,
) {
	panic("implement me")
}
