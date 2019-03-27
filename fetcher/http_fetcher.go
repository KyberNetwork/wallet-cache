package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/KyberNetwork/server-go/common"
	"github.com/KyberNetwork/server-go/ethereum"
	fCommon "github.com/KyberNetwork/server-go/fetcher/fetcher-common"
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
