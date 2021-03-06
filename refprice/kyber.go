package refprice

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
)

type KyberFetcher struct {}


func NewKyberFetcher() *KyberFetcher {
	return &KyberFetcher{}
}


type kyberPrices struct {
	PriceETH float64 `json:"price_ETH"`
	PriceUSD float64	`json:"price_USD"`
}

type kyberPricesResponse struct {
	Data map[string]kyberPrices `json:"data"`
	Error bool `json:"error"`
	Timestamp int64	`json:"timestamp"`
}

func (f *KyberFetcher) GetRefPrice(base, quote string) (*big.Float, error) {
	kyberAPIEndpoint := os.Getenv("KYBER_API_ENDPOINT")
	endpoint := fmt.Sprintf("%v/prices", kyberAPIEndpoint)
	errmsg := fmt.Sprintf("cannot get kyber rate for %v_%v", base, quote)
	b, err := HTTPCall(endpoint)
	if err != nil {
		return nil, err
	}

	var response kyberPricesResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	var (
		mapPrices = response.Data
		priceBase, priceQuote *big.Float
	)
	if v, ok := mapPrices[strings.ToUpper(base)]; ok {
		if v.PriceUSD == 0 {
			return big.NewFloat(0), errors.New("rate is zero")
		}
		priceBase = big.NewFloat(v.PriceUSD)
	} else {
		return nil, errors.New(errmsg)
	}

	if strings.ToUpper(quote) == "USD" {
		return priceBase, nil
	}

	if v, ok := mapPrices[strings.ToUpper(quote)]; ok {
		if v.PriceUSD == 0 {
			return big.NewFloat(0), errors.New("rate is zero")
		}
		priceQuote = big.NewFloat(v.PriceUSD)
	} else {
		return nil, errors.New(errmsg)
	}

	result := new(big.Float).Quo(priceBase, priceQuote)
	if result.Cmp(big.NewFloat(0)) == 0 {
		return big.NewFloat(0), errors.New("rate is zero")
	}
	return result, nil
}