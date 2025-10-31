package dto

import (
	"awesomeProject1/pkg/validator"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// CreateSubscriptionRequest — DTO для создания подписки
type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" example:"Yandex Plus"`
	Price       int    `json:"price" example:"400"`
	UserID      string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string `json:"start_date" example:"01-2025"` // Формат MM-YYYY
	EndDate     string `json:"end_date,omitempty" example:"12-2025"`
}

func (r *CreateSubscriptionRequest) IsValid() (bool, []string) {
	v := validator.New()
	v.CheckString(r.ServiceName, "ServiceName").IsMin(1).IsMax(255)
	v.CheckString(r.UserID, "UserID").IsUuid()
	v.CheckNumber(r.Price, "Price").IsMin(0)

	startDate, err := ParseMonthYear(r.StartDate)
	if err != nil {
		v.AddError(fmt.Sprintf("Invalid start_date format. Expected MM-YYYY (e.g., 01-2025). Got: %s", r.StartDate))
	}

	if r.EndDate != "" {
		endDate, err := ParseMonthYear(r.EndDate)
		if err != nil {
			v.AddError(fmt.Sprintf("Invalid end_date format. Expected MM-YYYY (e.g., 12-2025). Got: %s", r.EndDate))
		}

		if err == nil && endDate.Before(startDate) {
			v.AddError(fmt.Sprintf("end_date must be after start_date. Got: end_date=%s, start_date=%s", r.EndDate, r.StartDate))
		}
	}

	return !v.HasErrors(), v.GetErrors()
}

// ToSubscription конвертирует DTO в модель Subscription
func (r *CreateSubscriptionRequest) ToSubscription() (*Subscription, error) {
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user_id: %w", err)
	}

	startDate, err := ParseMonthYear(r.StartDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start_date: %w", err)
	}

	subscription := &Subscription{
		ServiceName: r.ServiceName,
		Price:       r.Price,
		UserID:      userID,
		StartDate:   startDate,
	}

	if r.EndDate != "" {
		endDate, err := ParseMonthYear(r.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse end_date: %w", err)
		}
		subscription.EndDate = sql.NullTime{Time: endDate, Valid: true}
	}

	return subscription, nil
}
