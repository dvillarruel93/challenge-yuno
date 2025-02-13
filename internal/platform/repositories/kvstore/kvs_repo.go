package kvstore

import (
	domain "challenge-yuno/internal/business/domain/order"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"sync"
	"time"
)

type OrderRepository struct {
	indexMap map[string]int
	orders   []orderDB
	mu       sync.Mutex
}

type orderDB struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Menu      []string
	Status    string
	Source    string
	Type      string
}

func NewOrderRepository() *OrderRepository {
	im := make(map[string]int)
	return &OrderRepository{
		indexMap: im,
		orders:   []orderDB{},
	}
}

func toOrderDB(o domain.Order) orderDB {
	now := time.Now().Truncate(time.Millisecond)
	return orderDB{
		ID:        uuid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
		Menu:      o.Menu,
		Status:    string(o.Status),
		Source:    string(o.Source),
		Type:      string(o.Type),
	}
}

func (o *orderDB) toOrderModel() *domain.Order {
	return &domain.Order{
		ID:        o.ID,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		Menu:      o.Menu,
		Status:    domain.Status(o.Status),
		Source:    domain.Source(o.Source),
		Type:      domain.OrderType(o.Type),
	}
}

func (r *OrderRepository) AddOrder(order domain.Order) (*domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	oDB := toOrderDB(order)

	if _, exists := r.indexMap[oDB.ID]; exists {
		log.Errorf("order %s already exists", oDB.ID)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "order already added")
	}

	r.orders = append(r.orders, oDB)
	r.indexMap[oDB.ID] = len(r.orders) - 1

	return oDB.toOrderModel(), nil
}

func (r *OrderRepository) GetOrder(orderID string) (*domain.Order, error) {
	var index int
	var exists bool

	if index, exists = r.indexMap[orderID]; !exists {
		log.Errorf("there is no order in db with id %s", orderID)
		return nil, echo.NewHTTPError(http.StatusNotFound, "order not found")
	}

	return r.orders[index].toOrderModel(), nil
}

func (r *OrderRepository) ListActiveOrders() ([]domain.Order, error) {
	var result []domain.Order

	for _, val := range r.orders {
		if val.Status == string(domain.Pending) {
			result = append(result, *val.toOrderModel())
		}
	}

	if len(result) == 0 {
		log.Errorf("there is no pending orders")
		return nil, echo.NewHTTPError(http.StatusNotFound, "orders not found")
	}

	return result, nil
}

func (r *OrderRepository) UpdateOrderStatus(orderID string, status domain.Status) (*domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var index int
	var exists bool

	if index, exists = r.indexMap[orderID]; !exists {
		log.Errorf("there is no order in db with id %s", orderID)
		return nil, echo.NewHTTPError(http.StatusNotFound, "order not found")
	}

	r.orders[index].Status = string(status)
	r.orders[index].UpdatedAt = time.Now().Truncate(time.Millisecond)

	return r.orders[index].toOrderModel(), nil
}

func (r *OrderRepository) GetAllOrders() []domain.Order {
	r.mu.Lock()
	defer r.mu.Unlock()

	var orders []domain.Order

	for _, o := range r.orders {
		orders = append(orders, *o.toOrderModel())
	}

	return orders
}
