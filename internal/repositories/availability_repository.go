package repositories

import (
	"applicationDesignTest/internal/domains"
	"errors"
	"time"
)

// TODO: fix and start using
type availabilityRepository interface {
	SetRoomPending(
		hotelID, roomID string,
		from, to time.Time,
	) (bool, error)
	Book(
		hotelID, roomID string,
		from, to time.Time,
	) (*domains.RoomAvailability, error)
}

type AvailabilityRepository struct {
	data []domains.RoomAvailability
}

var ErrRoomNotFound = errors.New("room not found")
var ErrRoomNotAvailable = errors.New("room not available for selected dates")

func isAvailabilityMatch(
	availability domains.RoomAvailability,
	hotelID, roomID string,
	date time.Time,
) bool {
	return availability.HotelID == hotelID && availability.RoomID == roomID &&
		availability.Date.Equal(date)
}

func (r *AvailabilityRepository) hasAvailabilityForDate(
	hotelID, roomID string,
	date time.Time,
) bool {
	for _, availability := range r.data {
		if isAvailabilityMatch(availability, hotelID, roomID, date) &&
			availability.Quota > 0 {
			return true
		}
	}
	return false
}

func (r *AvailabilityRepository) reduceQuotaForDate(
	hotelID, roomID string,
	date time.Time,
) error {
	for i, availability := range r.data {
		if isAvailabilityMatch(availability, hotelID, roomID, date) {
			r.data[i].Quota -= 1
			return nil
		}
	}
	return ErrRoomNotFound
}

func (r *AvailabilityRepository) SetRoomPending(
	hotelID, roomID string,
	from, to time.Time,
) error {
	from = toDay(from)
	to = toDay(to)

	for date := from; !date.After(to); date = date.AddDate(0, 0, 1) {
		if !r.hasAvailabilityForDate(hotelID, roomID, date) {
			return ErrRoomNotAvailable
		}
	}

	return nil
}

func (r *AvailabilityRepository) Book(
	hotelID, roomID string,
	from, to time.Time,
) error {
	from = toDay(from)
	to = toDay(to)

	for date := from; !date.After(to); date = date.AddDate(0, 0, 1) {
		if err := r.reduceQuotaForDate(hotelID, roomID, date); err != nil {
			return err
		}
	}

	return nil
}

func NewAvailabilityRepository() *AvailabilityRepository {
	initialData := []domains.RoomAvailability{
		{HotelID: "reddison", RoomID: "lux", Date: date(2024, 1, 1), Quota: 1},
		{HotelID: "reddison", RoomID: "lux", Date: date(2024, 1, 2), Quota: 1},
		{HotelID: "reddison", RoomID: "lux", Date: date(2024, 1, 3), Quota: 1},
		{HotelID: "reddison", RoomID: "lux", Date: date(2024, 1, 4), Quota: 1},
		{HotelID: "reddison", RoomID: "lux", Date: date(2024, 1, 5), Quota: 0},
	}
	return &AvailabilityRepository{data: initialData}
}

func toDay(timestamp time.Time) time.Time {
	return time.Date(
		timestamp.Year(),
		timestamp.Month(),
		timestamp.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
