package http

import (
	"fmt"
	"time"
)

type OrderDto struct {
	HotelID   string    `json:"hotel_id"`
	RoomID    string    `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}

func (req *OrderDto) Validate() error {
	if req.HotelID == "" {
		return fmt.Errorf("hotelID must not be empty")
	}

	if req.RoomID == "" {
		return fmt.Errorf("roomID must not be empty")
	}

	if req.UserEmail == "" {
		return fmt.Errorf("userEmail must not be empty")
	}

	if req.From.After(req.To) {
		return fmt.Errorf("to must be after from")
	}

	if req.From.Before(time.Now()) {
		return fmt.Errorf("from must be after now")
	}

	if req.To.Before(time.Now()) {
		return fmt.Errorf("to must be after now")
	}

	return nil
}
