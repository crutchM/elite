package repository

import (
	"context"
	"errors"

	"github.com/crutchm/elite/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VoteRepository struct {
	pool *pgxpool.Pool
}

func NewVoteRepository(pool *pgxpool.Pool) *VoteRepository {
	return &VoteRepository{pool: pool}
}

func (r *VoteRepository) CreateVote(ctx context.Context, vote *models.Vote) error {
	query := `
		INSERT INTO votes (tg_user_id, nominant_id, category_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (tg_user_id, category_id) DO NOTHING
		RETURNING tg_user_id, nominant_id, category_id
	`

	err := r.pool.QueryRow(ctx, query, vote.TGUserID, vote.NominantID, vote.CategoryID).
		Scan(&vote.TGUserID, &vote.NominantID, &vote.CategoryID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("vote already exists for this category")
		}
		return err
	}

	return nil
}

func (r *VoteRepository) GetNominantCategory(ctx context.Context, nominantID int64) (int64, error) {
	var categoryID int64
	query := `SELECT category_id FROM nominants WHERE id = $1`

	err := r.pool.QueryRow(ctx, query, nominantID).Scan(&categoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New("nominant not found")
		}
		return 0, err
	}

	return categoryID, nil
}

func (r *VoteRepository) CheckNominantExists(ctx context.Context, nominantID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM nominants WHERE id = $1)`

	err := r.pool.QueryRow(ctx, query, nominantID).Scan(&exists)
	return exists, err
}

