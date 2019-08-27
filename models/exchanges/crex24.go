package exchanges

// Crex24Rates is the response of the rates query on Crex24 Exchange
type Crex24Rates struct {
	Error   interface{} `json:"Error"`
	Tickers []struct {
		PairID        int     `json:"PairId"`
		PairName      string  `json:"PairName"`
		Last          float64 `json:"Last"`
		LowPrice      float64 `json:"LowPrice"`
		HighPrice     float64 `json:"HighPrice"`
		PercentChange float64 `json:"PercentChange"`
		BaseVolume    float64 `json:"BaseVolume"`
		QuoteVolume   float64 `json:"QuoteVolume"`
		VolumeInBtc   float64 `json:"VolumeInBtc"`
		VolumeInUsd   float64 `json:"VolumeInUsd"`
	} `json:"Tickers"`
}

// Crex24Markets is the response of the market depth query on Crex24 Exchange
type Crex24Markets struct {
	BuyLevels []struct {
		Price  float64 `json:"price"`
		Volume float64 `json:"volume"`
	} `json:"buyLevels"`
	SellLevels []struct {
		Price  float64 `json:"price"`
		Volume float64 `json:"volume"`
	} `json:"sellLevels"`
}
