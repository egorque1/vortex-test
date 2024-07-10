package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/egorque1/vortex-test/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

func TestGetOrderBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT version()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("mock_version"))

	gormDB, err := gorm.Open(clickhouse.New(clickhouse.Config{DriverName: "clickhouse", Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("error creating gorm DB: %v", err)
	}

	repo := NewOrderRepository(gormDB)

	mock.ExpectQuery("^SELECT \\* FROM `order_book_dtos` WHERE exchange = \\? AND pair = \\?$").
		WithArgs("exchange1", "pair1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "exchange", "pair", "asks", "bids"}).
			AddRow(1, "exchange1", "pair1", `[]`, `[]`))

	// Test case: valid data
	orderBooks, err := repo.GetOrderBook("exchange1", "pair1")
	assert.NoError(t, err)
	assert.NotNil(t, orderBooks)
	assert.Len(t, orderBooks, 1)

	// Test case: record not found (using gorm.ErrRecordNotFound)
	mock.ExpectQuery("^SELECT \\* FROM `order_book_dtos` WHERE exchange = \\? AND pair = \\?$").
		WithArgs("nonexistent_exchange", "nonexistent_pair").
		WillReturnError(gorm.ErrRecordNotFound)
	orderBooks, err = repo.GetOrderBook("nonexistent_exchange", "nonexistent_pair")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Nil(t, orderBooks)

	// Test case: database error (using sql.ErrNoRows)
	mock.ExpectQuery("^SELECT \\* FROM `order_book_dtos` WHERE exchange = \\? AND pair = \\?$").
		WithArgs("invalid_exchange", "invalid_pair").
		WillReturnError(sql.ErrNoRows)
	orderBooks, err = repo.GetOrderBook("invalid_exchange", "invalid_pair")
	assert.NotNil(t, err)
	assert.Nil(t, orderBooks)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSaveOrderBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT version()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("mock_version"))

	gormDB, err := gorm.Open(clickhouse.New(clickhouse.Config{DriverName: "clickhouse", Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("error creating gorm DB: %v", err)
	}

	repo := NewOrderRepository(gormDB)

	orderBook := []*entity.OrderBook{
		{
			ID:       1,
			Exchange: "exchange1",
			Pair:     "pair1",
			Asks:     nil,
			Bids:     nil,
		},
	}

	mock.ExpectExec("^INSERT INTO `order_book_dtos` (.+) VALUES (.+)$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SaveOrderBook(orderBook)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.Error(t, err)
}

func TestGetOrderHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT version()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("mock_version"))

	gormDB, err := gorm.Open(clickhouse.New(clickhouse.Config{DriverName: "clickhouse", Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("error creating gorm DB: %v", err)
	}

	repo := NewOrderRepository(gormDB)

	client := &entity.Client{
		ClientName:   "client1",
		ExchangeName: "exchange1",
		Label:        "label1",
		Pair:         "pair1",
	}

	mock.ExpectQuery("^SELECT \\* FROM `history_orders` WHERE client_name = \\? AND exchange_name = \\? AND label = \\? AND pair = \\?$").
		WithArgs("client1", "exchange1", "label1", "pair1").
		WillReturnRows(sqlmock.NewRows([]string{"client_name", "exchange_name", "label", "pair", "side", "type", "base_qty", "price", "algorithm_name_placed", "lowest_sell_prc", "highest_buy_prc", "commission_quote_qty", "time_placed"}).
			AddRow("client1", "exchange1", "label1", "pair1", "buy", "market", 1.0, 100.0, "algo1", 105.0, 95.0, 0.1, time.Now()))

	// Test case: valid data
	historyOrders, err := repo.GetOrderHistory(client)
	assert.NoError(t, err)
	assert.NotNil(t, historyOrders)
	assert.Len(t, historyOrders, 1)

	//Test case: invalid data
	client = &entity.Client{
		ClientName:   "nonexistent_client",
		ExchangeName: "nonexistent_exchange",
		Label:        "nonexistent_label",
		Pair:         "nonexistent_pair",
	}
	historyOrders, err = repo.GetOrderHistory(client)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Nil(t, historyOrders)

	// Verify all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSaveOrderHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT version()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("mock_version"))

	gormDB, err := gorm.Open(clickhouse.New(clickhouse.Config{DriverName: "clickhouse", Conn: db}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		t.Fatalf("error creating gorm DB: %v", err)
	}
	repo := NewOrderRepository(gormDB)

	orderHistory := entity.HistoryOrder{
		ClientName:   "test_client",
		ExchangeName: "test_exchange",
		Label:        "test_label",
		Pair:         "test_pair",
	}

	err = repo.SaveOrderHistory(orderHistory)
	assert.Error(t, err)
}
