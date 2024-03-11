package domains

import (
	"applicationDesignTest/internal/presentation/apimodels"
	"time"
)

type Order struct {
	HotelID   string
	RoomID    string
	UserEmail string
	From      time.Time
	To        time.Time
}

func OrderFromAPIRequest(
	apiRequestBody *apimodels.OrdersAPIRequest,
) Order {
	order := Order{
		HotelID:   apiRequestBody.HotelID,
		RoomID:    apiRequestBody.RoomID,
		UserEmail: apiRequestBody.UserEmail,
		From:      apiRequestBody.From,
		To:        apiRequestBody.To,
	}

	return order

}
