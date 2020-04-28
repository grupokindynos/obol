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

// StexTickers is the response of the market depth query on Stex Exchange
type StexTickers struct {
	Success bool `json:"success"`
	Data    []struct {
		ID               int    `json:"id"`
		AmountMultiplier int    `json:"amount_multiplier"`
		CurrencyCode     string `json:"currency_code"`
		MarketCode       string `json:"market_code"`
		CurrencyName     string `json:"currency_name"`
		MarketName       string `json:"market_name"`
		Symbol           string `json:"symbol"`
		GroupName        string `json:"group_name"`
		GroupID          int    `json:"group_id"`
		Ask              string `json:"ask"`
		Bid              string `json:"bid"`
		Last             string `json:"last"`
		Open             string `json:"open"`
		Low              string `json:"low"`
		High             string `json:"high"`
		Volume           string `json:"volume"`
		VolumeQuote      string `json:"volumeQuote"`
		Count            string `json:"count"`
		Timestamp        int64  `json:"timestamp"`
		GroupPosition    int    `json:"group_position"`
	} `json:"data"`
}
