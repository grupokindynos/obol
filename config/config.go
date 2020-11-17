package config

import (
	"errors"
	"net/http"
	"time"
)

var (
	// FixerRatesURL is the base URL for fiat rates based on FixerRates
	FixerRatesURL = "http://data.fixer.io/api/latest"
	// ErrorNoServiceForCoin returns when is not possible to load the service for an exchange
	ErrorNoServiceForCoin = errors.New("unable to load exchange for this coin")
	// ErrorNoFallBackServiceForCoin returns when is not possible to load the fallback service for an exchange
	ErrorNoFallBackServiceForCoin = errors.New("unable to load fallback exchange for this coin")
	// ErrorNoC2CWithSameCoin returns when a C2C rate is called using the same coins
	ErrorNoC2CWithSameCoin = errors.New("cannot use the same coin on both parameters")
	// ErrorInvalidAmountOnC2C returns when try to call a C2C rate with amount, and the amount is not properly formated.
	ErrorInvalidAmountOnC2C = errors.New("invalid amount to convert from coin to coin")
	// ErrorUnknownIdForCoin returns when an id for a coin is not defined
	ErrorUnknownIdForCoin = errors.New("unknown id for coin")
	// ErrorRequestTimeout returns when a request timeout threshold is reached
	ErrorRequestTimeout = errors.New("request timeout")
	// HttpClient is a wrapper to properly timeout http.get calls
	HttpClient = &http.Client{
		Timeout: time.Second * 45,
	}
)
