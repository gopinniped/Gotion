package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gopinniped/gotion/internal/domain/entity"
	"github.com/gopinniped/gotion/internal/domain/errs"
	"github.com/gopinniped/gotion/internal/transport/http/dto"
	"github.com/gopinniped/gotion/internal/transport/http/mapper"
)

type UserService interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	LoginUser(ctx context.Context, username, password string) (string, error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type UserHandler struct {
	svc UserService
}

func NewUserHandler(svc UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapper.SendError(w, errs.New(errs.TypeValidation, "invalid json"))
		return
	}

	if err := req.Validate(); err != nil {
		mapper.SendError(w, err)
		return
	}

	createdUser, err := h.svc.CreateUser(r.Context(), req.ToEntity())
	if err != nil {
		mapper.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToUserResponse(createdUser))
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapper.SendError(w, errs.New(errs.TypeValidation, "invalid json"))
		return
	}

	if err := req.Validate(); err != nil {
		mapper.SendError(w, err)
		return
	}

	tokenString, err := h.svc.LoginUser(r.Context(), req.Username, req.Password)
	if err != nil {
		mapper.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.LoginResponse{Token: tokenString})
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		mapper.SendError(w, errs.New(errs.TypeValidation, "invalid user id"))
		return
	}

	user, err := h.svc.GetUserByID(r.Context(), id)
	if err != nil {
		mapper.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToUserResponse(user))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		mapper.SendError(w, errs.New(errs.TypeValidation, "invalid user id"))
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapper.SendError(w, errs.New(errs.TypeValidation, "invalid json"))
		return
	}

	if err := req.Validate(); err != nil {
		mapper.SendError(w, err)
		return
	}

	userEntity := req.ToEntity()
	userEntity.ID = id

	updatedUser, err := h.svc.UpdateUser(r.Context(), userEntity)
	if err != nil {
		mapper.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToUserResponse(updatedUser))
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		mapper.SendError(w, errs.New(errs.TypeValidation, "invalid user id"))
		return
	}

	if err := h.svc.DeleteUser(r.Context(), id); err != nil {
		mapper.SendError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
