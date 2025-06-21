package usecase

import (
    "context"
    "errors"
    "ghorkov32/proletariat-budget-be/internal/core/domain"
    "ghorkov32/proletariat-budget-be/internal/core/port"
    "golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
    authRepo port.AuthRepo
}

func NewAuthUseCase(authRepo port.AuthRepo) *AuthUseCase {
    return &AuthUseCase{
        authRepo: authRepo,
    }
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*domain.AuthToken, *domain.User, error) {
    user, err := uc.authRepo.GetUserByEmail(ctx, email)
    if err != nil {
        return nil, nil, errors.New("invalid credentials")
    }
    
    // Compare password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, nil, errors.New("invalid credentials")
    }
    
    // Generate token
    token, err := uc.authRepo.CreateToken(ctx, user.ID)
    if err != nil {
        return nil, nil, err
    }
    
    return token, user, nil
}

func (uc *AuthUseCase) ValidateToken(ctx context.Context, tokenStr string) (*domain.User, error) {
    return uc.authRepo.ValidateToken(ctx, tokenStr)
}

func (uc *AuthUseCase) RefreshToken(ctx context.Context, userID string) (*domain.AuthToken, *domain.User, error) {
    user, err := uc.authRepo.GetUserByID(ctx, userID)
    if err != nil {
        return nil, nil, err
    }
    
    token, err := uc.authRepo.CreateToken(ctx, userID)
    if err != nil {
        return nil, nil, err
    }
    
    return token, user, nil
}