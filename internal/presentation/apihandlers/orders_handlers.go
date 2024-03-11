package apihandlers

import (
	"applicationDesignTest/internal/domains"
	"applicationDesignTest/internal/presentation/apierrors"
	"applicationDesignTest/internal/services"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/httplog/v2"
)

func CreateOrderHandler(
	bookingService services.BookingService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)

		req, err := getFromContext(ctx)
		if err != nil {
			apierrors.HandleError(w, err, ctx)
			return
		}

		order := domains.OrderFromAPIRequest(req)
		if err := bookingService.CreateOrder(order); err != nil {
			handleServiceErrors(w, err, ctx)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(req)
		if err != nil {
			oplog.Error("Unable to encode response", "error", err)
		}
	}
}

func handleServiceErrors(
	w http.ResponseWriter,
	err error,
	ctx context.Context,
) {
	if serviceErr, ok := err.(services.ServiceError); ok {
		switch serviceErr.(type) {
		case *services.RoomNotFoundError:
			apierrors.HandleError(
				w,
				apierrors.NewRoomNotFoundError(err),
				ctx,
			)
			return
		case *services.RoomNotAvailableError:
			apierrors.HandleError(
				w,
				apierrors.NewStatusConflictError(err),
				ctx,
			)
			return
		default:
			apierrors.HandleError(
				w,
				apierrors.NewInternalServerError(err),
				ctx,
			)
			return
		}
	} else {
		apierrors.HandleError(
			w,
			apierrors.NewInternalServerError(err),
			ctx,
		)
	}
}
