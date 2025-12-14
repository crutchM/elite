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

// CheckNominantExists проверяет существование номинанта
func (r *VoteRepository) CheckNominantExists(ctx context.Context, nominantID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM nominants WHERE id = $1)`

	err := r.pool.QueryRow(ctx, query, nominantID).Scan(&exists)
	return exists, err
}

// CheckNominantInCategory проверяет, участвует ли номинант в указанной категории
func (r *VoteRepository) CheckNominantInCategory(ctx context.Context, nominantID, categoryID int64) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM nominant_categories 
			WHERE nominant_id = $1 AND category_id = $2
		)
	`

	err := r.pool.QueryRow(ctx, query, nominantID, categoryID).Scan(&exists)
	return exists, err
}

// CheckCategoryExists проверяет существование категории
func (r *VoteRepository) CheckCategoryExists(ctx context.Context, categoryID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)`

	err := r.pool.QueryRow(ctx, query, categoryID).Scan(&exists)
	return exists, err
}
