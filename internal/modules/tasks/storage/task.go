package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/gopinniped/gotion/internal/modules/tasks/entity"
	"github.com/gopinniped/gotion/internal/modules/tasks/errs"
	"github.com/jmoiron/sqlx"
)

type TaskStorage struct {
	DB *sqlx.DB
}

func NewTaskStorage(db *sqlx.DB) *TaskStorage {
	return &TaskStorage{DB: db}
}

func (s *TaskStorage) Create(ctx context.Context, userID uuid.UUID, task *entity.Task) (*entity.Task, error) {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	taskQuery := `
		INSERT INTO tasks (name, description, is_completed, due_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, is_completed, due_date, created_at, updated_at`

	var created entity.Task
	err = tx.GetContext(ctx, &created, taskQuery, task.Name, task.Description, task.IsCompleted, task.DueDate)
	if err != nil {
		return nil, err
	}

	linkQuery := `INSERT INTO user_tasks (user_id, task_id) VALUES ($1, $2)`
	_, err = tx.ExecContext(ctx, linkQuery, userID, created.ID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &created, nil
}

func (s *TaskStorage) Update(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	query := `
		UPDATE tasks
		SET
			name = COALESCE(NULLIF($1, ''), name),
			description = COALESCE(NULLIF($2, ''), description),
			is_completed = $3,
			due_date = COALESCE($4, due_date),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING id, name, description, is_completed, due_date, created_at, updated_at`

	var updated entity.Task
	err := s.DB.GetContext(ctx, &updated, query, task.Name, task.Description, task.IsCompleted, task.DueDate, task.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrTaskNotFound
		}
		return nil, err
	}

	return &updated, nil
}

func (s *TaskStorage) GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	query := `SELECT id, name, description, is_completed, due_date, created_at, updated_at FROM tasks WHERE id = $1`

	var task entity.Task
	err := s.DB.GetContext(ctx, &task, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrTaskNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (s *TaskStorage) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error) {
	query := `
		SELECT t.id, t.name, t.description, t.is_completed, t.due_date, t.created_at, t.updated_at
		FROM tasks t
		INNER JOIN user_tasks ut ON ut.task_id = t.id
		WHERE ut.user_id = $1
		ORDER BY t.created_at DESC`

	var tasks []*entity.Task
	err := s.DB.SelectContext(ctx, &tasks, query, userID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskStorage) IsOwner(ctx context.Context, userID, taskID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_tasks WHERE user_id = $1 AND task_id = $2)`

	var exists bool
	err := s.DB.GetContext(ctx, &exists, query, userID, taskID)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *TaskStorage) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errs.ErrTaskNotFound
	}

	return nil
}
