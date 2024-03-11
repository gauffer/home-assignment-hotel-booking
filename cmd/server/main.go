package main

import (
	"applicationDesignTest/internal/infrastructure/logger"
	"applicationDesignTest/internal/infrastructure/unitofwork"
	"applicationDesignTest/internal/presentation/apihandlers"
	"applicationDesignTest/internal/repositories/roomavailability"
	"applicationDesignTest/internal/services"
	"errors"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

func setupChi(logger *httplog.Logger) *chi.Mux {
	r := chi.NewRouter()
	if logger != nil {
		r.Use(httplog.RequestLogger(logger))
	}
	r.Use(middleware.Heartbeat("/ping"))
	return r
}

func setupMutex() *sync.Mutex {
	var mutex sync.Mutex
	return &mutex
}

func setupOrdersRoute(r *chi.Mux) *chi.Mux {
	repo := roomavailability.NewAvailabilityRepository()

	// NOTE:
	// Alternative. Initialize and use PostgresUnitOfWork for the booking service.
	mu := setupMutex()
	uow := unitofwork.NewMutexUnitOfWork(mu)
	bookingService := services.NewBookingService(uow, repo)

	// NOTE:
	// CreateOrderHandler is responsible for handling HTTP request.
	// Setup for each request can be done via DI contrainers, middleware or another factory func.
	// Specifically: setting up the new unit of work per request.

	r.With(apihandlers.OrdersJsonDecoderMiddleware, apihandlers.OrdersValidationMiddleware).
		Post("/orders", apihandlers.CreateOrderHandler(bookingService))
	return r
}

func main() {
	logger := logger.SetupLogger()
	r := setupChi(logger)
	r = setupOrdersRoute(r)

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
