package main

import (
	"applicationDesignTest/internal/infrastructure/logger"
	"applicationDesignTest/internal/infrastructure/unitofwork"
	"applicationDesignTest/internal/presentation/httphandlers"
	"applicationDesignTest/internal/repositories"
	"applicationDesignTest/internal/services"
	"errors"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

func SetupMain(logger *httplog.Logger) *chi.Mux {
	r := chi.NewRouter()
	if logger != nil {
		r.Use(httplog.RequestLogger(logger))
	}
	r.Use(middleware.Heartbeat("/ping"))
	return r
}

func SetupOrders(r *chi.Mux) *chi.Mux {
	repo := repositories.NewAvailabilityRepository()
	var mutex sync.Mutex
	uow := unitofwork.NewMutexUnitOfWork(&mutex)
	bookingService := services.NewBookingService(*uow, repo)

	r.With(httphandlers.OrdersJsonDecoderMiddleware, httphandlers.OrdersValidationMiddleware).
		Post("/orders", httphandlers.CreateOrderHandler(bookingService))
	return r
}

func main() {
	logger := logger.SetupLogger()
	r := SetupMain(logger)
	r = SetupOrders(r)

	logger.Info("Server listening on localhost:8080")
	err := http.ListenAndServe(":8080", r)
	if errors.Is(err, http.ErrServerClosed) {
		logger.Info("Server closed")
	} else if err != nil {
		logger.Error("Server failed: %s", err)
		panic("Unexpected error")
	}
	// TODO: SIGINT, SIGTERM
}
