package apihandlers

import (
	"applicationDesignTest/internal/presentation/apierrors"
	"applicationDesignTest/internal/presentation/apimodels"
	"context"
	"encoding/json"
	"net/http"
)

const (
	// Error messages that used only in middlewares:
	DecodeFailureMessage          = "Failed to decode request"
	InvalidRequestDataMessage     = "Invalid request data"
	DatesMustBeInUTCMessage       = "`from` and `to` dates must be in UTC timezone"
	FromDateMustBeBeforeToMessage = "`from` date must be before `to` date"
)

func ValidateOrdersAPIRequest(requestBody *apimodels.OrdersAPIRequest) error {
	var missingFields []string
	if requestBody.HotelID == "" {
		missingFields = append(missingFields, "hotel_id")
	}
	if requestBody.RoomID == "" {
		missingFields = append(missingFields, "room_id")
	}
	if requestBody.UserEmail == "" {
		missingFields = append(missingFields, "email")
	}
	if requestBody.From.IsZero() {
		missingFields = append(missingFields, "from")
	}
	if requestBody.To.IsZero() {
		missingFields = append(missingFields, "to")
	}

	if len(missingFields) > 0 {
		return apierrors.NewValidationError(missingFields)
	}

	_, fromTimeZoneOffset := requestBody.From.Zone()
	_, toTimeZoneOffset := requestBody.To.Zone()
	if fromTimeZoneOffset != 0 || toTimeZoneOffset != 0 {
		return apierrors.NewBadRequestError(DatesMustBeInUTCMessage, nil)
	}

	if requestBody.From.After(requestBody.To) {
		return apierrors.NewBadRequestError(
			FromDateMustBeBeforeToMessage,
			nil,
		)
	}

	return nil
}

func OrdersJsonDecoderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// NOTE:
		// No need to defer closing the r.Body.
		// Stdlib server automatically handles closing the request body
		// after the handler has returned.
		var req apimodels.OrdersAPIRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			apierrors.HandleError(
				w,
				apierrors.NewBadRequestError(DecodeFailureMessage, err),
				ctx,
			)
			return
		}

		ctx = context.WithValue(ctx, ApiRequestBodyKey, &req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OrdersValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := getFromContext(ctx)
		if err != nil {
			apierrors.HandleError(w, err, ctx)
			return
		}

		if r.Method == http.MethodPost {
			if err := ValidateOrdersAPIRequest(req); err != nil {
				apierrors.HandleError(
					w,
					apierrors.NewUnprocessableEntityError(err.Error(), err),
					ctx,
				)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
