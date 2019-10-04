package exchanges

import "time"

// GraviexMarkets is the response of the market depth query on Graviex Exchange
type GraviexMarkets struct {
	Asks []struct {
		ID              int         `json:"id"`
		At              int         `json:"at"`
		Side            string      `json:"side"`
		OrdType         string      `json:"ord_type"`
		Price           string      `json:"price"`
		AvgPrice        string      `json:"avg_price"`
		State           string      `json:"state"`
		Market          string      `json:"market"`
		CreatedAt       time.Time   `json:"created_at"`
		Volume          string      `json:"volume"`
		RemainingVolume string      `json:"remaining_volume"`
		ExecutedVolume  string      `json:"executed_volume"`
		TradesCount     int         `json:"trades_count"`
		Strategy        interface{} `json:"strategy"`
	} `json:"asks"`
	Bids []struct {
		ID              int         `json:"id"`
		At              int         `json:"at"`
		Side            string      `json:"side"`
		OrdType         string      `json:"ord_type"`
		Price           string      `json:"price"`
		AvgPrice        string      `json:"avg_price"`
		State           string      `json:"state"`
		Market          string      `json:"market"`
		CreatedAt       time.Time   `json:"created_at"`
		Volume          string      `json:"volume"`
		RemainingVolume string      `json:"remaining_volume"`
		ExecutedVolume  string      `json:"executed_volume"`
		TradesCount     int         `json:"trades_count"`
		Strategy        interface{} `json:"strategy"`
	} `json:"bids"`
}
