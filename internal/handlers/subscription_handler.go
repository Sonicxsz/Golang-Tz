package handlers

import (
	"awesomeProject1/internal/dto"
	"awesomeProject1/internal/service"
	"awesomeProject1/pkg/httpHelpers"
	"awesomeProject1/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SubscriptionHandler struct {
	service service.ISubscriptionService
}

func NewCatalogHandler(service service.ISubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// GetTotal возвращает общую сумму по подпискам за указанный период.
//
// @Summary      Получить общую сумму подписок
// @Description  Возвращает суммарную стоимость подписок за указанный период с возможностью фильтрации по пользователю и сервису
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        start        query  string  true   "Дата начала периода (MM-YYYY)"  example("01-2025")
// @Param        end          query  string  true   "Дата окончания периода (MM-YYYY)"  example("12-2025")
// @Param        user_id      query  string  false  "ID пользователя (UUID)"  example("60601fee-2bf1-4721-ae6f-7636e79a0cba")
// @Param        service_name query  string  false  "Название сервиса"  example("Yandex Plus")
// @Success      200 {object} map[string]int
// @Failure      400  {object}  httpHelpers.ErrorMessage
// @Failure      500  {object}  httpHelpers.ErrorMessage
// @Router       /subscriptions/total [get]
func (c *SubscriptionHandler) GetTotal(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	start, err := dto.ParseMonthYear(params.Get("start"))

	if err != nil {
		httpHelpers.RespondError(w, http.StatusBadRequest, "Please provide start param in next format: mm-yyyy")
		return
	}

	end, err := dto.ParseMonthYear(params.Get("end"))
	if err != nil {
		httpHelpers.RespondError(w, http.StatusBadRequest, "Please provide end param in next format: mm-yyyy")
		return
	}

	serviceName := params.Get("service_name")
	userId := params.Get("user_id")

	req := dto.NewGetTotalSumRequest(start, end, userId, serviceName)

	if ok, errors := req.IsValid(); !ok {
		httpHelpers.RespondError(w, http.StatusBadRequest, strings.Join(errors, "; "))
		return
	}

	sum, sErr := c.service.GetTotalSum(r.Context(), req)

	if sErr != nil {
		httpHelpers.RespondError(w, sErr.Code, sErr.Message)
		return
	}

	httpHelpers.RespondSuccess(w, http.StatusOK, sum)
}

// Delete удаляет подписку по ID.
//
// @Summary      Удалить подписку
// @Description  Удаляет подписку по её уникальному идентификатору
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id path string true "ID подписки" example("123e4567-e89b-12d3-a456-426614174000")
// @Success      200  {object}  httpHelpers.SuccessMessage
// @Failure      400  {object}  httpHelpers.ErrorMessage
// @Failure      500  {object}  httpHelpers.ErrorMessage
// @Router       /subscription/{id} [delete]
func (c *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		httpHelpers.RespondError(w, http.StatusBadRequest, "Subscription id not provided")
		return
	}

	parsedId, err := uuid.Parse(id)

	if err != nil {
		httpHelpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Cannot parse provided id. Expected correct uuid. Got: %s", id))
		return
	}

	sError := c.service.Delete(r.Context(), parsedId)

	if sError != nil {
		httpHelpers.RespondError(w, sError.Code, sError.Message)
		return
	}

	httpHelpers.RespondSuccess(w, http.StatusOK, nil)
}

// GetById возвращает подписку по ID.
//
// @Summary      Получить подписку по ID
// @Description  Возвращает данные конкретной подписки по её уникальному идентификатору
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id path string true "ID подписки" example("123e4567-e89b-12d3-a456-426614174000")
// @Success      200 {object} dto.SubscriptionResponse
// @Failure      400  {object}  httpHelpers.ErrorMessage
// @Failure      500  {object}  httpHelpers.ErrorMessage
// @Router       /subscription/{id} [get]
func (c *SubscriptionHandler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		httpHelpers.RespondError(w, http.StatusBadRequest, "Subscription id not provided")
		return
	}

	parsedId, err := uuid.Parse(id)

	if err != nil {
		httpHelpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Cannot parse provided id. Expected correct uuid. Got: %s", id))
		return
	}

	item, sErr := c.service.GetById(r.Context(), parsedId)

	if sErr != nil {
		httpHelpers.RespondError(w, sErr.Code, sErr.Message)
		return
	}

	httpHelpers.RespondSuccess(w, http.StatusOK, item)
}

// Update обновляет существующую подписку.
//
// @Summary      Обновить подписку
// @Description  Обновляет данные подписки (частично или полностью)
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request body dto.UpdateSubscriptionRequest true "Данные для обновления подписки"
// @Success      200  {object} 	httpHelpers.SuccessMessage
// @Failure      400  {object}  httpHelpers.ErrorMessage
// @Failure      500  {object}  httpHelpers.ErrorMessage
// @Router       /subscription [put]
func (c *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	req := dto.UpdateSubscriptionRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpHelpers.RespondError(w, http.StatusBadRequest, httpHelpers.ErrorParse)
		return
	}

	if ok, errors := req.IsValid(); !ok {
		httpHelpers.RespondError(w, http.StatusBadRequest, strings.Join(errors, "; "))
		return
	}

	updateData, err := req.ToUpdateData()
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Subscription handler -> ToUpdateData Error -> err: %s", err.Error()))
		httpHelpers.RespondError(w, http.StatusInternalServerError, "Что-то пошло не так, попробуйте позже или проверьте данные")
		return
	}

	sError := c.service.Update(r.Context(), updateData)
	if sError != nil {
		httpHelpers.RespondError(w, sError.Code, sError.Message)
		return
	}

	httpHelpers.RespondSuccess(w, http.StatusOK, nil)
}

// Create создаёт новую подписку.
//
// @Summary      Создать подписку
// @Description  Создаёт новую подписку для пользователя
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateSubscriptionRequest true "Данные для создания подписки"
// @Success      201  {object}  dto.SubscriptionResponse
// @Failure      400  {object}  httpHelpers.ErrorMessage
// @Failure      500  {object}  httpHelpers.ErrorMessage
// @Router       /subscription [post]
func (c *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	req := dto.CreateSubscriptionRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpHelpers.RespondError(w, http.StatusBadRequest, httpHelpers.ErrorParse)
		return
	}

	if ok, errors := req.IsValid(); !ok {
		httpHelpers.RespondError(w, 400, strings.Join(errors, "; "))
		return
	}

	sub, err := req.ToSubscription()

	if err != nil {
		logger.Log.Error(fmt.Sprintf("Subscription handler -> ToSubscription Error -> err: %s", err.Error()))
		httpHelpers.RespondError(w, 500, "Что-то пошло не так, попробуйте позже или проверьте данные")
		return
	}

	created, sError := c.service.Create(r.Context(), sub)

	if sError != nil {
		httpHelpers.RespondError(w, sError.Code, sError.Message)
		return
	}

	httpHelpers.RespondSuccess(w, http.StatusCreated, created)
}

// GetAll возвращает список всех подписок.
//
// @Summary      Получить список подписок
// @Description  Возвращает список всех подписок с пагинацией
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        offset  query  int  false  "Смещение (по умолчанию 0)"  example(0)
// @Param        limit   query  int  false  "Лимит записей (по умолчанию 10)" example(10)
// @Success      200  {object}  dto.SubscriptionListResponse
// @Failure      400  {object}  httpHelpers.ErrorMessage
// @Failure      500  {object}  httpHelpers.ErrorMessage
// @Router       /subscriptions [get]
func (c *SubscriptionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	offsetStr := query.Get("offset")
	limitStr := query.Get("limit")

	offset := 0
	limit := 10

	if offsetStr != "" {
		fmt.Sscanf(offsetStr, "%d", &offset)
	}
	if limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}

	result, sErr := c.service.GetAll(r.Context(), offset, limit)
	if sErr != nil {
		httpHelpers.RespondError(w, sErr.Code, sErr.Message)
		return
	}

	httpHelpers.RespondSuccess(w, http.StatusOK, result)
}
