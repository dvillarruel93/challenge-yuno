package services

import (
	"challenge-yuno/internal/business/domain/order"
	"github.com/labstack/gommon/log"
)

type NotificationService struct {
	notificationClient string
}

func NewNotificationService(notificationClient string) *NotificationService {
	return &NotificationService{
		notificationClient: notificationClient,
	}
}

func (n *NotificationService) SendNotification(order *order.Order) {
	log.Infof("send notification of order %s through client %s \n", order.ID, n.notificationClient)
}
