package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"hotel/internal/entity"
	"hotel/internal/http"
	"hotel/internal/utils"
)

func TestBookingService_BookRoom(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("success", func(t *testing.T) {
		orderRepo := NewMockOrderRepo(ctrl)
		orderRepo.EXPECT().CreateOrder(entity.Order{
			HotelID:   "reddison",
			RoomID:    "lux",
			UserEmail: "my@mail.ru",
			From:      utils.Date(2024, 1, 1),
			To:        utils.Date(2024, 1, 2),
		})

		availabilityRepo := NewMockAvailabilityRepo(ctrl)
		availabilityRepo.EXPECT().
			GetAvailableDates("reddison", "lux").
			Return([]time.Time{
				utils.Date(2024, 1, 1),
				utils.Date(2024, 1, 2),
			})
		availabilityRepo.EXPECT().
			OccupyDates("reddison", "lux", utils.Date(2024, 1, 1), utils.Date(2024, 1, 2)).
			Return()

		orderDto := http.OrderDto{
			HotelID:   "reddison",
			RoomID:    "lux",
			UserEmail: "my@mail.ru",
			From:      utils.Date(2024, 1, 1),
			To:        utils.Date(2024, 1, 2),
		}

		srv := New(orderRepo, availabilityRepo)
		err := srv.BookRoom(orderDto)

		assert.NoError(t, err)
	})

	t.Run("no available days", func(t *testing.T) {
		orderRepo := NewMockOrderRepo(ctrl)
		availabilityRepo := NewMockAvailabilityRepo(ctrl)
		availabilityRepo.EXPECT().
			GetAvailableDates("reddison", "lux").
			Return([]time.Time{utils.Date(2024, 1, 1)})

		orderDto := http.OrderDto{
			HotelID:   "reddison",
			RoomID:    "lux",
			UserEmail: "my@mail.ru",
			From:      utils.Date(2024, 1, 1),
			To:        utils.Date(2024, 1, 2),
		}

		srv := New(orderRepo, availabilityRepo)
		err := srv.BookRoom(orderDto)

		assert.EqualError(t, err, "no available dates")
	})

	t.Run("failed to create order", func(t *testing.T) {
		orderRepo := NewMockOrderRepo(ctrl)
		orderRepo.EXPECT().
			CreateOrder(entity.Order{
				HotelID:   "reddison",
				RoomID:    "lux",
				UserEmail: "my@mail.ru",
				From:      utils.Date(2024, 1, 1),
				To:        utils.Date(2024, 1, 2),
			}).
			Return(fmt.Errorf("error"))

		availabilityRepo := NewMockAvailabilityRepo(ctrl)
		availabilityRepo.EXPECT().
			GetAvailableDates("reddison", "lux").
			Return([]time.Time{utils.Date(2024, 1, 1), utils.Date(2024, 1, 2)})

		orderDto := http.OrderDto{
			HotelID:   "reddison",
			RoomID:    "lux",
			UserEmail: "my@mail.ru",
			From:      utils.Date(2024, 1, 1),
			To:        utils.Date(2024, 1, 2),
		}

		srv := New(orderRepo, availabilityRepo)
		err := srv.BookRoom(orderDto)

		assert.EqualError(t, err, "failed to create order: error")
	})
}
