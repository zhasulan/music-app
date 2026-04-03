package service

import (
	"context"

	"github.com/freedom-music/user-service/internal/domain"
	"github.com/freedom-music/user-service/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Get(ctx context.Context, id int64) (*domain.UserProfile, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	placeholder := &domain.UserProfile{ID: id, Email: "", Name: "", Country: ""}
	return s.repo.Create(ctx, placeholder)
}

func (s *UserService) Update(ctx context.Context, id int64, name, country string) (*domain.UserProfile, error) {
	if _, err := s.Get(ctx, id); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, id, name, country)
}
