package sql

import (
	domain "challenge-yuno/internal/business/domain/order"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"sync"
	"time"
)

type OrderRepository struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	db.Exec("DROP TABLE IF EXISTS orders")
	db.Exec("DROP TABLE IF EXISTS order_dbs")
	db.AutoMigrate(&orders{})
	db.AutoMigrate(&orderDB{})
	return &OrderRepository{
		db: db,
	}
}

type orders struct {
	ID        string    `json:"id" gorm:"type:string;size:255;primary_key;"`
	CreatedAt time.Time `json:"created_at" gorm:"<-:create;type:time;not null;"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:time; not null"`
	Menu      []string  `json:"menu" gorm:"type:text[]; not null;"`
	Status    string    `json:"status" gorm:"type:string; size:255; not null;"`
	Source    string    `json:"order_source" gorm:"type:string; size:255; not null;"`
	Type      string    `json:"order_type" gorm:"type:string; size:255; not null;"`
}

type orderDB struct {
	ID        string    `json:"id" gorm:"type:string; size:255; primary_key;"`
	CreatedAt time.Time `json:"created_at" gorm:"<-:create; type:time; not null;"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:time; not null"`
	Menu      string    `json:"menu" gorm:"type:string; size:255; not null;"`
	Status    string    `json:"status" gorm:"type:string; size:255; not null;"`
	Source    string    `json:"order_source" gorm:"type:string; size:255; not null;"`
	Type      string    `json:"order_type" gorm:"type:string; size:255; not null;"`
	Priority  int       `json:"priority" gorm:"type:integer;not null;default:0"`
}

func toOrderDB(o domain.Order) orders {
	now := time.Now().Truncate(time.Millisecond)
	return orders{
		ID:        uuid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
		Menu:      o.Menu,
		Status:    string(o.Status),
		Source:    string(o.Source),
		Type:      string(o.Type),
	}
}

func toOrderDB2(o domain.Order, priority int) orderDB {
	now := time.Now().Truncate(time.Millisecond)
	return orderDB{
		ID:        uuid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
		Menu:      strings.Join(o.Menu, ","),
		Status:    string(o.Status),
		Source:    string(o.Source),
		Type:      string(o.Type),
		Priority:  priority + 1,
	}
}

func (o *orders) toOrderModel() *domain.Order {
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

func (o *orderDB) toOrderModel() *domain.Order {
	return &domain.Order{
		ID:        o.ID,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		Menu:      strings.Split(o.Menu, ","),
		Status:    domain.Status(o.Status),
		Source:    domain.Source(o.Source),
		Type:      domain.OrderType(o.Type),
		Priority:  o.Priority,
	}
}

func (r *OrderRepository) mapOrdersDBToOrdersModel(ordersDB []orderDB) []domain.Order {
	result := make([]domain.Order, 0, len(ordersDB))
	for _, oDB := range ordersDB {
		result = append(result, *oDB.toOrderModel())
	}

	return result
}

func (r *OrderRepository) AddOrder(order domain.Order) (*domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	priority, err := r.getDailyPriority()
	if err != nil {
		return nil, err
	}
	oDB := toOrderDB2(order, *priority)
	err = r.db.Create(&oDB).Error
	if err != nil {
		log.Errorf("error saving order: %v", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "order wasn't created")
	}

	return oDB.toOrderModel(), nil
}

func (r *OrderRepository) GetOrder(orderID string) (*domain.Order, error) {
	var oDB orderDB

	err := r.db.Model(&orderDB{}).First(&oDB, "id = ?", orderID).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, echo.NewHTTPError(http.StatusNotFound, "order not found")
		}
		log.Errorf("error getting order: %v", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "error getting order")
	}

	return oDB.toOrderModel(), nil
}

func (r *OrderRepository) ListActiveOrders() ([]domain.Order, error) {
	var ordersDB []orderDB

	err := r.db.Where("status = ?", domain.Pending).
		Order("priority ASC").
		Order("created_at ASC").
		Find(&ordersDB).
		Error
	if err != nil {
		log.Errorf("error getting order: %v", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "error getting order")
	}
	if len(ordersDB) == 0 {
		return nil, echo.NewHTTPError(http.StatusNotFound, "there is no active orders")
	}

	return r.mapOrdersDBToOrdersModel(ordersDB), nil
}

func (r *OrderRepository) UpdateOrder(orderID string, status domain.Status, priority *int) (*domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var err error
	if priority == nil {
		err = r.db.Model(&orderDB{}).Where("id = ?", orderID).Update("status", status).Error
	} else {
		err = r.db.Model(&orderDB{}).
			Where("id = ?", orderID).
			Update("status", status).
			Update("priority", *priority).
			Error
	}
	if err != nil {
		log.Errorf("error updating order %s: %v", orderID, err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "error updating order")
	}

	return r.GetOrder(orderID)
}

func (r *OrderRepository) GetAllOrders() ([]domain.Order, error) {
	var ordersDB []orderDB

	err := r.db.
		Order("priority ASC").
		Order("created_at ASC").
		Find(&ordersDB).
		Error
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "error getting all order")
	}
	if len(ordersDB) == 0 {
		return nil, echo.NewHTTPError(http.StatusNotFound, "orders not found")
	}

	result := make([]domain.Order, 0, len(ordersDB))
	for _, oDB := range ordersDB {
		result = append(result, *oDB.toOrderModel())
	}

	return r.mapOrdersDBToOrdersModel(ordersDB), nil
}

func (r *OrderRepository) getDailyPriority() (*int, error) {
	var count int64

	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)

	// Contar las órdenes del día
	err := r.db.Model(&orderDB{}).
		Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).
		Count(&count).Error

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "error getting daily priority")
	}

	countResult := int(count)

	return &countResult, nil
}
