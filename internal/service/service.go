package service

import (
	"fmt"
	"sync"
	"time"

	"hotel/internal/entity"
	"hotel/internal/http"
	"hotel/internal/utils"
)

// BookingService сервис бронирования.
type BookingService struct {
	orderRepo        OrderRepo
	availabilityRepo AvailabilityRepo

	lock sync.Mutex
}

// New возвращает новый экземпляр сервиса.
func New(orderRepo OrderRepo, availabilityRepo AvailabilityRepo) *BookingService {
	return &BookingService{
		orderRepo:        orderRepo,
		availabilityRepo: availabilityRepo,
	}
}

// BookRoom бронирует номер.
func (s *BookingService) BookRoom(orderDto http.OrderDto) error {
	order := entity.Order{
		HotelID:   orderDto.HotelID,
		RoomID:    orderDto.RoomID,
		UserEmail: orderDto.UserEmail,
		From:      orderDto.From,
		To:        orderDto.To,
	}

	// имитация транзакции при проверке доступности номера и последующего бронирования, чтобы избежать овербукинга
	s.lock.Lock()
	defer s.lock.Unlock()

	// проверяем доступен ли номер на заданные даты
	ok := s.isRoomAvailableForDates(order.HotelID, order.RoomID, order.From, order.To)

	if !ok {
		return fmt.Errorf("no available dates")
	}

	// создаем бронь
	err := s.orderRepo.CreateOrder(order)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// резервируем номер на выбранные даты
	s.availabilityRepo.OccupyDates(order.HotelID, order.RoomID, order.From, order.To)

	return nil
}

// isRoomAvailableForDates определяет доступен ли номер для бронирования для заданного периода.
func (s *BookingService) isRoomAvailableForDates(hotelID, roomID string, from, to time.Time) bool {
	availableDates := s.availabilityRepo.GetAvailableDates(hotelID, roomID)
	requiredDates := utils.PeriodToDateList(from, to)

	availableDatesMap := make(map[time.Time]struct{})
	for _, d := range availableDates {
		availableDatesMap[d] = struct{}{}
	}

	for _, d := range requiredDates {
		if _, ok := availableDatesMap[d]; !ok {
			return false
		}
	}

	return true
}
