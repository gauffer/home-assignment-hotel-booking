package httphandlers

import (
	"applicationDesignTest/internal/presentation/httpmodels"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

func OrdersJsonDecoderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data httpmodels.OrdersAPIRequest

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "jsonBody", &data)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OrdersValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// OrderAPIRequest это структура без методов
		jsonBody, ok := r.Context().Value("jsonBody").(*httpmodels.OrdersAPIRequest)
		if !ok {
			http.Error(w, "Invalid request data", http.StatusBadRequest)
			return
		}
		if r.Method == http.MethodPost {
			if err := ValidateOrdersAPIRequest(jsonBody); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// TODO: проверить что время в UTC и написать тест
func ValidateOrdersAPIRequest(jsonBody *httpmodels.OrdersAPIRequest) error {
	if jsonBody.HotelID == "" {
		return errors.New("hotel_id is required")
	}
	if jsonBody.RoomID == "" {
		return errors.New("room_id is required")
	}
	if jsonBody.UserEmail == "" {
		return errors.New("email is required")
	}
	if jsonBody.From.After(jsonBody.To) {
		return errors.New("from date must be before to date")
	}
	return nil
}
