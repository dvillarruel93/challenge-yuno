package interfaces

import model "challenge-yuno/internal/business/domain/order"

type OrderUsecase interface {
	AddOrder(order model.Order) (*model.Order, error)
	GetOrder(orderID string) (*model.Order, error)
	ListActiveOrders() ([]model.Order, error)
	UpdateOrder(orderID string, status model.Status, priority *int) (*model.Order, error)
	GetAllOrders() ([]model.Order, error)
}
