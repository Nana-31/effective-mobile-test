package repository

import (
	"context"
	"fmt"
	"test_em/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/google/uuid"
)

type PostgresRepository struct {
	db *pgx.Conn
}

func NewPostgresRepository(db *pgx.Conn) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (db *PostgresRepository) Create(ctx context.Context, model domain.Subscription) error {
	query := "INSERT INTO subscriptions(id, user_id, service_name, price, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := db.db.Exec(ctx, query, model.ID, model.UserID, model.ServiceName, model.Price, model.StartDate, model.EndDate)

	if err != nil {
		return fmt.Errorf("Error with create %w", err)
	}

	return nil
}

func (db *PostgresRepository) Update(ctx context.Context, model domain.Subscription) error {
	query := "UPDATE subscriptions SET service_name = $1, price = $2, start_date = $3, end_date = $4 WHERE id = $5"

	_, err := db.db.Exec(ctx, query, model.ServiceName, model.Price, model.StartDate, model.EndDate, model.ID)

	if err != nil {
		return fmt.Errorf("Error with update %w", err)
	}

	return nil
}

func (db *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM subscriptions WHERE id = $1"

	_, err := db.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("Error with delete %w", err)
	}

	return nil
}

func (db *PostgresRepository) Get(ctx context.Context, id uuid.UUID) (domain.Subscription, error) {
	query := "SELECT id, user_id, service_name, price, start_date, end_date FROM subscriptions WHERE id = $1"

	var sub domain.Subscription
	row := db.db.QueryRow(ctx, query, id)
	err := row.Scan(&sub.ID, &sub.UserID, &sub.ServiceName, &sub.Price, &sub.StartDate, &sub.EndDate)

	if err != nil {
		return domain.Subscription{}, fmt.Errorf("Error: %w", err)
	}

	return sub, err
}

func (db *PostgresRepository) List(ctx context.Context, UserId uuid.UUID) ([]domain.Subscription, error) {
	query := "SELECT id, user_id, service_name, price, start_date, end_date FROM subscriptions WHERE user_id = $1"
	subs, err := db.db.Query(ctx, query, UserId)
	if err != nil {
		return nil, fmt.Errorf("Error: %w", err)
	}
	defer subs.Close()

	var result []domain.Subscription
	for subs.Next() {
		var sub domain.Subscription

		err := subs.Scan(&sub.ID, &sub.UserID, &sub.ServiceName, &sub.Price, &sub.StartDate, &sub.EndDate)
		if err != nil {
			return nil, fmt.Errorf("Error with scan: %w", err)
		}
		result = append(result, sub)
	}

	if err = subs.Err(); err != nil {
		return nil, fmt.Errorf("Error with rows: %w", err)
	}

	return result, err
}

func (db *PostgresRepository) GetTotalCost(ctx context.Context, user_id uuid.UUID, start_date time.Time, end_date time.Time, service_name string) (int, error) {
	query := "SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE user_id = $1 AND start_date BETWEEN $2 AND $3"
	args := []any{user_id, start_date, end_date}
	if service_name != "" {
		query += " AND service_name = $4"
		args = append(args, service_name)
	}

	var sum int
	err := db.db.QueryRow(ctx, query, args...).Scan(&sum)
	if err != nil {
		return 0, err
	}

	return sum, nil
}
