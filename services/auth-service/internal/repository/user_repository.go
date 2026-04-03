package repository

import (
	"context"
	"errors"

	"github.com/freedom-music/auth-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, `
        INSERT INTO users (email, password_hash)
        VALUES ($1, $2)
        RETURNING id, email, password_hash, created_at, updated_at
    `, email, passwordHash)

	return scanUser(row)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email=$1`, email)
	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	row := r.db.QueryRow(ctx, `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE id=$1`, id)
	return scanUser(row)
}

func scanUser(row pgx.Row) (*domain.User, error) {
	var u domain.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
