package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"body-smith-be/internal/model"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, email, hashedPassword string) (*model.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, `
		SELECT id, email, password, created_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, email, hashedPassword string) (*model.User, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO users (email, password)
		VALUES (?, ?)
	`, email, hashedPassword)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := r.db.GetContext(ctx, &user, `
		SELECT id, email, password, created_at
		FROM users
		WHERE id = ?
	`, id); err != nil {
		return nil, err
	}

	return &user, nil
}
