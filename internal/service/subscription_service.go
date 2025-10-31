package service

import (
	"awesomeProject1/internal/dto"
	"awesomeProject1/internal/repository"
	"awesomeProject1/pkg/httpHelpers"
	"awesomeProject1/pkg/logger"
	"awesomeProject1/pkg/queryBuilder"
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type ISubscriptionService interface {
	Create(cxt context.Context, req *dto.Subscription) (*dto.SubscriptionResponse, *httpHelpers.ServiceError)
	Delete(cxt context.Context, id uuid.UUID) *httpHelpers.ServiceError
	Update(cxt context.Context, req *dto.UpdateData) *httpHelpers.ServiceError
	GetById(ctx context.Context, id uuid.UUID) (*dto.SubscriptionResponse, *httpHelpers.ServiceError)
	GetTotalSum(ctx context.Context, req *dto.GetTotalSumRequest) (int, *httpHelpers.ServiceError)
	GetAll(ctx context.Context, offset, limit int) (*dto.SubscriptionListResponse, *httpHelpers.ServiceError)
}

type SubscriptionService struct {
	SubscriptionRepository repository.ISubscriptionRepository
}

func NewSubscriptionService(repo repository.ISubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{SubscriptionRepository: repo}
}

func (c *SubscriptionService) Create(cxt context.Context, req *dto.Subscription) (*dto.SubscriptionResponse, *httpHelpers.ServiceError) {
	item, err := c.SubscriptionRepository.Create(cxt, req)

	if err != nil {
		logger.Log.Error("SubscriptionService -> Create -> err -> " + err.Error())
		return nil, httpHelpers.NewServiceError(http.StatusInternalServerError, httpHelpers.Error500)
	}

	return item.ToResponse(), nil
}

func (c *SubscriptionService) Delete(ctx context.Context, id uuid.UUID) *httpHelpers.ServiceError {
	ok, err := c.SubscriptionRepository.Delete(ctx, id)

	if err != nil {
		logger.Log.Error("SubscriptionService -> Delete -> err -> " + err.Error())
		return httpHelpers.NewServiceError(http.StatusInternalServerError, httpHelpers.Error500)
	}

	if !ok {
		return httpHelpers.NewServiceError(http.StatusBadRequest, fmt.Sprintf("Cant delete Subscription by id: %s", id.String()))
	}

	return nil
}

func (c *SubscriptionService) GetById(ctx context.Context, id uuid.UUID) (*dto.SubscriptionResponse, *httpHelpers.ServiceError) {
	item, ok, err := c.SubscriptionRepository.FindById(ctx, id)

	if err != nil {
		logger.Log.Error("SubscriptionService -> GetById -> err -> " + err.Error())
		return nil, httpHelpers.NewServiceError(http.StatusInternalServerError, httpHelpers.Error500)
	}

	if !ok {
		return nil, httpHelpers.NewServiceError(http.StatusBadRequest, httpHelpers.ErrorNotFoundById)
	}

	return item.ToResponse(), nil
}

func (c *SubscriptionService) GetTotalSum(ctx context.Context, req *dto.GetTotalSumRequest) (int, *httpHelpers.ServiceError) {
	sum, err := c.SubscriptionRepository.GetTotal(ctx, req.Start, req.End, req.ServiceName, req.UserId)

	if err != nil {
		logger.Log.Error("SubscriptionService -> GetTotalSum -> err -> ", err.Error())
		return 0, httpHelpers.NewServiceError(500, httpHelpers.Error500)
	}

	return sum, nil
}

func (c *SubscriptionService) Update(cxt context.Context, req *dto.UpdateData) *httpHelpers.ServiceError {
	qb := queryBuilder.NewQueryBuilder(true).
		Set("user_id", req.UserID).
		Set("price", req.Price).
		Set("service_name", req.ServiceName).
		Set("start_date", req.StartDate).
		Set("end_date", req.EndDate)

	query, values := qb.BuildUpdateQuery("public.Subscriptions", "id", req.ID)
	ok, err := c.SubscriptionRepository.Update(cxt, query, values)

	if err != nil {
		logger.Log.Error("SubscriptionService -> Update -> err -> " + err.Error())

		return httpHelpers.NewServiceError(http.StatusInternalServerError, httpHelpers.Error500)
	}

	if !ok {
		logger.Log.Error(fmt.Sprintf("SubscriptionService -> Update -> err -> "+"Cant update Subscription item id: %d", req.ID))
		return httpHelpers.NewServiceError(http.StatusBadRequest, httpHelpers.ErrorNotFoundById)
	}

	return nil
}

func (s *SubscriptionService) GetAll(ctx context.Context, offset, limit int) (*dto.SubscriptionListResponse, *httpHelpers.ServiceError) {
	items, total, err := s.SubscriptionRepository.FindAll(ctx, offset, limit)

	if err != nil {
		logger.Log.Error(err.Error())
		return nil, httpHelpers.NewServiceError(http.StatusInternalServerError, "Ошибка при получении подписок")
	}

	responses := make([]*dto.SubscriptionResponse, len(items))
	for i, sub := range items {
		responses[i] = sub.ToResponse()
	}

	return &dto.SubscriptionListResponse{
		Total:         total,
		Offset:        offset,
		Limit:         limit,
		Subscriptions: responses,
	}, nil
}
