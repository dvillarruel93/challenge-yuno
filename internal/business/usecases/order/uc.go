package order

import (
	model "challenge-yuno/internal/business/domain/order"
	"challenge-yuno/internal/business/interfaces"
)

type OrderUsecase struct {
	KVSOrderRepository  interfaces.KVSOrderRepository
	SQLOrderRepository  interfaces.SQLOrderRepository
	NotificationService interfaces.INotificationService
}

func NewOrderUsecase(kvsOrderRepository interfaces.KVSOrderRepository,
	sqlOrderRepository interfaces.SQLOrderRepository, notifService interfaces.INotificationService) *OrderUsecase {
	return &OrderUsecase{
		KVSOrderRepository:  kvsOrderRepository,
		SQLOrderRepository:  sqlOrderRepository,
		NotificationService: notifService,
	}
}

func (u *OrderUsecase) AddOrder(order model.Order) (*model.Order, error) {
	return u.SQLOrderRepository.AddOrder(order)
}

func (u *OrderUsecase) GetOrder(orderID string) (*model.Order, error) {
	return u.SQLOrderRepository.GetOrder(orderID)
}

func (u *OrderUsecase) ListActiveOrders() ([]model.Order, error) {
	return u.SQLOrderRepository.ListActiveOrders()
}

func (u *OrderUsecase) UpdateOrder(orderID string, status model.Status, priority *int) (*model.Order, error) {
	order, err := u.SQLOrderRepository.UpdateOrder(orderID, status, priority)
	if err != nil {
		return nil, err
	}

	if order.Status == model.Finished {
		u.NotificationService.SendNotification(order)
	}

	return order, err
}

func (u *OrderUsecase) GetAllOrders() ([]model.Order, error) {
	return u.SQLOrderRepository.GetAllOrders()
}
