package service

import (
	"testing"
	"time"

	"github.com/egorque1/vortex-test/internal/entity"
	"github.com/egorque1/vortex-test/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetOrderBook(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	mockService := NewOrderService(mockRepo)

	exchange := "exchange1"
	pair := "pair1"
	expected := []*entity.OrderBook{
		{
			ID:       1,
			Exchange: exchange,
			Pair:     pair,
			Asks:     nil,
			Bids:     nil,
		},
	}

	mockRepo.On("GetOrderBook", exchange, pair).Return(expected, nil)

	result, err := mockService.GetOrderBook(exchange, pair)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestSaveOrderBook(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	mockService := NewOrderService(mockRepo)

	orderBook := []*entity.OrderBook{
		{
			ID:       1,
			Exchange: "exchange1",
			Pair:     "pair1",
			Asks:     nil,
			Bids:     nil,
		},
	}

	mockRepo.On("SaveOrderBook", orderBook).Return(nil)

	err := mockService.SaveOrderBook(orderBook)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestGetOrderHistory(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	mockService := NewOrderService(mockRepo)

	client := &entity.Client{
		ClientName:   "client1",
		ExchangeName: "exchange1",
		Label:        "label1",
		Pair:         "pair1",
	}

	expected := []*entity.HistoryOrder{
		{
			ClientName:          client.ClientName,
			ExchangeName:        client.ExchangeName,
			Label:               client.Label,
			Pair:                client.Pair,
			Side:                "buy",
			Type:                "market",
			BaseQty:             1.0,
			Price:               100.0,
			AlgorithmNamePlaced: "algo1",
			LowestSellPrc:       105.0,
			HighestBuyPrc:       95.0,
			CommissionQuoteQty:  1.0,
			TimePlaced:          time.Now(),
		},
	}

	mockRepo.On("GetOrderHistory", client).Return(expected, nil)

	result, err := mockService.GetOrderHistory(client)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestSaveOrderHistory(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	mockService := NewOrderService(mockRepo)

	order := entity.HistoryOrder{
		ClientName:          "client1",
		ExchangeName:        "exchange1",
		Label:               "label1",
		Pair:                "pair1",
		Side:                "sell",
		Type:                "limit",
		BaseQty:             1.0,
		Price:               200.0,
		AlgorithmNamePlaced: "algo2",
		LowestSellPrc:       205.0,
		HighestBuyPrc:       195.0,
		CommissionQuoteQty:  2.0,
		TimePlaced:          time.Now(),
	}

	mockRepo.On("SaveOrderHistory", order).Return(nil)

	err := mockService.SaveOrderHistory(order)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
