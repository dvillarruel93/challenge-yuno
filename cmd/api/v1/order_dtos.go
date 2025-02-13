package v1

import model "challenge-yuno/internal/business/domain/order"

type Order struct {
	Menu   []string         `json:"menu" validate:"required"`
	Status model.Status     `json:"status" validate:"required"`
	Source model.Source     `json:"source" validate:"required"`
	Type   *model.OrderType `json:"type,omitempty"`
}

func (o *Order) ToModel() model.Order {
	order := model.Order{
		Menu:   o.Menu,
		Status: o.Status,
		Source: o.Source,
		Type:   model.Normal,
	}

	if o.Type != nil {
		order.Type = *o.Type
	}

	return order
}

type OrderUpdate struct {
	Status   model.Status `json:"status" validate:"required"`
	Priority *int         `json:"priority,omitempty"`
}
