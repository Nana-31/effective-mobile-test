package usecase

import (
	"context"
	"fmt"
	"test_em/internal/domain"
	"test_em/internal/repository"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, d domain.Subscription) error {
	if d.UserID == uuid.Nil {
		return fmt.Errorf("user ID cannot be nil")
	}

	if d.Price <= 0 {
		return fmt.Errorf("price cannot be <= 0")
	}

	if !d.EndDate.IsZero() && !d.EndDate.After(d.StartDate) {
		return fmt.Errorf("end date must be after start date")
	}

	err := s.repo.Create(ctx, d)
	if err != nil {
		return fmt.Errorf("error in repo create: %w", err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, d domain.Subscription) error {
	if d.ID == uuid.Nil {
		return fmt.Errorf("subscription ID cannot be nil")
	}
	if d.UserID == uuid.Nil {
		return fmt.Errorf("user ID cannot be nil")
	}

	if d.Price <= 0 {
		return fmt.Errorf("price cannot be <= 0")
	}

	if !d.EndDate.IsZero() && !d.EndDate.After(d.StartDate) {
		return fmt.Errorf("end date must be after start date")
	}

	err := s.repo.Update(ctx, d)
	if err != nil {
		return fmt.Errorf("error in repo update: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("id cannot be nil")
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error in repo delete: %w", err)
	}

	return nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (domain.Subscription, error) {
	if id == uuid.Nil {
		return domain.Subscription{}, fmt.Errorf("id cannot be nil")
	}

	sub, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("error in repo get: %w", err)
	}

	return sub, nil
}

func (s *Service) List(ctx context.Context, userId uuid.UUID) ([]domain.Subscription, error) {
	if userId == uuid.Nil {
		return nil, fmt.Errorf("user ID cannot be nil")
	}

	subs, err := s.repo.List(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("error in repo list: %w", err)
	}

	return subs, nil
}

func (s *Service) GetTotalCost(ctx context.Context, user_id uuid.UUID, start_date time.Time, end_date time.Time, service_name string) (int, error) {
	if user_id == uuid.Nil {
		return 0, fmt.Errorf("user ID cannot be nil")
	}

	if !end_date.IsZero() && !end_date.After(start_date) {
		return 0, fmt.Errorf("end date must be after start date")
	}

	res, err := s.repo.GetTotalCost(ctx, user_id, start_date, end_date, service_name)
	if err != nil {
		return 0, fmt.Errorf("error in repo get total cost: %w", err)
	}

	return res, nil
}
