package repository

import (
	"awesomeProject1/internal/dto"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type SubscriptionRepository struct {
	db *pgxpool.Pool
}

type ISubscriptionRepository interface {
	FindAll(ctx context.Context, offset, limit int) ([]*dto.Subscription, int, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	Create(ctx context.Context, category *dto.Subscription) (*dto.Subscription, error)
	Update(ctx context.Context, queryParts string, values []any) (bool, error)
	FindById(ctx context.Context, id uuid.UUID) (*dto.Subscription, bool, error)
	GetTotal(ctx context.Context, start, end time.Time, serviceName, userId string) (int, error)
}

func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}
func (c *SubscriptionRepository) GetTotal(ctx context.Context, start, end time.Time, serviceName, userId string) (int, error) {
	query := `
        SELECT COALESCE(SUM(price), 0)
        FROM subscriptions
        WHERE start_date >= $1
          AND (start_date <= $2)
          AND ($3::uuid IS NULL OR user_id = $3)
          AND ($4::text IS NULL OR service_name = $4);
    `

	var total int
	err := c.db.QueryRow(ctx, query, start, end, nullString(userId), nullString(serviceName)).Scan(&total)
	return total, err
}

func nullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (c *SubscriptionRepository) Update(ctx context.Context, query string, values []any) (bool, error) {
	tag, err := c.db.Exec(ctx, query, values...)

	if err != nil {
		return false, err
	}

	return tag.RowsAffected() != 0, nil
}
func (c *SubscriptionRepository) FindById(ctx context.Context, id uuid.UUID) (*dto.Subscription, bool, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM public.subscriptions
		WHERE id = $1
	`

	item := &dto.Subscription{}

	err := c.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.ServiceName,
		&item.Price,
		&item.UserID,
		&item.StartDate,
		&item.EndDate,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return item, true, nil
}

func (c *SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	query := "delete from public.Subscriptions where id = $1"
	tag, err := c.db.Exec(ctx, query, id)

	if err != nil {
		return false, err
	}

	return tag.RowsAffected() != 0, nil
}
func (c *SubscriptionRepository) FindAll(ctx context.Context, offset, limit int) ([]*dto.Subscription, int, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM public.subscriptions
		ORDER BY created_at DESC
		OFFSET $1 LIMIT $2;
	`

	rows, err := c.db.Query(ctx, query, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var subscriptions []*dto.Subscription
	for rows.Next() {
		item := &dto.Subscription{}
		err := rows.Scan(
			&item.ID,
			&item.ServiceName,
			&item.Price,
			&item.UserID,
			&item.StartDate,
			&item.EndDate,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		subscriptions = append(subscriptions, item)
	}

	// Получаем общее количество подписок
	var total int
	countQuery := `SELECT COUNT(*) FROM public.subscriptions;`
	err = c.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}

func (c *SubscriptionRepository) Create(ctx context.Context, ci *dto.Subscription) (*dto.Subscription, error) {
	query := "insert into public.Subscriptions (service_name, start_date, price, end_date, user_id) values ($1, $2, $3, $4, $5) returning id"
	err := c.db.QueryRow(ctx, query,
		ci.ServiceName,
		ci.StartDate,
		ci.Price,
		ci.EndDate,
		ci.UserID).Scan(&ci.ID)

	if err != nil {
		return ci, err
	}

	return ci, nil
}
