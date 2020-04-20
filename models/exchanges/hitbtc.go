package exchanges

import "time"

// HitBTCRate is the response of the rates query on HitBTC Exchange
type HitBTCRate map[string]data

type data struct {
	Symbol string `json:"symbol"`
	Ask    []struct {
		Price string `json:"price"`
		Size  string `json:"size"`
	} `json:"ask"`
	Bid []struct {
		Price string `json:"price"`
		Size  string `json:"size"`
	} `json:"bid"`
	Timestamp time.Time `json:"timestamp"`
}
