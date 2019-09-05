package exchanges

// StexMarkets is the response of the market depth query on Stex Exchange
type StexMarkets struct {
	Success bool `json:"success"`
	Data    struct {
		Ask []struct {
			CurrencyPairID   int     `json:"currency_pair_id"`
			Amount           string  `json:"amount"`
			Price            string  `json:"price"`
			Amount2          string  `json:"amount2"`
			Count            int     `json:"count"`
			CumulativeAmount float64 `json:"cumulative_amount"`
		} `json:"ask"`
		Bid []struct {
			CurrencyPairID   int     `json:"currency_pair_id"`
			Amount           string  `json:"amount"`
			Price            string  `json:"price"`
			Amount2          string  `json:"amount2"`
			Count            int     `json:"count"`
			CumulativeAmount float64 `json:"cumulative_amount"`
		} `json:"bid"`
		AskTotalAmount float64 `json:"ask_total_amount"`
		BidTotalAmount float64 `json:"bid_total_amount"`
	} `json:"data"`
}
