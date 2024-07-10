package entity

import (
	"encoding/json"
	"fmt"
)

type OrderBook struct {
	ID       int64
	Exchange string
	Pair     string
	Asks     []DepthOrder
	Bids     []DepthOrder
}

type OrderBookDTO struct {
	ID       int64  `json:"id"`
	Exchange string `json:"exchange"`
	Pair     string `json:"pair"`
	Asks     string `json:"asks"`
	Bids     string `json:"bids"`
}

func ToOrderBookDTO(orderBook *OrderBook) (*OrderBookDTO, error) {
	asksJSON, err := json.Marshal(orderBook.Asks)
	if err != nil {
		return nil, fmt.Errorf("error marshalling asks to JSON: %w", err)
	}

	bidsJSON, err := json.Marshal(orderBook.Bids)
	if err != nil {
		return nil, fmt.Errorf("error marshalling bids to JSON: %w", err)
	}

	return &OrderBookDTO{
		ID:       orderBook.ID,
		Exchange: orderBook.Exchange,
		Pair:     orderBook.Pair,
		Asks:     string(asksJSON), // Строковое представление JSON данных
		Bids:     string(bidsJSON), // Строковое представление JSON данных
	}, nil
}

func ToOrderBookEntity(dto *OrderBookDTO) (*OrderBook, error) {
	var asks []DepthOrder
	if err := json.Unmarshal([]byte(dto.Asks), &asks); err != nil {
		return nil, err
	}

	var bids []DepthOrder
	if err := json.Unmarshal([]byte(dto.Bids), &bids); err != nil {
		return nil, err
	}

	return &OrderBook{
		ID:       dto.ID,
		Exchange: dto.Exchange,
		Pair:     dto.Pair,
		Asks:     asks,
		Bids:     bids,
	}, nil
}
