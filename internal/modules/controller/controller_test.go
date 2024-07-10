package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/egorque1/vortex-test/internal/entity"
	"github.com/egorque1/vortex-test/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetOrderBookHandler(t *testing.T) {
	mockService := &mocks.MockOrderService{}
	mockRepo := &mocks.MockOrderRepository{}
	controller := NewController(mockRepo, mockService)

	reqBody := &entity.OrderBookRequest{
		Exchange_name: "Binance",
		Pair:          "BTC/USDT",
	}
	reqBytes, _ := json.Marshal(reqBody)

	mockOrderBook := []*entity.OrderBook{
		{
			ID:       1,
			Exchange: "Binance",
			Pair:     "BTC/USDT",
			Asks:     []entity.DepthOrder{{Price: 30000, BaseQty: 0.5}, {Price: 30010, BaseQty: 0.2}},
			Bids:     []entity.DepthOrder{{Price: 29900, BaseQty: 0.5}, {Price: 29850, BaseQty: 1}},
		},
	}
	mockService.On("GetOrderBook", "Binance", "BTC/USDT").Return(mockOrderBook, nil)

	req, err := http.NewRequest("POST", "/order-book", bytes.NewBuffer(reqBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(controller.GetOrderBookHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedBody, _ := json.Marshal(mockOrderBook)
	assert.Equal(t, expectedBody, rr.Body.Bytes())

	mockService.AssertCalled(t, "GetOrderBook", "Binance", "BTC/USDT")
}

func TestGetOrderHandler_BadRequest(t *testing.T) {
	mockService := &mocks.MockOrderService{}
	mockRepo := &mocks.MockOrderRepository{}
	controller := NewController(mockRepo, mockService)

	req := httptest.NewRequest("GET", "/order/book?exchange=BYBIT&pair=ETH-BTC", nil)
	rr := httptest.NewRecorder()

	controller.GetOrderBookHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSaveOrderBookHandler_BadRequest(t *testing.T) {
	mockService := &mocks.MockOrderService{}
	mockRepo := &mocks.MockOrderRepository{}
	controller := NewController(mockRepo, mockService)

	reqBody := `invalid json`
	reqBytes := []byte(reqBody)

	req, err := http.NewRequest("POST", "/order-books", bytes.NewBuffer(reqBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(controller.SaveOrderBookHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetOrderHistoryHandler(t *testing.T) {
	mockService := &mocks.MockOrderService{}
	mockRepo := &mocks.MockOrderRepository{}
	controller := NewController(mockRepo, mockService)

	reqBody := []byte(`{"client_name": "client1"}`)
	req := httptest.NewRequest("POST", "/getOrderHistory", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mockService.On("GetOrderHistory", mock.Anything).Return([]*entity.HistoryOrder{}, nil)

	http.HandlerFunc(controller.GetOrderHistoryHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %d", rr.Code)
	}

	expectedBody := `[]`
	if rr.Body.String() != expectedBody {
		t.Errorf("expected body %s; got %s", expectedBody, rr.Body.String())
	}

	mockService.AssertCalled(t, "GetOrderHistory", mock.Anything)
}

func TestGetOrderHistoryHandler_BadRequest(t *testing.T) {
	reqBody := []byte(`invalid-json`)

	req := httptest.NewRequest("POST", "/getOrderHistory", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	controller := NewController(nil, nil)

	http.HandlerFunc(controller.GetOrderHistoryHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest; got %d", rr.Code)
	}
}

func TestSaveOrderHistoryHandler(t *testing.T) {
	mockService := &mocks.MockOrderService{}
	mockRepo := &mocks.MockOrderRepository{}
	controller := NewController(mockRepo, mockService)

	order := entity.HistoryOrder{
		ClientName: "client1",
	}
	reqBody, _ := json.Marshal(order)
	req := httptest.NewRequest("POST", "/saveOrderHistory", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mockService.On("SaveOrderHistory", order).Return(nil)

	http.HandlerFunc(controller.SaveOrderHistoryHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %d", rr.Code)
	}

	mockService.AssertCalled(t, "SaveOrderHistory", order)
}

func TestSaveOrderHistoryHandler_BadRequest(t *testing.T) {
	reqBody := []byte(`invalid-json`)

	req := httptest.NewRequest("POST", "/saveOrderHistory", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	controller := NewController(nil, nil)

	http.HandlerFunc(controller.SaveOrderHistoryHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest; got %d", rr.Code)
	}
}
