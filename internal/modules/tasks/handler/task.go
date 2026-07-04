package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gopinniped/gotion/internal/modules/tasks/dto"
	"github.com/gopinniped/gotion/internal/modules/tasks/entity"
	"github.com/gopinniped/gotion/internal/shared/errs"
	"github.com/gopinniped/gotion/internal/shared/mapper"
	"github.com/gopinniped/gotion/internal/shared/middleware"
)

type TaskService interface {
	CreateTask(ctx context.Context, userID uuid.UUID, task *entity.Task) (*entity.Task, error)
	UpdateTask(ctx context.Context, userID uuid.UUID, task *entity.Task) (*entity.Task, error)
	GetTaskByID(ctx context.Context, userID, taskID uuid.UUID) (*entity.Task, error)
	GetTasksByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error)
	DeleteTask(ctx context.Context, userID, taskID uuid.UUID) error
}

type TaskHandler struct {
	svc TaskService
}

func NewTaskHandler(svc TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		mapper.SendError(w, errs.New(errs.TypeUnauthenticated, "user not authenticated"))
		return
	}

	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapper.SendError(w, errs.New(errs.TypeValidation, "invalid json"))
		return
	}

	if err := req.Validate(); err != nil {
		mapper.SendError(w, err)
		return
	}

	created, err := h.svc.CreateTask(r.Context(), userID, req.ToEntity())
	if err != nil {
		mapper.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToTaskResponse(created))
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		mapper.SendError(w, errs.New(errs.TypeUnauthenticated, "user not authenticated"))
		return
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		mapper.SendError(w, errs.New(errs.TypeValidation, "invalid task id"))
		return
	}

	task, err := h.svc.GetTaskByID(r.Context(), userID, taskID)
	if err != nil {
		mapper.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToTaskResponse(task))
}

func (h *TaskHandler) GetMyTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		mapper.SendError(w, errs.New(errs.TypeUnauthenticated, "user not authenticated"))
		return
	}

	tasks, err := h.svc.GetTasksByUserID(r.Context(), userID)
	if err != nil {
		mapper.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToTaskListResponse(tasks))
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		mapper.SendError(w, errs.New(errs.TypeUnauthenticated, "user not authenticated"))
		return
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		mapper.SendError(w, errs.New(errs.TypeValidation, "invalid task id"))
		return
	}

	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapper.SendError(w, errs.New(errs.TypeValidation, "invalid json"))
		return
	}

	if err := req.Validate(); err != nil {
		mapper.SendError(w, err)
		return
	}

	taskEntity := req.ToEntity()
	taskEntity.ID = taskID

	updated, err := h.svc.UpdateTask(r.Context(), userID, taskEntity)
	if err != nil {
		mapper.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToTaskResponse(updated))
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		mapper.SendError(w, errs.New(errs.TypeUnauthenticated, "user not authenticated"))
		return
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		mapper.SendError(w, errs.New(errs.TypeValidation, "invalid task id"))
		return
	}

	if err := h.svc.DeleteTask(r.Context(), userID, taskID); err != nil {
		mapper.SendError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
