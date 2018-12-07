package nFetcher

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/KyberNetwork/server-go/ethereum"
	fCommon "github.com/KyberNetwork/server-go/fetcher/fetcher-common"
)

type CMCFetcher struct {
	APIV1      string
	APIV2      string
	typeMarket string
}

func NewCMCFetcher() *CMCFetcher {
	return &CMCFetcher{
		APIV1:      "https://api.coinmarketcap.com/v1",
		APIV2:      "https://api.coinmarketcap.com/v2",
		typeMarket: "cmc",
	}
}

func (self *CMCFetcher) GetRateUsdEther() (string, error) {
	// typeMarket := self.typeMarket
	url := self.APIV1 + "/ticker/ethereum"
	b, err := fCommon.HTTPCall(url)
	if err != nil {
		log.Print(err)
		return "", err
	}
	rateItem := make([]ethereum.RateUSD, 0)
	err = json.Unmarshal(b, &rateItem)
	if err != nil {
		log.Print(err)
		return "", err
	}
	return rateItem[0].PriceUsd, nil
}

func (self *CMCFetcher) GetGeneralInfo(usdId string) (*ethereum.TokenGeneralInfo, error) {
	url := self.APIV2 + "/ticker/" + usdId + "/?convert=ETH"
	b, err := fCommon.HTTPCall(url)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	tokenItem := map[string]ethereum.TokenGeneralInfo{}
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

func (self *CMCFetcher) GetTypeMarket() string {
	return self.typeMarket
}
