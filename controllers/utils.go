package controllers

import "errors"

func getExchangeRateMargin(exchange string) (float64, error){
	mp := getExchangesMapRateMargin()
	if val, ok := mp[exchange]; ok {
		return val, nil
	}
	return 0, errors.New("exchange not found on MapRateMargin")
}

func getExchangesMapRateMargin() map[string]float64 {
	return map[string]float64 {
		"binance" : 1.05,
		"stex" : 1.05,
		"bittrex" : 1.05,
		"crex24" : 2.0,
		"southxchange" : 2.0,
	}
}