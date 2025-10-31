package dto

import (
	"github.com/google/uuid"
	"time"
)

// SubscriptionResponse — DTO для ответа API
type SubscriptionResponse struct {
	ID          uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ServiceName string    `json:"service_name" example:"Yandex Plus"`
	Price       int       `json:"price" example:"400"`
	UserID      uuid.UUID `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string    `json:"start_date" example:"01-2025"` // Формат: MM-YYYY
	EndDate     *string   `json:"end_date,omitempty" example:"12-2025"`
	CreatedAt   time.Time `json:"created_at" example:"2025-10-28T10:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2025-10-28T10:00:00Z"`
}

// SubscriptionListResponse — DTO для списка подписок
type SubscriptionListResponse struct {
	Total         int                     `json:"total" example:"100"`
	Offset        int                     `json:"offset" example:"0"`
	Limit         int                     `json:"limit" example:"10"`
	Subscriptions []*SubscriptionResponse `json:"subscriptions"`
}
