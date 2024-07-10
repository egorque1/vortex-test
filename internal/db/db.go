package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

func Connect(env string) (*gorm.DB, error) {
	err := godotenv.Load(env)
	if err != nil {
		return nil, err
	}

	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")

	dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s?dial_timeout=10s&read_timeout=20s", user, pass, host, port, name)

	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	orderBooksTable := `
		CREATE TABLE IF NOT EXISTS order_books (
			id Int64,
			exchange String,
			pair String,
			asks String,
			bids String
		) ENGINE = MergeTree()
		PRIMARY KEY (exchange, pair)
		ORDER BY (exchange, pair);
	`
	if err := db.Exec(orderBooksTable).Error; err != nil {
		return fmt.Errorf("error creating order_books table: %w", err)
	}

	historyOrdersTable := `
		CREATE TABLE IF NOT EXISTS history_orders (
			client_name String,
			exchange_name String,
			label String,
			pair String,
			side String,
			type String,
			base_qty Float64,
			price Float64,
			algorithm_name_placed String,
			lowest_sell_prc Float64,
			highest_buy_prc Float64,
			commission_quote_qty Float64,
			time_placed DateTime
		) ENGINE = MergeTree()
		PRIMARY KEY (client_name, exchange_name, pair)
		ORDER BY (client_name, exchange_name, pair);
	`
	if err := db.Exec(historyOrdersTable).Error; err != nil {
		return fmt.Errorf("error creating history_orders table: %w", err)
	}

	return nil
}
