package dto

import (
	"awesomeProject1/pkg/validator"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// UpdateSubscriptionRequest — DTO для обновления подписки
type UpdateSubscriptionRequest struct {
	ID          string  `json:"id"`
	ServiceName *string `json:"service_name,omitempty" example:"Yandex Plus"`
	Price       *int    `json:"price,omitempty" example:"400"`
	UserID      *string `json:"user_id,omitempty" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   *string `json:"start_date,omitempty" example:"01-2025"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2025"`
}

// UpdateData — структура для передачи обновлённых данных в слой репозитория
type UpdateData struct {
	ID          uuid.UUID
	ServiceName *string
	Price       *int
	UserID      *uuid.UUID
	StartDate   *time.Time
	EndDate     *sql.NullTime
}

// IsValid проверяет корректность данных запроса
func (c *UpdateSubscriptionRequest) IsValid() (bool, []string) {
	v := validator.New()
	v.CheckString(c.ID, "id").IsUuid()

	if c.ServiceName != nil {
		v.CheckString(*c.ServiceName, "ServiceName").IsMin(1).IsMax(255)
	}

	if c.Price != nil {
		v.CheckNumber(*c.Price, "Price").IsMin(0)
	}

	if c.UserID != nil {
		v.CheckString(*c.UserID, "UserID").IsUuid()
	}

	var startDate time.Time
	var startDateValid bool

	if c.StartDate != nil {
		parsed, err := ParseMonthYear(*c.StartDate)
		if err != nil {
			v.AddError(fmt.Sprintf("Invalid start_date format. Expected MM-YYYY (e.g., 01-2025). Got: %s", *c.StartDate))
		} else {
			startDate = parsed
			startDateValid = true
		}
	}

	if c.EndDate != nil {
		endDate, err := ParseMonthYear(*c.EndDate)
		if err != nil {
			v.AddError(fmt.Sprintf("Invalid end_date format. Expected MM-YYYY (e.g., 12-2025). Got: %s", *c.EndDate))
		} else if startDateValid && endDate.Before(startDate) {
			v.AddError(fmt.Sprintf("end_date must be after start_date. Got: end_date=%s, start_date=%s", *c.EndDate, *c.StartDate))
		}
	}

	return !v.HasErrors(), v.GetErrors()
}

// ToUpdateData конвертирует DTO в структуру UpdateData
func (c *UpdateSubscriptionRequest) ToUpdateData() (*UpdateData, error) {
	id, err := uuid.Parse(c.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse id: %w", err)
	}

	data := &UpdateData{
		ID:          id,
		ServiceName: c.ServiceName,
		Price:       c.Price,
	}

	if c.UserID != nil {
		userID, err := uuid.Parse(*c.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse user_id: %w", err)
		}
		data.UserID = &userID
	}

	if c.StartDate != nil {
		startDate, err := ParseMonthYear(*c.StartDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse start_date: %w", err)
		}
		data.EndDate = &sql.NullTime{Valid: false}
		data.StartDate = &startDate
	}

	if c.EndDate != nil && *c.EndDate != "" {
		endDate, err := ParseMonthYear(*c.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse end_date: %w", err)
		}
		data.EndDate = &sql.NullTime{Time: endDate, Valid: true}
	}

	return data, nil
}
