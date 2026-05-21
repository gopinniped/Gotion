package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gopinniped/gotion/internal/domain/entity"
	"github.com/gopinniped/gotion/internal/domain/errs"
	"github.com/gopinniped/gotion/pkg/hash"
	"github.com/gopinniped/gotion/pkg/token"
)

type UserRepositoryI interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserUseCase struct {
	userRepo   UserRepositoryI
	tokenMaker *token.TokenMaker
}

func NewUserUseCase(userRepo UserRepositoryI, tokenMaker *token.TokenMaker) *UserUseCase {
	return &UserUseCase{userRepo: userRepo, tokenMaker: tokenMaker}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	hashedPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	newUser := entity.User{
		Username: user.Username,
		Password: hashedPassword,
	}

	createdUser, err := uc.userRepo.Create(ctx, &newUser)
	if err != nil {
		return nil, fmt.Errorf("create user in repo: %w", err)
	}

	return createdUser, nil
}

func (uc *UserUseCase) LoginUser(ctx context.Context, username, password string) (string, error) {
	user, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", errs.ErrInvalidCredentials
	}

	if !hash.CheckPassword(password, user.Password) {
		return "", errs.ErrInvalidCredentials
	}

	tokenString, err := uc.tokenMaker.Generate(user.ID, user.Username)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return tokenString, nil
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	if user.Password != "" {
		hashedPassword, err := hash.HashPassword(user.Password)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}
		user.Password = hashedPassword
	}

	updatedUser, err := uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("update user in repo: %w", err)
	}

	return updatedUser, nil
}

func (uc *UserUseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getbyid user in repo: %w", err)
	}

	return user, nil
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := uc.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete user in repo: %w", err)
	}

	return nil
}
