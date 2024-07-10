package service

import (
	"github.com/egorque1/vortex-test/internal/entity"
	"github.com/egorque1/vortex-test/internal/modules/repository"
)

type OrderService interface {
	GetOrderBook(exchange_name, pair string) ([]*entity.OrderBook, error)
	SaveOrderBook(orderBook []*entity.OrderBook) error
	GetOrderHistory(client *entity.Client) ([]*entity.HistoryOrder, error)
	SaveOrderHistory(order entity.HistoryOrder) error
}

type orderServiceImpl struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderServiceImpl{repo: repo}
}

/*
GetOrderBook returns the order book for a specific exchange and trading pair.
Also returns an error if one occures.
*/

func (s *orderServiceImpl) GetOrderBook(exchange_name, pair string) ([]*entity.OrderBook, error) {
	return s.repo.GetOrderBook(exchange_name, pair)
}

/*
SaveOrderBook saves the order book to ClickHouse.
Returns an error if one occures.
*/

func (s *orderServiceImpl) SaveOrderBook(orderBook []*entity.OrderBook) error {
	return s.repo.SaveOrderBook(orderBook)
}

/*
GetOrderHistory returns the order history for a specific client.
Also returns an error if one occures.
*/

func (s *orderServiceImpl) GetOrderHistory(client *entity.Client) ([]*entity.HistoryOrder, error) {
	return s.repo.GetOrderHistory(client)
}

/*
SaveOrderHistory saves the order history to ClickHouse.
Returns an error if one occures.
*/
func (s *orderServiceImpl) SaveOrderHistory(order entity.HistoryOrder) error {
	return s.repo.SaveOrderHistory(order)
}
