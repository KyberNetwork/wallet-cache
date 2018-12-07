package fetcher

import (
	"github.com/KyberNetwork/server-go/ethereum"
	nFetcher "github.com/KyberNetwork/server-go/fetcher/normal-fetcher"
)

type FetcherNormalInterface interface {
	GetRateUsdEther() (string, error)
	GetGeneralInfo(string) (*ethereum.TokenGeneralInfo, error)
	GetTypeMarket() string
}

//var transactionPersistent = models.NewTransactionPersister()

func NewFetcherNormalIns(typeName string) FetcherNormalInterface {
	var fetcher FetcherNormalInterface
	switch typeName {
	case "cmc":
		fetcher = nFetcher.NewCMCFetcher()
		break
	case "coingecko":
		fetcher = nFetcher.NewCGFetcher()
		break
	}
	return fetcher
}
