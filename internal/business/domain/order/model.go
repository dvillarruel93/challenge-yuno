package order

import "time"

type Order struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Menu      []string  `json:"menu"`
	Status    Status    `json:"status"`
	Source    Source    `json:"order_source"`
	Type      OrderType `json:"order_type"`
	Priority  int       `json:"priority"`
}

type Status string

const (
	Pending       Status = "PENDING"
	InPreparation Status = "IN_PREPARATION"
	Finished      Status = "FINISHED"
	Delivered     Status = "DELIVERED"
	Canceled      Status = "CANCELED"
)

type Source string

const (
	InPerson Source = "IN_PERSON"
	Delivery Source = "DELIVERY"
	Phone    Source = "PHONE"
)

type OrderType string

const (
	Normal OrderType = "NORMAL"
	VIP    OrderType = "VIP"
)
