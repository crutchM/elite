package repository

import (
	"context"

	"github.com/crutchm/elite/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) GetOrCreateUser(ctx context.Context, tgID int64) (*models.User, error) {
	user := &models.User{TGID: tgID}

	query := `
		INSERT INTO users (tg_id)
		VALUES ($1)
		ON CONFLICT (tg_id) DO NOTHING
		RETURNING tg_id
	`

	err := r.pool.QueryRow(ctx, query, tgID).Scan(&user.TGID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	if err == pgx.ErrNoRows {
		query = `SELECT tg_id FROM users WHERE tg_id = $1`
		err = r.pool.QueryRow(ctx, query, tgID).Scan(&user.TGID)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

