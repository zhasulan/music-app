package service

import (
	"context"
	"errors"
	"time"

	"github.com/freedom-music/auth-service/internal/domain"
	"github.com/freedom-music/auth-service/internal/repository"
	sharedAuth "github.com/freedom-music/shared/auth"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailExists  = errors.New("email already registered")
	ErrInvalidLogin = errors.New("invalid email or password")
	ErrInvalidToken = errors.New("invalid token")
)

type AuthService struct {
	repo       *repository.UserRepository
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(repo *repository.UserRepository, jwtSecret string, accessTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*domain.User, string, string, error) {
	existing, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", "", err
	}
	if existing != nil {
		return nil, "", "", ErrEmailExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", err
	}

	user, err := s.repo.Create(ctx, email, string(hash))
	if err != nil {
		return nil, "", "", err
	}

	access, refresh, err := s.issueTokens(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, access, refresh, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.User, string, string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", "", err
	}
	if user == nil {
		return nil, "", "", ErrInvalidLogin
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", "", ErrInvalidLogin
	}

	access, refresh, err := s.issueTokens(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, access, refresh, nil
}

func (s *AuthService) Refresh(token string) (string, string, error) {
	claims, err := sharedAuth.ParseToken(token, s.jwtSecret)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	access, refresh, err := s.issueWithSubject(claims.Subject)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func (s *AuthService) issueTokens(user *domain.User) (string, string, error) {
	return s.issueWithSubject(user.IDString())
}

func (s *AuthService) issueWithSubject(sub string) (string, string, error) {
	access, err := sharedAuth.GenerateToken(sub, s.jwtSecret, s.accessTTL)
	if err != nil {
		return "", "", err
	}
	refresh, err := sharedAuth.GenerateToken(sub, s.jwtSecret, s.refreshTTL)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}
