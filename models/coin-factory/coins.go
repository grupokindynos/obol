package coin_factory

import (
	"github.com/grupokindynos/obol/config"
	"strings"
)

type Coin struct {
	Tag      string
	Name     string
	Exchange string
}

var (
	Bitcoin = Coin{
		Tag:      "BTC",
		Name:     "bitcoin",
		Exchange: "bitso",
	}
	ColossusXT = Coin{
		Tag:      "COLX",
		Name:     "colossus",
		Exchange: "cryptobridge",
	}
	Dash = Coin{
		Tag:      "DASH",
		Name:     "dash",
		Exchange: "binance",
	}
	DigiByte = Coin{
		Tag:      "DGB",
		Name:     "digibyte",
		Exchange: "binance",
	}
	GroestlCoin = Coin{
		Tag:      "GRS",
		Name:     "groestlcoin",
		Exchange: "binance",
	}
	Litecoin = Coin{
		Tag:      "LTC",
		Name:     "litecoin",
		Exchange: "binance",
	}
	MNPCoin = Coin{
		Tag:      "MNP",
		Name:     "mnpcoin",
		Exchange: "crex24",
	}
	Polis = Coin{
		Tag:      "POLIS",
		Name:     "polis",
		Exchange: "cryptobridge",
	}
	SnowGem = Coin{
		Tag:      "XSG",
		Name:     "snowgem",
		Exchange: "stex",
	}
	ZCoin = Coin{
		Tag:      "XZC",
		Name:     "zcoin",
		Exchange: "binance",
	}
)

type Coins []Coin

var CoinFactory = Coins{
	Bitcoin,
	Polis,
	Dash,
	Litecoin,
	DigiByte,
	GroestlCoin,
	ZCoin,
	ColossusXT,
	MNPCoin,
	SnowGem,
}

func GetCoin(tag string) (*Coin, error) {
	for _, v := range CoinFactory {
		if v.Tag == strings.ToUpper(tag) {
			return &v, nil
		}
	}
	return nil, config.ErrorCoinNotAvailable
}
