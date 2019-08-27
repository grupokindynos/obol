package exchanges

// StexRate is the response of the rates query on Stex Exchange
type StexRate struct {
	Success bool `json:"success"`
	Data    struct {
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
		FiatsRate        struct {
			USD float64 `json:"USD"`
			EUR float64 `json:"EUR"`
			UAH int     `json:"UAH"`
			AUD float64 `json:"AUD"`
			IDR int     `json:"IDR"`
			CNY int     `json:"CNY"`
			KRW int     `json:"KRW"`
			JPY int     `json:"JPY"`
			VND int     `json:"VND"`
			INR int     `json:"INR"`
			GBP float64 `json:"GBP"`
			CAD float64 `json:"CAD"`
			BRL int     `json:"BRL"`
			RUB int     `json:"RUB"`
		} `json:"fiatsRate"`
		Timestamp int64 `json:"timestamp"`
	} `json:"data"`
}

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
