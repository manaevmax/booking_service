package http

import (
	"encoding/json"
	"net/http"
)

type BookingService interface {
	BookRoom(order OrderDto) error
}

// Handler обрабатывает запросы.
type Handler struct {
	srv BookingService
}

// New возвращает новый экземпляр обработчика.
func New(srv BookingService) *Handler {
	return &Handler{
		srv: srv,
	}
}

// CreateOrder создает новую бронь.
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order OrderDto

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := order.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.srv.BookRoom(order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
