package service

//go:generate mockgen -package $GOPACKAGE -source deps.go -destination mocks.go

import (
	"time"

	"hotel/internal/entity"
)

type OrderRepo interface {
	CreateOrder(order entity.Order) error
}

type AvailabilityRepo interface {
	GetAvailableDates(hotelID, roomID string) []time.Time
	OccupyDates(hotelID, roomID string, from, to time.Time)
}
