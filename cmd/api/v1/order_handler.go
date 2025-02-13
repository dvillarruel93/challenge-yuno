package v1

import (
	model "challenge-yuno/internal/business/domain/order"
	"challenge-yuno/internal/business/interfaces"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"sync"
)

type OrderHandler struct {
	OrderUsecase interfaces.OrderUsecase
}

func NewOrderHandler(e *echo.Echo, orderUsecase interfaces.OrderUsecase) {
	handler := &OrderHandler{
		OrderUsecase: orderUsecase,
	}

	e.POST("/order", handler.AddOrder)
	e.GET("/order/active", handler.ListActiveOrders)
	e.GET("/order/:ID", handler.GetOrder)
	e.PUT("/order/:ID/cancel", handler.CancelOrder)
	e.PUT("/order/:ID/status", handler.UpdateOrder)

	e.POST("/order/test", handler.TestOrders)
	e.GET("/order/all", handler.GetAllOrders)
}

func (h *OrderHandler) AddOrder(c echo.Context) error {
	order := Order{}
	if err := c.Bind(&order); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error binding order body")
	}

	if err := model.Validate(order); err != nil {
		return err
	}

	response, err := h.OrderUsecase.AddOrder(order.ToModel())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *OrderHandler) GetOrder(c echo.Context) error {
	orderID := c.Param("ID")
	if len(orderID) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "ID param can't be empty")
	}

	response, err := h.OrderUsecase.GetOrder(orderID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) ListActiveOrders(c echo.Context) error {

	response, err := h.OrderUsecase.ListActiveOrders()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) CancelOrder(c echo.Context) error {
	orderID := c.Param("ID")
	if len(orderID) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "ID param can't be empty")
	}

	response, err := h.OrderUsecase.UpdateOrder(orderID, model.Canceled, nil)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) UpdateOrder(c echo.Context) error {
	order := OrderUpdate{}
	if err := c.Bind(&order); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error binding order body")
	}

	if err := model.Validate(order); err != nil {
		return err
	}

	orderID := c.Param("ID")
	if len(orderID) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "ID param can't be empty")
	}

	response, err := h.OrderUsecase.UpdateOrder(orderID, order.Status, order.Priority)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *OrderHandler) TestOrders(c echo.Context) error {
	var wg sync.WaitGroup
	sources := []model.Source{model.InPerson, model.Phone, model.Delivery}
	statuses := []model.Status{model.Pending, model.InPreparation, model.Finished, model.Delivered, model.Canceled}

	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			order := model.Order{
				Menu:   []string{fmt.Sprintf("Plato # %d", i)},
				Status: statuses[i%len(statuses)],
				Source: sources[i%len(sources)],
				Type:   model.Normal,
			}

			_, err := h.OrderUsecase.AddOrder(order)
			if err != nil {
				log.Errorf("error saving order: %w", err)
			}
		}(i)
	}

	wg.Wait()

	return c.JSON(http.StatusCreated, "all orders created")
}

func (h *OrderHandler) GetAllOrders(c echo.Context) error {
	result, err := h.OrderUsecase.GetAllOrders()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}
