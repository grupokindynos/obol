package exchanges

// LukkiMarkets is the response of the market depth query on Binance Exchange
type LukkiMarkets struct {
	Total string `json:"total"`
	Data  []struct {
		Direction int    `json:"direction"`
		Price     string `json:"price"`
		Amount    string `json:"amount"`
		Ticker    string `json:"ticker"`
	} `json:"data"`
	Page int `json:"page"`
}
