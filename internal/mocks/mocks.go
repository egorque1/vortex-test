package mocks

import (
	"github.com/egorque1/vortex-test/internal/entity"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) GetOrderBook(exchangeName, pair string) ([]*entity.OrderBook, error) {
	args := m.Called(exchangeName, pair)
	return args.Get(0).([]*entity.OrderBook), args.Error(1)
}

func (m *MockOrderRepository) SaveOrderBook(books []*entity.OrderBook) error {
	args := m.Called(books)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrderHistory(client *entity.Client) ([]*entity.HistoryOrder, error) {
	args := m.Called(client)
	return args.Get(0).([]*entity.HistoryOrder), args.Error(1)
}

func (m *MockOrderRepository) SaveOrderHistory(order entity.HistoryOrder) error {
	args := m.Called(order)
	return args.Error(0)
}

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) GetOrderBook(exchange, pair string) ([]*entity.OrderBook, error) {
	args := m.Called(exchange, pair)
	return args.Get(0).([]*entity.OrderBook), args.Error(1)
}

func (m *MockOrderService) SaveOrderBook(books []*entity.OrderBook) error {
	args := m.Called(books)
	return args.Error(0)
}

func (m *MockOrderService) GetOrderHistory(client *entity.Client) ([]*entity.HistoryOrder, error) {
	args := m.Called(client)
	return args.Get(0).([]*entity.HistoryOrder), args.Error(1)
}

func (m *MockOrderService) SaveOrderHistory(order entity.HistoryOrder) error {
	args := m.Called(order)
	return args.Error(0)
}
