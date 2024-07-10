package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/egorque1/vortex-test/docs"
	"github.com/egorque1/vortex-test/internal/db"
	"github.com/egorque1/vortex-test/internal/modules/controller"
	"github.com/egorque1/vortex-test/internal/modules/repository"
	"github.com/egorque1/vortex-test/internal/modules/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/httprate"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title swagger Order Management API
// @version 1.0
// @description This is a sample server for managing orders.
// @BasePath /

func main() {
	database, err := db.Connect(".env")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to ClickHouse!")

	if err := db.Migrate(database); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	orderBookRepo := repository.NewOrderRepository(database)
	orderBookService := service.NewOrderService(orderBookRepo)
	orderBookController := controller.NewController(orderBookRepo, orderBookService)

	r := chi.NewMux()

	r.Mount("/swagger", httpSwagger.WrapHandler)

	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(100, 1*time.Second))

		r.Get("/orderbook", orderBookController.GetOrderBookHandler)
		r.Get("/history", orderBookController.GetOrderHistoryHandler)
	})
	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(200, 1*time.Second))

		r.Post("/orderbook", orderBookController.SaveOrderBookHandler)
		r.Post("/history", orderBookController.SaveOrderHistoryHandler)
	})

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}
