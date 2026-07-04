package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gopinniped/gotion/internal/modules/tasks/entity"
	"github.com/gopinniped/gotion/internal/modules/tasks/errs"
)

type TaskRepositoryI interface {
	Create(ctx context.Context, userID uuid.UUID, task *entity.Task) (*entity.Task, error)
	Update(ctx context.Context, task *entity.Task) (*entity.Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
	IsOwner(ctx context.Context, userID, taskID uuid.UUID) (bool, error)
}

type TaskUseCase struct {
	taskRepo TaskRepositoryI
}

func NewTaskUseCase(taskRepo TaskRepositoryI) *TaskUseCase {
	return &TaskUseCase{taskRepo: taskRepo}
}

func (uc *TaskUseCase) CreateTask(ctx context.Context, userID uuid.UUID, task *entity.Task) (*entity.Task, error) {
	created, err := uc.taskRepo.Create(ctx, userID, task)
	if err != nil {
		return nil, fmt.Errorf("create task in repo: %w", err)
	}

	return created, nil
}

func (uc *TaskUseCase) UpdateTask(ctx context.Context, userID uuid.UUID, task *entity.Task) (*entity.Task, error) {
	if err := uc.checkOwnership(ctx, userID, task.ID); err != nil {
		return nil, err
	}

	updated, err := uc.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("update task in repo: %w", err)
	}

	return updated, nil
}

func (uc *TaskUseCase) GetTaskByID(ctx context.Context, userID, taskID uuid.UUID) (*entity.Task, error) {
	if err := uc.checkOwnership(ctx, userID, taskID); err != nil {
		return nil, err
	}

	task, err := uc.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("get task by id in repo: %w", err)
	}

	return task, nil
}

func (uc *TaskUseCase) GetTasksByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error) {
	tasks, err := uc.taskRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get tasks by user id in repo: %w", err)
	}

	return tasks, nil
}

func (uc *TaskUseCase) DeleteTask(ctx context.Context, userID, taskID uuid.UUID) error {
	if err := uc.checkOwnership(ctx, userID, taskID); err != nil {
		return err
	}

	if err := uc.taskRepo.Delete(ctx, taskID); err != nil {
		return fmt.Errorf("delete task in repo: %w", err)
	}

	return nil
}

func (uc *TaskUseCase) checkOwnership(ctx context.Context, userID, taskID uuid.UUID) error {
	owns, err := uc.taskRepo.IsOwner(ctx, userID, taskID)
	if err != nil {
		return fmt.Errorf("check task ownership: %w", err)
	}
	if !owns {
		return errs.ErrTaskForbidden
	}
	return nil
}
