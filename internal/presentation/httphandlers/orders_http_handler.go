package httphandlers

import (
	"applicationDesignTest/internal/presentation/httpmodels"
	"applicationDesignTest/internal/services"
	"encoding/json"
	"net/http"
)

func CreateOrderHandler(bookingService *services.BookingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonBody, ok := r.Context().Value("jsonBody").(*httpmodels.OrdersAPIRequest)
		if !ok {
			http.Error(w, "Invalid request data", http.StatusBadRequest)
			return
		}

		if err := bookingService.CreateOrder(jsonBody.HotelID, jsonBody.RoomID, jsonBody.From, jsonBody.To); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(jsonBody)
	}
}
