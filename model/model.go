package model

type Ticker struct {
	Exchange       string  `json:"exchange"`
	Currency       string  `json:"currency"`
	Price          float64 `json:"price"`
	YesterdayPrice float64 `json:"yesterday_price"`
	Change         float64 `json:"change"`
	ChangeRate     float64 `json:"change_rate"`
	Volume         uint    `json:"volume"`
}
