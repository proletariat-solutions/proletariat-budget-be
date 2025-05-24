package mongodb

import (
	"context"
	"errors"

	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	DuplicateErrorCode = 11000
)

type UserRepo struct {
	client *mongo.Client
}

func NewUserRepo(client *mongo.Client) *UserRepo {
	return &UserRepo{client: client}
}

func (r *UserRepo) Save(ctx context.Context, input domain.User) error {
	_, err := r.client.Database("test").Collection("users").InsertOne(ctx, input)

	switch {
	case isDuplicatedKey(err):
		return domain.ErrUserAlreadyExists
	case err != nil:
		return nil
	}

	return nil
}

func isDuplicatedKey(err error) bool {
	var e mongo.WriteException
	if !errors.As(err, &e) {
		return false
	}

	for _, we := range e.WriteErrors {
		if we.Code == DuplicateErrorCode {
			return true
		}
	}

	return false
}
