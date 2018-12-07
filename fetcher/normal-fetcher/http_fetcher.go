package nFetcher

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/KyberNetwork/server-go/ethereum"
	fCommon "github.com/KyberNetwork/server-go/fetcher/fetcher-common"
)

const (
	TIME_TO_DELETE = 18000
)

type HTTPFetcher struct {
	tradingAPIEndpoint string
}

func NewHTTPFetcher(tradingAPIEndpoint string) *HTTPFetcher {
	return &HTTPFetcher{
		tradingAPIEndpoint: tradingAPIEndpoint,
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
