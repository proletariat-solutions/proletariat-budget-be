package mysql

import (
	"context"
	"database/sql"
	"errors"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AuthRepoImpl struct {
	db     *sql.DB
	secret []byte // JWT secret
}

func NewAuthRepository(db *sql.DB, secret string) port.AuthRepo {
	return &AuthRepoImpl{
		db:     db,
		secret: []byte(secret),
	}
}

func (r *AuthRepoImpl) CreateUser(ctx context.Context, user domain.User) (string, error) {
	return "", errors.New("not implemented") // Replace with actual implementation
}

func (r *AuthRepoImpl) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {

	return nil, errors.New("not implemented") // Replace with actual implementation
}

func (r *AuthRepoImpl) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return nil, errors.New("not implemented") // Replace with actual implementation
}

func (r *AuthRepoImpl) UpdateUser(ctx context.Context, id string, user domain.User) error {
	return errors.New("not implemented") // Replace with actual implementation
}

func (r *AuthRepoImpl) CreateToken(ctx context.Context, userID string) (*domain.AuthToken, error) {
	// Create JWT token
	expiresAt := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(r.secret)
	if err != nil {
		return nil, err
	}

	// Store token in database for revocation capability
	// ...

	return &domain.AuthToken{
		Token:     tokenString,
		ExpiresAt: expiresAt,
		UserID:    userID,
	}, nil
}

func (r *AuthRepoImpl) ValidateToken(ctx context.Context, tokenStr string) (*domain.User, error) {
	// Parse and validate JWT token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return r.secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check if token is revoked
	// ...

	// Extract user ID from claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	// Get user from database
	return r.GetUserByID(ctx, userID)
}

func (r *AuthRepoImpl) RevokeToken(ctx context.Context, token string) error {
	return errors.New("not implemented") // Replace with actual implementation
}
