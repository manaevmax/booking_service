package repository

import (
	"sync"
	"time"

	"hotel/internal/entity"
	"hotel/internal/utils"
)

var availability = []entity.RoomAvailability{
	{"reddison", "lux", utils.Date(2024, 9, 1), 1},
	{"reddison", "lux", utils.Date(2024, 9, 2), 1},
	{"reddison", "lux", utils.Date(2024, 9, 3), 1},
	{"reddison", "lux", utils.Date(2024, 9, 4), 1},
	{"reddison", "lux", utils.Date(2024, 9, 5), 0},
}

type AvailabilityRepository struct {
	availability []entity.RoomAvailability

	rwLock sync.RWMutex
}

func NewAvailabilityRepository() *AvailabilityRepository {
	return &AvailabilityRepository{
		availability: availability,
	}
}

func (r *AvailabilityRepository) GetAvailableDates(hotelID, roomID string) []time.Time {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()

	var availableDates []time.Time

	for _, a := range r.availability {
		if a.HotelID == hotelID && a.RoomID == roomID && a.Quota > 0 {
			availableDates = append(availableDates, a.Date)
		}
	}

	return availableDates
}

func (r *AvailabilityRepository) OccupyDates(hotelID, roomID string, from, to time.Time) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	occupiedDates := make(map[time.Time]struct{})
	for _, d := range utils.PeriodToDateList(from, to) {
		occupiedDates[d] = struct{}{}
	}

	for i, a := range r.availability {
		if _, ok := occupiedDates[a.Date]; ok && a.HotelID == hotelID && a.RoomID == roomID {
			r.availability[i].Quota -= 1
		}
	}
}
