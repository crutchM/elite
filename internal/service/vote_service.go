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
	exists, err := s.voteRepo.CheckNominantExists(ctx, request.NominantID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("nominant not found")
	}

	categoryID, err := s.voteRepo.GetNominantCategory(ctx, request.NominantID)
	if err != nil {
		return err
	}

	vote := &models.Vote{
		TGUserID:   tgUserID,
		NominantID: request.NominantID,
		CategoryID: categoryID,
	}

	return s.voteRepo.CreateVote(ctx, vote)
}

