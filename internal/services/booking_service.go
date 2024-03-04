package services

import (
	"applicationDesignTest/internal/infrastructure/unitofwork"
	"applicationDesignTest/internal/repositories"
	"time"
)

type BookingService struct {
	uow  unitofwork.MutexUnitOfWork
	repo *repositories.AvailabilityRepository
}

func NewBookingService(
	uow unitofwork.MutexUnitOfWork,
	repo *repositories.AvailabilityRepository,
) *BookingService {
	return &BookingService{
		uow:  uow,
		repo: repo,
	}
}

func (s *BookingService) CreateOrder(
	hotelID, roomID string,
	from, to time.Time,
) error {
	s.uow.Begin()
	defer s.uow.Commit()

	if err := s.repo.SetRoomPending(hotelID, roomID, from, to); err != nil {
		return err
	}
	if err := s.repo.Book(hotelID, roomID, from, to); err != nil {
		return err
	}

	return nil
}
