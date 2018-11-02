package nFetcher

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/KyberNetwork/server-go/ethereum"
	fCommon "github.com/KyberNetwork/server-go/fetcher/fetcher-common"
)

type CGFetcher struct {
	API string
}

func NewCGFetcher() *CGFetcher {
	return &CGFetcher{
		API: "https://api.coingecko.com/api/v3",
	}
}

func (self *CMCFetcher) GetRateUsdEther() (string, error) {
	url := self.API + "/coins/ethereum"
	b, err := fCommon.HTTPCall(url)
	if err != nil {
		log.Print(err)
		return "", err
	}
	rateItem := make([]RateUSD, 0)
	err = json.Unmarshal(b, &rateItem)
	if err != nil {
		log.Print(err)
		return "", err
	}
	return rateItem[0].PriceUsd, nil
}

func (self *CMCFetcher) GetGeneralInfo(coinID string) (*ethereum.TokenGeneralInfo, error) {
	url := self.API + "/coins/" + coinID
	b, err := fCommon.HTTPCall(url)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	tokenItem := map[string]ethereum.TokenInfoCoinGecko{}
	err = json.Unmarshal(b, &tokenItem)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if data, ok := tokenItem["data"]; ok {
		data.MarketCap = data.Quotes["ETH"].MarketCap
		return &data, nil
	}
	err = errors.New("Cannot find data key in return quotes of ticker")
	log.Print(err)
	return nil, err
}
