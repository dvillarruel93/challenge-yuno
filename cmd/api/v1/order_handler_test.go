package v1

import (
	"bytes"
	"challenge-yuno/internal/business/domain/order"
	"challenge-yuno/internal/mocks"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type OrderHandlerTestSuite struct {
	suite.Suite
	orderHandler *OrderHandler
	orderUseCase *mocks.MockOrderUsecase
}

func (s *OrderHandlerTestSuite) SetupTest() {
	s.orderUseCase = new(mocks.MockOrderUsecase)
	s.orderHandler = &OrderHandler{s.orderUseCase}
}

func TestOrderHandler(t *testing.T) {
	suite.Run(t, new(OrderHandlerTestSuite))
}

func (s *OrderHandlerTestSuite) TestAddOrder() {
	var tests = []struct {
		name                 string
		payload              []byte
		mockExpectedResponse *order.Order
		mockExpectedError    error
		expectedResponse     *order.Order
		expectedError        error
	}{
		{
			name:                 "error_wrong_payload",
			payload:              []byte(`{bad payload!}`),
			mockExpectedResponse: &order.Order{ID: "123456"},
			mockExpectedError:    nil,
			expectedResponse:     nil,
			expectedError:        echo.NewHTTPError(http.StatusBadRequest, "error binding order body"),
		},
		{
			name:                 "error_validating_payload",
			payload:              []byte(`{"menu": ["food"]}`),
			mockExpectedResponse: &order.Order{ID: "123456"},
			mockExpectedError:    nil,
			expectedResponse:     nil,
			expectedError:        echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error validating model: %s", "Key: 'Order.Status' Error:Field validation for 'Status' failed on the 'required' tag\nKey: 'Order.Source' Error:Field validation for 'Source' failed on the 'required' tag")),
		},
		{
			name:                 "error_adding_order",
			payload:              []byte(`{"menu": ["drink"], "status": "DELIVERED", "source": "IN_PERSON", "number": 1}`),
			mockExpectedResponse: &order.Order{Menu: []string{"drink"}, Status: order.Delivered, Source: order.InPerson, Type: order.Normal},
			mockExpectedError:    fmt.Errorf("mock error"),
			expectedResponse:     nil,
			expectedError:        fmt.Errorf("mock error"),
		},
		{
			name:                 "success",
			payload:              []byte(`{"menu": ["food"], "status": "DELIVERED", "source": "IN_PERSON", "number": 1}`),
			mockExpectedResponse: &order.Order{Menu: []string{"food"}, Status: order.Delivered, Source: order.InPerson, Type: order.Normal},
			mockExpectedError:    nil,
			expectedResponse:     &order.Order{Menu: []string{"food"}, Status: order.Delivered, Source: order.InPerson, Type: order.Normal},
			expectedError:        nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req, err := http.NewRequest(http.MethodPost, "/order", bytes.NewReader(tt.payload))
			s.Require().NoError(err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			recorder := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, recorder)

			s.orderUseCase.On("AddOrder", *tt.mockExpectedResponse).
				Return(tt.mockExpectedResponse, tt.mockExpectedError)

			err = s.orderHandler.AddOrder(ctx)

			if tt.expectedError != nil {
				s.Require().Error(err)
				s.Equal(tt.expectedError, err)
				return
			}

			s.Require().NoError(err)
			response := &order.Order{}
			err = json.Unmarshal(recorder.Body.Bytes(), response)
			s.Require().NoError(err)
			s.Equal(tt.expectedResponse, response)
		})
	}
}

func (s *OrderHandlerTestSuite) TestGetOrder() {
	req, err := http.NewRequest(http.MethodGet, "/order", nil)
	s.Require().NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	var tests = []struct {
		name                 string
		orderID              string
		mockExpectedResponse *order.Order
		mockExpectedError    error
		expectedResponse     *order.Order
		expectedError        error
	}{
		{
			name:                 "error_empty_param",
			orderID:              "",
			mockExpectedResponse: &order.Order{ID: "123456"},
			mockExpectedError:    nil,
			expectedResponse:     nil,
			expectedError:        echo.NewHTTPError(http.StatusBadRequest, "ID param can't be empty"),
		},
		{
			name:                 "error_getting_order",
			orderID:              "123789",
			mockExpectedResponse: nil,
			mockExpectedError:    fmt.Errorf("mock error"),
			expectedResponse:     nil,
			expectedError:        fmt.Errorf("mock error"),
		},
		{
			name:                 "success",
			orderID:              "123456",
			mockExpectedResponse: &order.Order{ID: "123456"},
			mockExpectedError:    nil,
			expectedResponse:     &order.Order{ID: "123456"},
			expectedError:        nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			recorder := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, recorder)
			ctx.SetParamNames("ID")
			ctx.SetParamValues(tt.orderID)

			s.orderUseCase.On("GetOrder", tt.orderID).
				Return(tt.mockExpectedResponse, tt.mockExpectedError)

			err = s.orderHandler.GetOrder(ctx)

			if tt.expectedError != nil {
				s.Require().Error(err)
				s.Equal(tt.expectedError, err)
				return
			}

			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, recorder.Code)
			response := &order.Order{}
			err = json.Unmarshal(recorder.Body.Bytes(), response)
			s.Require().NoError(err)
			s.Equal(tt.expectedResponse, response)
		})
	}
}

func (s *OrderHandlerTestSuite) TestListActiveOrders() {
	req, err := http.NewRequest(http.MethodGet, "/order/active", nil)
	s.Require().NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	var tests = []struct {
		name                 string
		mockExpectedResponse []order.Order
		mockExpectedError    error
		expectedResponse     []order.Order
		expectedError        error
	}{
		{
			name:                 "error_getting_orders",
			mockExpectedResponse: nil,
			mockExpectedError:    fmt.Errorf("mock error"),
			expectedResponse:     nil,
			expectedError:        fmt.Errorf("mock error"),
		},
		{
			name:                 "success",
			mockExpectedResponse: []order.Order{{ID: "123456"}},
			mockExpectedError:    nil,
			expectedResponse:     []order.Order{{ID: "123456"}},
			expectedError:        nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			recorder := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, recorder)

			s.orderUseCase.On("ListActiveOrders").
				Return(tt.mockExpectedResponse, tt.mockExpectedError).Once()

			err = s.orderHandler.ListActiveOrders(ctx)

			if tt.expectedError != nil {
				s.Require().Error(err)
				s.Equal(tt.expectedError, err)
				return
			}

			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, recorder.Code)
			var response []order.Order
			err = json.Unmarshal(recorder.Body.Bytes(), &response)
			s.Require().NoError(err)
			s.Equal(tt.expectedResponse, response)
		})
	}
}

func (s *OrderHandlerTestSuite) TestCancelOrder() {
	req, err := http.NewRequest(http.MethodPut, "/order", nil)
	s.Require().NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	var tests = []struct {
		name                 string
		orderID              string
		mockExpectedResponse *order.Order
		mockExpectedError    error
		expectedResponse     *order.Order
		expectedError        error
	}{
		{
			name:                 "error_empty_param",
			orderID:              "",
			mockExpectedResponse: &order.Order{ID: "123456"},
			mockExpectedError:    nil,
			expectedResponse:     nil,
			expectedError:        echo.NewHTTPError(http.StatusBadRequest, "ID param can't be empty"),
		},
		{
			name:                 "error_getting_order",
			orderID:              "123789",
			mockExpectedResponse: nil,
			mockExpectedError:    fmt.Errorf("mock error"),
			expectedResponse:     nil,
			expectedError:        fmt.Errorf("mock error"),
		},
		{
			name:                 "success",
			orderID:              "123456",
			mockExpectedResponse: &order.Order{ID: "123456"},
			mockExpectedError:    nil,
			expectedResponse:     &order.Order{ID: "123456"},
			expectedError:        nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			recorder := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, recorder)
			ctx.SetParamNames("ID")
			ctx.SetParamValues(tt.orderID)

			s.orderUseCase.On("UpdateOrder", tt.orderID, order.Canceled, mock.Anything).
				Return(tt.mockExpectedResponse, tt.mockExpectedError)

			err = s.orderHandler.CancelOrder(ctx)

			if tt.expectedError != nil {
				s.Require().Error(err)
				s.Equal(tt.expectedError, err)
				return
			}

			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, recorder.Code)
			response := &order.Order{}
			err = json.Unmarshal(recorder.Body.Bytes(), response)
			s.Require().NoError(err)
			s.Equal(tt.expectedResponse, response)
		})
	}
}

func (s *OrderHandlerTestSuite) TestUpdateOrder() {
	var tests = []struct {
		name                 string
		orderID              string
		payload              []byte
		mockExpectedResponse *order.Order
		mockExpectedError    error
		expectedResponse     *order.Order
		expectedError        error
	}{
		{
			name:                 "error_wrong_payload",
			orderID:              "",
			payload:              []byte(`{bad payload!}`),
			mockExpectedResponse: &order.Order{ID: "123456"},
			mockExpectedError:    nil,
			expectedResponse:     nil,
			expectedError:        echo.NewHTTPError(http.StatusBadRequest, "error binding order body"),
		},
		{
			name:                 "error_validating_payload",
			orderID:              "123789",
			payload:              []byte(`{"test": "bad payload"}`),
			mockExpectedResponse: nil,
			mockExpectedError:    nil,
			expectedResponse:     nil,
			expectedError:        echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error validating model: %s", "Key: 'OrderUpdate.Status' Error:Field validation for 'Status' failed on the 'required' tag")),
		},
		{
			name:                 "error_empty_param",
			orderID:              "",
			payload:              []byte(`{"status": "DELIVERED"}`),
			mockExpectedResponse: &order.Order{ID: "123456"},
			mockExpectedError:    nil,
			expectedResponse:     nil,
			expectedError:        echo.NewHTTPError(http.StatusBadRequest, "ID param can't be empty"),
		},
		{
			name:                 "success",
			orderID:              "123456",
			payload:              []byte(`{"status": "DELIVERED"}`),
			mockExpectedResponse: &order.Order{ID: "123456", Status: order.Delivered},
			mockExpectedError:    nil,
			expectedResponse:     &order.Order{ID: "123456", Status: order.Delivered},
			expectedError:        nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req, err := http.NewRequest(http.MethodPut, "/order", bytes.NewReader(tt.payload))
			s.Require().NoError(err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			recorder := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, recorder)
			ctx.SetParamNames("ID")
			ctx.SetParamValues(tt.orderID)

			s.orderUseCase.On("UpdateOrder", tt.orderID, order.Delivered, mock.Anything).
				Return(tt.mockExpectedResponse, tt.mockExpectedError)

			err = s.orderHandler.UpdateOrder(ctx)

			if tt.expectedError != nil {
				s.Require().Error(err)
				s.Equal(tt.expectedError, err)
				return
			}

			s.Require().NoError(err)
			response := &order.Order{}
			err = json.Unmarshal(recorder.Body.Bytes(), response)
			s.Require().NoError(err)
			s.Equal(tt.expectedResponse, response)
		})
	}
}
