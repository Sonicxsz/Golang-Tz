package httpHelpers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const (
	Error500          = "Something went wrong, try later..."
	ErrorParse        = "Cant parse data, please check provided data"
	ErrorNotFoundById = "Nothing found, please check the provided id"
)

type ServiceError struct {
	Code    int
	Message string
	Err     error
}

func NewServiceError(code int, message string) *ServiceError {
	return &ServiceError{Code: code, Message: message}
}

type SuccessMessage struct {
	Status  int  `json:"status"`
	Data    any  `json:"data,omitempty"`
	Success bool `json:"success"`
}

type ErrorMessage struct {
	Status  int       `json:"status"`
	Error   string    `json:"error"`
	Success bool      `json:"success"`
	Id      uuid.UUID `json:"id"`
	Time    time.Time `json:"time"`
}

func NewSuccessMessage(status int, data any) *SuccessMessage {
	return &SuccessMessage{
		Status:  status,
		Data:    data,
		Success: true,
	}
}

func NewErrorMessage(error string, status int) *ErrorMessage {
	return &ErrorMessage{
		Status:  status,
		Error:   error,
		Success: false,
		Id:      uuid.New(),
		Time:    time.Now(),
	}
}

func RespondSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(NewSuccessMessage(status, data))
}

func RespondError(w http.ResponseWriter, status int, error string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(NewErrorMessage(error, status))
}
