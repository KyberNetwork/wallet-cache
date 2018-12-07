package nfetcher

import (
	"errors"

	"github.com/KyberNetwork/server-go/ethereum"
)

type HTTPFetcher struct {
}

func NewCGFetcher() *CGFetcher {
	return &CGFetcher{
		API:        "https://api.coingecko.com/api/v3",
		typeMarket: "coingecko",
	}
}

func (self *CGFetcher) GetRateUsdEther() (string, error) {
	return nil, errors.New("not support this func")
}

func (self *CGFetcher) GetGeneralInfo(coinID string) (*ethereum.TokenGeneralInfo, error) {
	return nil, errors.New("not support this func")
}

// func (self *CGFetcher) GetTypeMarket() string {
// 	return self.typeMarket
// }

func (self *) GetListToken() {
	
}