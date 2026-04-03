package repository

import (
	"context"
	"errors"

	"github.com/freedom-music/user-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.UserProfile, error) {
	row := r.db.QueryRow(ctx, `SELECT id, email, name, country, created_at, updated_at FROM users WHERE id=$1`, id)
	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.UserProfile) (*domain.UserProfile, error) {
	row := r.db.QueryRow(ctx, `
        INSERT INTO users (id, email, name, country)
        VALUES ($1, $2, $3, $4)
        RETURNING id, email, name, country, created_at, updated_at
    `, user.ID, user.Email, user.Name, user.Country)
	return scanUser(row)
}

func (r *UserRepository) Update(ctx context.Context, id int64, name, country string) (*domain.UserProfile, error) {
	row := r.db.QueryRow(ctx, `
        UPDATE users SET name=$2, country=$3, updated_at=NOW()
        WHERE id=$1
        RETURNING id, email, name, country, created_at, updated_at
    `, id, name, country)
	return scanUser(row)
}

func scanUser(row pgx.Row) (*domain.UserProfile, error) {
	var u domain.UserProfile
	if err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Country, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
