package models

type MarketOrder struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}

type Rate struct {
	Code string  `json:"code"`
	Name string  `json:"name"`
	Rate float64 `json:"rate"`
}

type BitpayRates struct {
	Data []Rate `json:"data"`
}
