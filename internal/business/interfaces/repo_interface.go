package interfaces

import model "challenge-yuno/internal/business/domain/order"

type KVSOrderRepository interface {
	AddOrder(order model.Order) (*model.Order, error)
	GetOrder(orderID string) (*model.Order, error)
	ListActiveOrders() ([]model.Order, error)
	UpdateOrderStatus(orderID string, status model.Status) (*model.Order, error)
	GetAllOrders() []model.Order
}

type SQLOrderRepository interface {
	AddOrder(order model.Order) (*model.Order, error)
	GetOrder(orderID string) (*model.Order, error)
	ListActiveOrders() ([]model.Order, error)
	UpdateOrder(orderID string, status model.Status, priority *int) (*model.Order, error)
	GetAllOrders() ([]model.Order, error)
}
