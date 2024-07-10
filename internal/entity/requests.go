package entity

type OrderBookRequest struct {
	Exchange_name string `json:"exchange"`
	Pair          string `json:"pair"`
}
