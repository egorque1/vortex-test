package repository

import (
	"fmt"

	"github.com/egorque1/vortex-test/internal/entity"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetOrderBook(exchange_name, pair string) ([]*entity.OrderBook, error)
	SaveOrderBook(orderBook []*entity.OrderBook) error
	GetOrderHistory(client *entity.Client) ([]*entity.HistoryOrder, error)
	SaveOrderHistory(order entity.HistoryOrder) error
}

type orderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepositoryImpl{db: db}
}

/*
GetOrderBook retrieves order books from the database for a specified exchange and trading pair.
It converts the retrieved order book DTOs to entities and returns them.
If no records are found, it returns a "record not found" error.
If a database error occurs, it returns the error.
*/

func (r *orderRepositoryImpl) GetOrderBook(exchange_name, pair string) ([]*entity.OrderBook, error) {
	var orderBookDTOs []*entity.OrderBookDTO
	tx := r.db.Where("exchange = ?", exchange_name).
		Where("pair = ?", pair).
		Find(&orderBookDTOs)

	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if tx.Error != nil {
		return nil, tx.Error
	}

	var orderBooks []*entity.OrderBook
	for _, orderBookDTO := range orderBookDTOs {
		orderBook, err := entity.ToOrderBookEntity(orderBookDTO)
		if err != nil {
			return nil, fmt.Errorf("error converting to Entity: %w", err)
		}
		orderBooks = append(orderBooks, orderBook)
	}

	return orderBooks, nil
}

/*
SaveOrderBook saves an order book entity to the database.
It converts array of order books entities to DTOs and attempts to save them.
If the conversion or save operation fails, it returns an error.
*/

func (r *orderRepositoryImpl) SaveOrderBook(order []*entity.OrderBook) error {
	for _, o := range order {
		orderBookDTO, err := entity.ToOrderBookDTO(o)
		if err != nil {
			return fmt.Errorf("error converting to DTO: %w", err)
		}

		if err := r.db.Create(orderBookDTO).Error; err != nil {
			return fmt.Errorf("error saving OrderBook: %w", err)
		}
	}
	return nil
}

/*
GetOrderHistory retrieves the order history for a specific client from the database.
It filters the history based on the client's name, exchange name, label, and trading pair.
If no records are found, it returns a "record not found" error.
If a database error occurs, it returns the error.
*/

func (r *orderRepositoryImpl) GetOrderHistory(client *entity.Client) ([]*entity.HistoryOrder, error) {
	var orderHistory []*entity.HistoryOrder
	tx := r.db.
		Where("client_name = ?", client.ClientName).
		Where("exchange_name = ?", client.ExchangeName).
		Where("label = ?", client.Label).
		Where("pair = ?", client.Pair).
		Find(&orderHistory)

	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if tx.Error != nil {
		return nil, tx.Error
	}

	return orderHistory, nil
}

/*
SaveOrderHistory saves a history order entity to the database.
If the save operation fails, it returns the error.
*/

func (r *orderRepositoryImpl) SaveOrderHistory(order entity.HistoryOrder) error {
	tx := r.db.Create(&order)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
