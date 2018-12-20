package fetcher

import (
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"time"

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
	trackerEndpoint    string
}

func NewHTTPFetcher(tradingAPIEndpoint, gasStationEndpoint, trackerEndpoint string) *HTTPFetcher {
	return &HTTPFetcher{
		tradingAPIEndpoint: tradingAPIEndpoint,
		gasStationEndPoint: gasStationEndpoint,
		trackerEndpoint:    trackerEndpoint,
	}
}

func (self *HTTPFetcher) GetListToken() (map[string]ethereum.Token, error) {
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
	listToken := make(map[string]ethereum.Token)
	for _, token := range result.Data {
		if token.DelistTime == 0 || uint64(time.Now().UTC().Unix()) <= TIME_TO_DELETE+token.DelistTime {
			listToken[token.Symbol] = token
		}
	}
	return listToken, nil
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

func (self *HTTPFetcher) GetTrackerData() (map[string]*ethereum.Rates, error) {
	trackerAPI := self.trackerEndpoint + "/api/tokens/rates?api_key=" + API_KEY_TRACKER
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
