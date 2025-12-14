package service

import (
	"context"
	"errors"

	"github.com/crutchm/elite/internal/models"
	"github.com/crutchm/elite/internal/repository"
)

type VoteService struct {
	voteRepo *repository.VoteRepository
}

func NewVoteService(voteRepo *repository.VoteRepository) *VoteService {
	return &VoteService{voteRepo: voteRepo}
}

func (s *VoteService) CreateVote(ctx context.Context, tgUserID int64, request *models.VoteRequest) error {
	// Проверяем существование номинанта
	exists, err := s.voteRepo.CheckNominantExists(ctx, request.NominantID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("nominant not found")
	}

	// Проверяем существование категории
	categoryExists, err := s.voteRepo.CheckCategoryExists(ctx, request.CategoryID)
	if err != nil {
		return err
	}
	if !categoryExists {
		return errors.New("category not found")
	}

	// Проверяем, что номинант участвует в указанной категории
	inCategory, err := s.voteRepo.CheckNominantInCategory(ctx, request.NominantID, request.CategoryID)
	if err != nil {
		return err
	}
	if !inCategory {
		return errors.New("nominant is not participating in this category")
	}

	vote := &models.Vote{
		TGUserID:   tgUserID,
		NominantID: request.NominantID,
		CategoryID: request.CategoryID,
	}

	return s.voteRepo.CreateVote(ctx, vote)
}
