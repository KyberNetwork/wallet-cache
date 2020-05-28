package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/KyberNetwork/cache/common"
	"github.com/KyberNetwork/cache/ethereum"
	fCommon "github.com/KyberNetwork/cache/fetcher/fetcher-common"
)

const (
	TIME_TO_DELETE  = 18000
	API_KEY_TRACKER = "jHGlaMKcGn5cCBxQCGwusS4VcnH0C6tN"
)

type HTTPFetcher struct {
	tradingAPIEndpoint string
	gasStationEndPoint string
	apiEndpoint        string
}

func NewHTTPFetcher(tradingAPIEndpoint, gasStationEndpoint, apiEndpoint string) *HTTPFetcher {
	return &HTTPFetcher{
		tradingAPIEndpoint: tradingAPIEndpoint,
		gasStationEndPoint: gasStationEndpoint,
		apiEndpoint:        apiEndpoint,
	}
}

func (self *HTTPFetcher) GetListToken() ([]ethereum.Token, error) {
	b, err := fCommon.HTTPCall(self.tradingAPIEndpoint)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	var result ethereum.TokenConfig
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if result.Success == false {
		err = errors.New("Cannot get list token")
		return nil, err
	}
	data := result.Data
	if len(data) == 0 {
		err = errors.New("list token from api is empty")
		return nil, err
	}
	return data, nil
}

type GasStation struct {
	Fast     float64 `json:"fast"`
	Standard float64 `json:"average"`
	Low      float64 `json:"safeLow"`
}

func (self *HTTPFetcher) GetGasPrice() (*ethereum.GasPrice, error) {
	b, err := fCommon.HTTPCall(self.gasStationEndPoint)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	var gasPrice GasStation
	err = json.Unmarshal(b, &gasPrice)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	fast := big.NewFloat(gasPrice.Fast / 10)
	standard := big.NewFloat((gasPrice.Fast + gasPrice.Standard) / 20)
	low := big.NewFloat(gasPrice.Low / 10)
	defaultGas := standard

	return &ethereum.GasPrice{
		fast.String(), standard.String(), low.String(), defaultGas.String(),
	}, nil
}

// get data from tracker.kyber

func (self *HTTPFetcher) GetRate7dData() (map[string]*ethereum.Rates, error) {
	trackerAPI := self.apiEndpoint + "/rates7d"
	b, err := fCommon.HTTPCall(trackerAPI)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	trackerData := map[string]*ethereum.Rates{}
	err = json.Unmarshal(b, &trackerData)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return trackerData, nil
}

func (self *HTTPFetcher) GetUserInfo(url string) (*common.UserInfo, error) {
	userInfo := &common.UserInfo{}
	b, err := fCommon.HTTPCall(url)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	err = json.Unmarshal(b, userInfo)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return userInfo, nil
}

type TokenPrice struct {
	Data []struct {
		Symbol string  `json:"symbol"`
		Price  float64 `json:"price"`
	} `json:"data"`
	Error      bool   `json:"error"`
	TimeUpdate uint64 `json:"timeUpdated"`
}

// GetRateUsdEther get usd from api
func (self *HTTPFetcher) GetRateUsdEther() (string, error) {
	var ethPrice string
	url := fmt.Sprintf("%s/token_price?currency=USD", self.apiEndpoint)
	b, err := fCommon.HTTPCall(url)
	if err != nil {
		log.Print(err)
		return ethPrice, err
	}
	var tokenPrice TokenPrice
	err = json.Unmarshal(b, &tokenPrice)
	if err != nil {
		log.Println(err)
		return ethPrice, err
	}
	if tokenPrice.Error {
		return ethPrice, errors.New("cannot get token price from api")
	}
	for _, v := range tokenPrice.Data {
		if v.Symbol == common.ETHSymbol {
			ethPrice = fmt.Sprintf("%.6f", v.Price)
			break
		}
	}
	return ethPrice, nil
}

type TokenRate struct {
	BaseSymbol  string  `json:"base_symbol"`
	QuoteSymbol string  `json:"quote_symbol"`
	RateSell    float64 `json:"current_bid"`
	RateBuy     float64 `json:"current_ask"`
}

type MarketData struct {
	Data []TokenRate `json:"data"`
	Err  bool        `json:"error"`
}

// GetRateUsdEther get usd from api
func (self *HTTPFetcher) GetRate() ([]ethereum.Rate, error) {

	url := fmt.Sprintf("%s/market", self.apiEndpoint)
	b, err := fCommon.HTTPCall(url)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	tokenRate := MarketData{}
	err = json.Unmarshal(b, &tokenRate)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if !tokenRate.Err {
		rates := make([]ethereum.Rate, 0)

		for _, rate := range tokenRate.Data {
			rates = append(rates, getRateBuy(rate))
			rates = append(rates, getRateSell(rate))
		}

		return rates, nil
	}

	return nil, errors.New("Cannot get rate")
}

func getRateBuy(rate TokenRate) ethereum.Rate {
	if rate.RateBuy == 0 {
		return ethereum.Rate{
			Source:  rate.QuoteSymbol,
			Dest:    rate.BaseSymbol,
			Rate:    "0",
			Minrate: "0",
		}
	}
	rateBuy := 1 / rate.RateBuy
	minRate := rateBuy * 0.97

	rateBig := common.ToWei(rateBuy, 18)
	minRateBig := common.ToWei(minRate, 18)

	return ethereum.Rate{
		Source:  rate.QuoteSymbol,
		Dest:    rate.BaseSymbol,
		Rate:    rateBig.String(),
		Minrate: minRateBig.String(),
	}
}

func getRateSell(rate TokenRate) ethereum.Rate {
	rateBig := common.ToWei(rate.RateSell, 18)
	minRateBig := common.ToWei(rate.RateSell*0.97, 18)

	return ethereum.Rate{
		Source:  rate.BaseSymbol,
		Dest:    rate.QuoteSymbol,
		Rate:    rateBig.String(),
		Minrate: minRateBig.String(),
	}
}

type QuoteData struct {
	Data string `json:"data"`
	Err  bool   `json:"error"`
}

func (self *HTTPFetcher) GetQuoteAmount(quote string, base string, baseAmount string, typeQ string) (string, error) {
	url := fmt.Sprintf("%s/quote_amount?base=%s&quote=%s&base_amount=%s&type=%s", self.apiEndpoint, base, quote, baseAmount, typeQ)
	b, err := fCommon.HTTPCall(url)
	if err != nil {
		log.Print(err)
		return "", err
	}

	quoteData := QuoteData{}
	err = json.Unmarshal(b, &quoteData)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if !quoteData.Err {
		return quoteData.Data, nil
	}

	return "", errors.New("Cannot get quote data")
}
