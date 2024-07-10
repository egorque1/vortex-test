package controller

import (
	"encoding/json"
	"net/http"

	"github.com/egorque1/vortex-test/internal/entity"
	"github.com/egorque1/vortex-test/internal/modules/repository"
	"github.com/egorque1/vortex-test/internal/modules/service"
	"gorm.io/gorm"
)

type OrderController interface {
	GetOrderBookHandler(w http.ResponseWriter, r *http.Request)
	SaveOrderBookHandler(w http.ResponseWriter, r *http.Request)
	GetOrderHistoryHandler(w http.ResponseWriter, r *http.Request)
	SaveOrderHistoryHandler(w http.ResponseWriter, r *http.Request)
}

type orderControllerImpl struct {
	repo repository.OrderRepository
	svc  service.OrderService
}

func NewController(repo repository.OrderRepository, svc service.OrderService) OrderController {
	return &orderControllerImpl{repo: repo, svc: svc}
}

// @Summary Get Order Book
// @Description Retrieve the order book for a specific exchange and trading pair.
// @Tags order
// @Accept json
// @Produce json
// @Param orderBookRequest body entity.OrderBookRequest true "Order Book Request"
// @Success 200 {array} entity.OrderBook
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /orderbook [get]
func (c *orderControllerImpl) GetOrderBookHandler(w http.ResponseWriter, r *http.Request) {
	var req entity.OrderBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ob, err := c.svc.GetOrderBook(req.Exchange_name, req.Pair)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(ob)
	w.Write(bytes)
}

// @Summary Save Order Book
// @Description Save a new order book entry.
// @Tags order
// @Accept json
// @Produce json
// @Param orderBooks body []entity.OrderBook true "Order Books"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /orderbook [post]
func (c *orderControllerImpl) SaveOrderBookHandler(w http.ResponseWriter, r *http.Request) {
	var req []*entity.OrderBook
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := c.svc.SaveOrderBook(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Get Order History
// @Description Retrieve the order history for a specific client.
// @Tags order
// @Accept json
// @Produce json
// @Param client body entity.Client true "Client"
// @Success 200 {array} entity.HistoryOrder
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /orderhistory [get]
func (c *orderControllerImpl) GetOrderHistoryHandler(w http.ResponseWriter, r *http.Request) {
	var req *entity.Client
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ho, err := c.svc.GetOrderHistory(req)
	if err == gorm.ErrRecordNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(ho)
	w.Write(bytes)
}

// @Summary Save Order History
// @Description Save a new history order entry.
// @Tags order
// @Accept json
// @Produce json
// @Param historyOrder body entity.HistoryOrder true "History Order"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /orderhistory [post]
func (c *orderControllerImpl) SaveOrderHistoryHandler(w http.ResponseWriter, r *http.Request) {
	var req entity.HistoryOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := c.svc.SaveOrderHistory(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
