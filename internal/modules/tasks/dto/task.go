package dto

import (
	"time"

	"github.com/gopinniped/gotion/internal/modules/tasks/entity"
	"github.com/gopinniped/gotion/internal/shared/errs"
)

type CreateTaskRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
}

func (r *CreateTaskRequest) Validate() error {
	if r.Name == "" {
		return errs.New(errs.TypeValidation, "task name is required")
	}
	return nil
}

func (r *CreateTaskRequest) ToEntity() *entity.Task {
	return &entity.Task{
		Name:        r.Name,
		Description: r.Description,
		DueDate:     r.DueDate,
	}
}

type UpdateTaskRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IsCompleted *bool      `json:"is_completed"`
	DueDate     *time.Time `json:"due_date"`
}

func (r *UpdateTaskRequest) Validate() error {
	if r.Name == "" && r.Description == "" && r.IsCompleted == nil && r.DueDate == nil {
		return errs.New(errs.TypeValidation, "at least one field must be provided")
	}
	return nil
}

func (r *UpdateTaskRequest) ToEntity() *entity.Task {
	task := &entity.Task{
		Name:        r.Name,
		Description: r.Description,
		DueDate:     r.DueDate,
	}
	if r.IsCompleted != nil {
		task.IsCompleted = *r.IsCompleted
	}
	return task
}

type TaskResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IsCompleted bool       `json:"is_completed"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func ToTaskResponse(task *entity.Task) TaskResponse {
	if task == nil {
		return TaskResponse{}
	}

	return TaskResponse{
		ID:          task.ID.String(),
		Name:        task.Name,
		Description: task.Description,
		IsCompleted: task.IsCompleted,
		DueDate:     task.DueDate,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func ToTaskListResponse(tasks []*entity.Task) []TaskResponse {
	result := make([]TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		result = append(result, ToTaskResponse(t))
	}
	return result
}
