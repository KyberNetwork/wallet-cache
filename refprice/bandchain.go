package refprice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	ENDPOINT = "https://poa-api.bandchain.org/oracle/request_prices"
)

type BandchainFetcher struct {

}

func NewBandchainFetcher() *BandchainFetcher {
	return &BandchainFetcher{}
}

type bandchainRequestPrices struct {
	Symbols  []string `json:"symbols"`
	MinCount uint64   `json:"min_count"`
	AskCount uint64   `json:"ask_count"`
}

func (f *BandchainFetcher) GetRefPrice(base, quote string) (*big.Float, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	var errMsg = fmt.Sprintf("cannot get price of %v_%v from bandchain", base, quote)

	reqBodyBytes := new(bytes.Buffer)
	symbols := []string{base}

	if strings.ToUpper(quote) != "USD" {
		symbols = append(symbols, quote)
	}

	if err := json.NewEncoder(reqBodyBytes).Encode(bandchainRequestPrices{
		Symbols: symbols,
		MinCount: 3,
		AskCount: 4},
	); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", ENDPOINT, reqBodyBytes)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var (
		prices = make(map[string]*big.Float)
	)
	if res.StatusCode == 200 {
		rawPrices, err := decodeSuccess(res)
		if err != nil {
			return nil, err
		}
		for idx, rawPrice := range rawPrices {
			price, err := decodePrice(rawPrice)
			if err != nil {
				return nil, err
			}
			prices[symbols[idx]] = GetTokenPrice(new(big.Int).SetUint64(price.Px), new(big.Int).SetUint64(price.Multiplier))
		}
	} else {
		return nil, errors.New(errMsg)
	}
	if strings.ToUpper(quote) != "USD" && len(prices) > 1 {
		return new(big.Float).Quo(prices[base], prices[quote]), nil
	} else {
		return prices[base], nil
	}
	return nil, errors.New(errMsg)
}

type bandchainPrice struct {
	Multiplier  uint64 `json:"multiplier"`
	Px          uint64 `json:"px"`
	ResolveTime int64  `json:"resolve_time"`
}

type bandchainRawPrice struct {
	Multiplier  string `json:"multiplier"`
	Px          string `json:"px"`
	ResolveTime string `json:"resolve_time"`
}

func decodeErr(res *http.Response) (map[string]float64, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return nil, errors.New(fmt.Sprintf("%v", result["error"]))
}

func decodeSuccess(res *http.Response) ([]bandchainRawPrice, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result struct {
		Height string     `json:"height"`
		Result []bandchainRawPrice `json:"result"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Result, nil
}

func decodePrice(rawPrice bandchainRawPrice) (bandchainPrice, error) {
	multiplier, err := strconv.ParseUint(rawPrice.Multiplier, 10, 64)
	if err != nil {
		return bandchainPrice{}, err
	}
	px, err := strconv.ParseUint(rawPrice.Px, 10, 64)
	if err != nil {
		return bandchainPrice{}, err
	}
	resolveTime, err := strconv.ParseInt(rawPrice.ResolveTime, 10, 64)
	if err != nil {
		return bandchainPrice{}, err
	}
	return bandchainPrice{
		Multiplier:  multiplier,
		Px:          px,
		ResolveTime: resolveTime,
	}, nil
}
