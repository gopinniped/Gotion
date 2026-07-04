package dto

import (
	"github.com/gopinniped/gotion/internal/modules/users/entity"
	"github.com/gopinniped/gotion/internal/shared/errs"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *CreateUserRequest) Validate() error {
	if r.Username == "" {
		return errs.New(errs.TypeValidation, "username is required")
	}
	if r.Password == "" {
		return errs.New(errs.TypeValidation, "password is required")
	}
	return nil
}

func (r *CreateUserRequest) ToEntity() *entity.User {
	return &entity.User{
		Username: r.Username,
		Password: r.Password,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if r.Username == "" {
		return errs.New(errs.TypeValidation, "username is required")
	}
	if r.Password == "" {
		return errs.New(errs.TypeValidation, "password is required")
	}
	return nil
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *UpdateUserRequest) Validate() error {
	if r.Username == "" && r.Password == "" {
		return errs.New(errs.TypeValidation, "at least one field must be provided")
	}
	return nil
}

func (r *UpdateUserRequest) ToEntity() *entity.User {
	return &entity.User{
		Username: r.Username,
		Password: r.Password,
	}
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func ToUserResponse(user *entity.User) UserResponse {
	if user == nil {
		return UserResponse{}
	}

	return UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
	}
}
