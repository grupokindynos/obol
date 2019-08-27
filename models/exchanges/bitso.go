package exchanges

import "time"

// BitsoRates is the response of the rates query on Bitso Exchange
type BitsoRates struct {
	Success bool `json:"success"`
	Payload struct {
		High      string    `json:"high"`
		Last      string    `json:"last"`
		CreatedAt time.Time `json:"created_at"`
		Book      string    `json:"book"`
		Volume    string    `json:"volume"`
		Vwap      string    `json:"vwap"`
		Low       string    `json:"low"`
		Ask       string    `json:"ask"`
		Bid       string    `json:"bid"`
		Change24  string    `json:"change_24"`
	} `json:"payload"`
}
