package exchanges

// NovaExchangeMarkets is the response of the market depth query on Nova Exchange
type NovaExchangeMarkets struct {
	Items []struct {
		Amount         string `json:"amount"`
		Baseamount     string `json:"baseamount"`
		Basecurrency   string `json:"basecurrency"`
		Currency       string `json:"currency"`
		Datestamp      string `json:"datestamp"`
		Orderid        int    `json:"orderid"`
		Price          string `json:"price"`
		Tradetype      string `json:"tradetype"`
		UnixTDatestamp int    `json:"unix_t_datestamp"`
	} `json:"items"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
