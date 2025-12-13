package service

import (
	"context"

	"github.com/crutchm/elite/internal/models"
	"github.com/crutchm/elite/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetOrCreateUser(ctx context.Context, tgID int64) (*models.User, error) {
	return s.userRepo.GetOrCreateUser(ctx, tgID)
}

