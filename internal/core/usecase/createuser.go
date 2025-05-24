package usecase

import (
	"context"

	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"github.com/rs/zerolog/log"
)

type CreateUser struct {
	userRepo port.StoringUser
}

func NewCreateUser(userRepo port.StoringUser) *CreateUser {
	return &CreateUser{userRepo: userRepo}
}

func (s *CreateUser) Create(ctx context.Context, input domain.User) error {
	err := s.userRepo.Save(ctx, input)
	if err != nil {
		log.Error().Err(err).Msg("unable to create user")

		return err
	}

	return nil
}
