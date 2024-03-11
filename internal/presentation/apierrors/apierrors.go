package apierrors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/httplog/v2"
)

const (
	missingRequiredFieldsMessage = "Missing required fields: %s"
	internalServerErrorMessage   = "Internal server error"
	notAvailableMessage          = "Room not available for selected dates"
	roomNotFoundMessage          = "Room not found"
)

type HTTPError struct {
	Code    int      `json:"-"`
	Message string   `json:"message"`
	Details []Detail `json:"details,omitempty"`
	Err     error    `json:"-"`
}

type Detail struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

func NewBadRequestError(msg string, err error) *HTTPError {
	return &HTTPError{Code: http.StatusBadRequest, Message: msg, Err: err}
}

func NewUnprocessableEntityError(msg string, err error) *HTTPError {
	return &HTTPError{
		Code:    http.StatusUnprocessableEntity,
		Message: msg,
		Err:     err,
	}
}

func NewInternalServerError(err error) *HTTPError {
	return &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: internalServerErrorMessage,
		Err:     err,
	}
}

type ValidationError struct {
	MissingFields []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf(
		missingRequiredFieldsMessage,
		strings.Join(e.MissingFields, ", "),
	)
}

func NewValidationError(missingFields []string) *HTTPError {
	return &HTTPError{
		Code:    http.StatusUnprocessableEntity,
		Message: (&ValidationError{MissingFields: missingFields}).Error(),
	}
}

func NewStatusConflictError(err error) *HTTPError {
	return &HTTPError{
		Code:    http.StatusConflict,
		Message: notAvailableMessage,
		Err:     err,
	}
}

func NewRoomNotFoundError(err error) *HTTPError {
	return &HTTPError{
		Code:    http.StatusNotFound,
		Message: roomNotFoundMessage,
		Err:     err,
	}
}
func HandleError(w http.ResponseWriter, err error, ctx context.Context) {
	oplog := httplog.LogEntry(ctx)
	oplog.Error("Error in processing request", "error", err)

	w.Header().Set("Content-Type", "application/json")

	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		w.WriteHeader(httpErr.Code)

		err := json.NewEncoder(w).Encode(httpErr)
		if err != nil {
			oplog.Error("Unable to encode response", "error", err)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)

		err := json.NewEncoder(w).Encode(NewInternalServerError(err))
		if err != nil {
			oplog.Error("Unable to encode response", "error", err)
		}
	}
}
