package coinfactory

import (
	"github.com/grupokindynos/obol/config"
	"strings"
)

// Coin is the basic coin structure to get the correct properties for each coin.
type Coin struct {
	Tag      string
	Name     string
	Exchange string
}

var (
	bitcoin = Coin{
		Tag:      "BTC",
		Name:     "bitcoin",
		Exchange: "bitso",
	}
	colossus = Coin{
		Tag:      "COLX",
		Name:     "colossus",
		Exchange: "cryptobridge",
	}
	dash = Coin{
		Tag:      "DASH",
		Name:     "dash",
		Exchange: "binance",
	}
	digibyte = Coin{
		Tag:      "DGB",
		Name:     "digibyte",
		Exchange: "bittrex",
	}
	groestlcoin = Coin{
		Tag:      "GRS",
		Name:     "groestlcoin",
		Exchange: "binance",
	}
	litecoin = Coin{
		Tag:      "LTC",
		Name:     "litecoin",
		Exchange: "binance",
	}
	mnpcoin = Coin{
		Tag:      "MNP",
		Name:     "mnpcoin",
		Exchange: "crex24",
	}
	polis = Coin{
		Tag:      "POLIS",
		Name:     "polis",
		Exchange: "cryptobridge",
	}
	snowgem = Coin{
		Tag:      "XSG",
		Name:     "snowgem",
		Exchange: "stex",
	}
	zcoin = Coin{
		Tag:      "XZC",
		Name:     "zcoin",
		Exchange: "binance",
	}
)

//Coins is the main array where used coins are stored
type Coins []Coin

//CoinFactory refers to the coins that are being used on the API instance
var CoinFactory = Coins{
	bitcoin,
	polis,
	dash,
	litecoin,
	digibyte,
	groestlcoin,
	zcoin,
	colossus,
	mnpcoin,
	snowgem,
}

//GetCoin is the safe way to check if a coin exists and retrieve the coin data
func GetCoin(tag string) (*Coin, error) {
	for _, v := range CoinFactory {
		if v.Tag == strings.ToUpper(tag) {
			return &v, nil
		}
	}
	return nil, config.ErrorCoinNotAvailable
}
