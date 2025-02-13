package kvstore

import (
	domain "challenge-yuno/internal/business/domain/order"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	orderRepo *OrderRepository
}

func (s *OrderRepositoryTestSuite) SetupTest() {
	s.orderRepo = NewOrderRepository()
}

func TestOrderRepository(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (s *OrderRepositoryTestSuite) TestAddOrder() {
	order := domain.Order{
		Menu:   []string{"food", "drink"},
		Status: domain.InPreparation,
		Source: domain.Delivery,
		Type:   domain.Normal,
	}

	response, err := s.orderRepo.AddOrder(order)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().NotNil(response.ID)
	s.Require().NotNil(response.CreatedAt)
}

func (s *OrderRepositoryTestSuite) TestGetOrder() {
	response, err := s.orderRepo.GetOrder("test-id")
	s.Require().Nil(response)
	s.Require().Error(err)
	s.Require().Equal(echo.NewHTTPError(http.StatusNotFound, "order not found"), err)

	order := domain.Order{
		Menu:   []string{"food", "drink"},
		Status: domain.Pending,
		Source: domain.Delivery,
		Type:   domain.Normal,
	}
	response, err = s.orderRepo.AddOrder(order)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().NotNil(response.ID)
	s.Require().NotNil(response.CreatedAt)

	getResponse, err := s.orderRepo.GetOrder(response.ID)
	s.Require().NoError(err)
	s.Require().NotNil(getResponse)
	s.Require().Equal(response, getResponse)
}

func (s *OrderRepositoryTestSuite) TestListActiveOrders() {
	listOrders, err := s.orderRepo.ListActiveOrders()
	s.Require().Nil(listOrders)
	s.Require().Error(err)
	s.Require().Equal(echo.NewHTTPError(http.StatusNotFound, "orders not found"), err)

	order := domain.Order{
		Menu:   []string{"food", "drink"},
		Status: domain.Pending,
		Source: domain.Delivery,
		Type:   domain.Normal,
	}

	response, err := s.orderRepo.AddOrder(order)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().NotNil(response.ID)
	s.Require().NotNil(response.CreatedAt)

	listOrders, err = s.orderRepo.ListActiveOrders()
	s.Require().NoError(err)
	s.Require().Equal(1, len(listOrders))
	s.Require().Equal(*response, listOrders[0])
}

func (s *OrderRepositoryTestSuite) TestUpdateOrderStatus() {
	orderUpdated, err := s.orderRepo.UpdateOrderStatus("some-id", domain.InPreparation)
	s.Require().Nil(orderUpdated)
	s.Require().Error(err)
	s.Require().Equal(echo.NewHTTPError(http.StatusNotFound, "order not found"), err)

	order := domain.Order{
		Menu:   []string{"food", "drink"},
		Status: domain.Pending,
		Source: domain.Delivery,
		Type:   domain.Normal,
	}

	response, err := s.orderRepo.AddOrder(order)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().NotNil(response.ID)
	s.Require().NotNil(response.CreatedAt)

	orderUpdated, err = s.orderRepo.UpdateOrderStatus(response.ID, domain.InPreparation)
	s.Require().NoError(err)
	s.Require().NotNil(orderUpdated)
	s.Require().Equal(response.ID, orderUpdated.ID)
	s.Require().Equal(domain.InPreparation, orderUpdated.Status)
}
