package interfaces

import "challenge-yuno/internal/business/domain/order"

type INotificationService interface {
	SendNotification(order *order.Order)
}
