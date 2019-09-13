package models

import "time"

// MarketOrder is just a simple struct to homologue the exchange responses
type MarketOrder struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}

// FiatRates is the struct to handle internally the OpenRate response
type FiatRates struct {
	Rates       map[string]float64
	LastUpdated time.Time
}

// Rate is the struct to homologue exchange rate responses
type Rate struct {
	Code string  `json:"code"`
	Name string  `json:"name"`
	Rate float64 `json:"rate"`
}

// OpenRates is the structure of the OpenRate response
type OpenRates struct {
	Rates struct {
		CAD float64 `json:"CAD"`
		HKD float64 `json:"HKD"`
		ISK float64 `json:"ISK"`
		PHP float64 `json:"PHP"`
		DKK float64 `json:"DKK"`
		HUF float64 `json:"HUF"`
		CZK float64 `json:"CZK"`
		AUD float64 `json:"AUD"`
		RON float64 `json:"RON"`
		SEK float64 `json:"SEK"`
		IDR float64 `json:"IDR"`
		INR float64 `json:"INR"`
		BRL float64 `json:"BRL"`
		RUB float64 `json:"RUB"`
		HRK float64 `json:"HRK"`
		JPY float64 `json:"JPY"`
		THB float64 `json:"THB"`
		CHF float64 `json:"CHF"`
		SGD float64 `json:"SGD"`
		PLN float64 `json:"PLN"`
		BGN float64 `json:"BGN"`
		TRY float64 `json:"TRY"`
		CNY float64 `json:"CNY"`
		NOK float64 `json:"NOK"`
		NZD float64 `json:"NZD"`
		ZAR float64 `json:"ZAR"`
		USD float64 `json:"USD"`
		MXN float64 `json:"MXN"`
		EUR float64 `json:"EUR"`
		ILS float64 `json:"ILS"`
		GBP float64 `json:"GBP"`
		KRW float64 `json:"KRW"`
		MYR float64 `json:"MYR"`
	} `json:"rates"`
	Base string `json:"base"`
	Date string `json:"date"`
}

var OpenRateNames = map[string]string{
	"CAD": "Canadian Dollar",
	"HKD": "Hong Kong Dollar",
	"ISK": "Iceland Krona",
	"PHP": "Philippine Peso",
	"DKK": "Danish Krone",
	"HUF": "Forint",
	"CZK": "Czech Koruna",
	"AUD": "Australian Dollar",
	"RON": "Leu",
	"SEK": "Swedish Krona",
	"IDR": "Rupiah",
	"INR": "Indian Rupee",
	"BRL": "Brazilian Real",
	"RUB": "Russian Ruble",
	"HRK": "Croatian Kuna",
	"JPY": "Yen",
	"THB": "Baht",
	"CHF": "Swiss Franc",
	"SGD": "Singapore Dollar",
	"PLN": "PZloty",
	"BGN": "Bulgarian Lev",
	"TRY": "Turkish Lira",
	"CNY": "Yuan",
	"NOK": "Norwegian Krone",
	"NZD": "New Zealand Dollar",
	"ZAR": "Rand",
	"USD": "US Dollar",
	"MXN": "Mexican Peso",
	"EUR": "Euro",
	"ILS": "New Israeli Shekel",
	"GBP": "Pound Sterling",
	"KRW": "South Korean Won",
	"MYR": "Malaysian Ringgit",
	"BTC": "Bitcoin",
}
