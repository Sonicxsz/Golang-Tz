package dto

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

// Subscription — базовая модель подписки в БД
type Subscription struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	ServiceName string       `json:"service_name" db:"service_name"`
	Price       int          `json:"price" db:"price"`
	UserID      uuid.UUID    `json:"user_id" db:"user_id"`
	StartDate   time.Time    `json:"start_date" db:"start_date"`
	EndDate     sql.NullTime `json:"end_date,omitempty" db:"end_date"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}

// ToResponse конвертирует Subscription в SubscriptionResponse для API
func (s *Subscription) ToResponse() *SubscriptionResponse {
	response := SubscriptionResponse{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID,
		StartDate:   FormatMonthYear(s.StartDate),
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}

	if s.EndDate.Valid {
		endDate := FormatMonthYear(s.EndDate.Time)
		response.EndDate = &endDate
	}

	return &response
}
