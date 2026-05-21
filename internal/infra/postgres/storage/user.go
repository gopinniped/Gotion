package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/gopinniped/gotion/internal/domain/entity"
	"github.com/gopinniped/gotion/internal/domain/errs"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserStorage struct {
	DB *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{DB: db}
}

func (s *UserStorage) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`

	err := s.DB.QueryRowContext(ctx, query, user.Username, user.Password).Scan(&user.ID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return nil, errs.ErrUserAlreadyExists
		}
		return nil, err
	}

	return user, nil
}

func (s *UserStorage) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `
        UPDATE users 
        SET 
            username = COALESCE(NULLIF($1, ''), username), 
            password = COALESCE(NULLIF($2, ''), password) 
        WHERE id = $3 
        RETURNING id, username, password`

	var updatedUser entity.User
	err := s.DB.GetContext(ctx, &updatedUser, query, user.Username, user.Password, user.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return nil, errs.ErrUserAlreadyExists
		}
		return nil, err
	}

	return &updatedUser, nil
}

func (s *UserStorage) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `SELECT id, username, password FROM users WHERE id = $1`

	var user entity.User
	err := s.DB.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserStorage) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `SELECT id, username, password FROM users WHERE username = $1`

	var user entity.User
	err := s.DB.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserStorage) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}
