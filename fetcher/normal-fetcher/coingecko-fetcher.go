package nFetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/KyberNetwork/server-go/ethereum"
	fCommon "github.com/KyberNetwork/server-go/fetcher/fetcher-common"
)

type CGFetcher struct {
	API        string
	typeMarket string
}

func NewCGFetcher() *CGFetcher {
	return &CGFetcher{
		API:        "https://api.coingecko.com/api/v3",
		typeMarket: "coingecko",
	}
}

func (self *CGFetcher) GetRateUsdEther() (string, error) {
	// typeMarket := self.typeMarket
	url := self.API + "/coins/ethereum"
	b, err := fCommon.HTTPCall(url)
	if err != nil {
		log.Print(err)
		return "", err
	}
	rateItem := ethereum.RateUSDCG{}
	err = json.Unmarshal(b, &rateItem)
	if err != nil {
		log.Print(err)
		return "", err
	}
	rateString := fmt.Sprintf("%.6f", rateItem.MarketData.CurrentPrice.USD)
	return rateString, nil
}

func (self *CGFetcher) GetGeneralInfo(coinID string) (*ethereum.TokenGeneralInfo, error) {
	url := self.API + "/coins/" + coinID
	b, err := fCommon.HTTPCall(url)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	tokenItem := ethereum.TokenInfoCoinGecko{}
	err = json.Unmarshal(b, &tokenItem)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	tokenGenalInfo := tokenItem.ToTokenInfoCMC()
	return &tokenGenalInfo, nil
	err = errors.New("Cannot find data key in return quotes of ticker")
	log.Print(err)
	return nil, err
}

// func (self *CGFetcher) GetTypeMarket() string {
// 	return self.typeMarket
// }
