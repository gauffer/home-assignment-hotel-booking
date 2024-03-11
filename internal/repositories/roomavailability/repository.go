package roomavailability

import (
	"applicationDesignTest/internal/domains"
	"time"
)

type AvailabilityRepository interface {
	EnsureRoomAvailability(hotelID, roomID string, from, to time.Time) error
	Book(hotelID, roomID string, from, to time.Time) error
}

type Repository struct {
	data []domains.RoomAvailability
}

type RoomNotAvailableError struct {
}

func (e *RoomNotAvailableError) Error() string {
	return "Room not available for selected dates"
}

func NewRoomAvailabilityError() *RoomNotAvailableError {
	return &RoomNotAvailableError{}
}

func isAvailabilityMatch(
	availability domains.RoomAvailability,
	hotelID, roomID string,
	date time.Time,
) bool {
	return availability.HotelID == hotelID && availability.RoomID == roomID &&
		availability.Date.Equal(date)
}

func (r *Repository) hasAvailabilityForDate(
	hotelID, roomID string,
	date time.Time,
) bool {
	available := false
	for _, availability := range r.data {
		if isAvailabilityMatch(availability, hotelID, roomID, date) {
			if availability.Quota > 0 {
				available = true
				break
			}
		}
	}
	return available
}

func (r *Repository) reduceQuotaForDate(
	hotelID, roomID string,
	date time.Time,
) {
	for i, availability := range r.data {
		if isAvailabilityMatch(availability, hotelID, roomID, date) {
			r.data[i].Quota -= 1
			return
		}
	}
}

func (r *Repository) EnsureRoomAvailability(
	hotelID, roomID string,
	from, to time.Time,
) error {
	from = toDay(from)
	to = toDay(to)

	for date := from; !date.After(to); date = date.AddDate(0, 0, 1) {
		if !r.hasAvailabilityForDate(hotelID, roomID, date) {
			return NewRoomAvailabilityError()
		}
	}

	return nil
}

func (r *Repository) Book(
	hotelID, roomID string,
	from, to time.Time,
) error {
	from = toDay(from)
	to = toDay(to)

	for date := from; !date.After(to); date = date.AddDate(0, 0, 1) {
		r.reduceQuotaForDate(hotelID, roomID, date)
	}

	return nil
}

func NewAvailabilityRepository() *Repository {
	initialData := []domains.RoomAvailability{
		{HotelID: "reddison", RoomID: "lux", Date: date(2024, 1, 1), Quota: 1},
		{HotelID: "reddison", RoomID: "lux", Date: date(2024, 1, 2), Quota: 1},
		{HotelID: "reddison", RoomID: "lux", Date: date(2024, 1, 3), Quota: 1},
		{HotelID: "reddison", RoomID: "lux", Date: date(2024, 1, 4), Quota: 1},
		{HotelID: "reddison", RoomID: "lux", Date: date(2024, 1, 5), Quota: 0},
	}
	return &Repository{data: initialData}
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
