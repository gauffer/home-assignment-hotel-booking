package services

import (
	"applicationDesignTest/internal/domains"
	"applicationDesignTest/internal/infrastructure/unitofwork"
	"applicationDesignTest/internal/repositories/roomavailability"
	"fmt"
	"time"
)

type BookingService interface {
	CreateOrder(
		order domains.Order,
	) error
}

type bookingService struct {
	uow  unitofwork.UnitOfWork
	repo roomavailability.AvailabilityRepository
}

func NewBookingService(
	uow unitofwork.UnitOfWork,
	repo roomavailability.AvailabilityRepository,
) BookingService {
	return &bookingService{
		uow:  uow,
		repo: repo,
	}
}

type ServiceError interface {
	error
}

type RoomNotFoundError struct {
}

func (e *RoomNotFoundError) Error() string {
	return "room not found"
}

func NewRoomNotFoundError(hotelID, roomID string) *RoomNotFoundError {
	return &RoomNotFoundError{}
}

type RoomNotAvailableError struct {
}

func (e *RoomNotAvailableError) Error() string {
	return "room not available for"
}

func NewRoomNotAvailableError(
	hotelID, roomID string,
	dates []time.Time,
) *RoomNotAvailableError {
	return &RoomNotAvailableError{}
}

func (s *bookingService) CreateOrder(
	order domains.Order,
) error {
	s.uow.Begin()
	defer s.uow.Commit()

	err := s.repo.EnsureRoomAvailability(
		order.HotelID,
		order.RoomID,
		order.From,
		order.To,
	)
	if err != nil {
		switch err.(type) {
		case *roomavailability.RoomNotAvailableError:
			return NewRoomNotAvailableError(order.HotelID, order.RoomID, []time.Time{})
		default:
			return fmt.Errorf("error during order creation: %w", err)
		}
	}

	if err := s.repo.Book(order.HotelID, order.RoomID, order.From, order.To); err != nil {
		return err
	}

	return nil
}
