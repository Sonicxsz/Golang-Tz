package dto

import (
	"awesomeProject1/pkg/validator"
	"time"
)

// GetTotalSumRequest — DTO для получения суммарной стоимости
type GetTotalSumRequest struct {
	Start       time.Time `example:"01-2025"`
	End         time.Time `example:"12-2025"`
	UserId      string    `example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName string    `example:"Yandex Plus"`
}

// IsValid — валидация параметров запроса
func (r *GetTotalSumRequest) IsValid() (bool, []string) {
	v := validator.New()

	if r.UserId != "" {
		v.CheckString(r.UserId, "UserId").IsUuid()
	}

	if r.ServiceName != "" {
		v.CheckString(r.ServiceName, "ServiceName").IsMin(1).IsMax(255)
	}

	return !v.HasErrors(), v.GetErrors()
}

// NewGetTotalSumRequest — конструктор
func NewGetTotalSumRequest(start, end time.Time, userId, serviceName string) *GetTotalSumRequest {
	return &GetTotalSumRequest{
		Start:       start,
		End:         end,
		UserId:      userId,
		ServiceName: serviceName,
	}
}
