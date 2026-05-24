package repository

import (
	"context"
	"test_em/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, d domain.Subscription) error
	Update(ctx context.Context, d domain.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (domain.Subscription, error)
	List(ctx context.Context, userId uuid.UUID) ([]domain.Subscription, error)
	GetTotalCost(ctx context.Context, user_id uuid.UUID, start_date time.Time, end_date time.Time, service_name string) (int, error)
}
