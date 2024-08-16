package repository

import (
	"sync"

	"hotel/internal/entity"
)

type OrderRepository struct {
	orders []entity.Order

	rwLock sync.RWMutex
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) CreateOrder(order entity.Order) error {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	r.orders = append(r.orders, order)
	return nil
}
