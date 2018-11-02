package fetcher

import (
	"github.com/KyberNetwork/server-go/ethereum"
	bFetcher "github.com/KyberNetwork/server-go/fetcher/blockchain-fetcher"
)

type FetcherNormalInterface interface {
	GetRateUsdEther() (string, error)
	GetGeneralInfo(string) (*ethereum.TokenGeneralInfo, error)
}

//var transactionPersistent = models.NewTransactionPersister()

func NewFetcherNormalIns(typeName string, endpoint string, apiKey string) (FetcherNormalInterface, error) {
	var fetcher FetcherInterface
	var err error
	switch typeName {
	case "etherscan":
		fetcher, err = bFetcher.NewEtherScan(typeName, endpoint, apiKey)
		break
	case "node":
		fetcher, err = bFetcher.NewBlockchainFetcher(typeName, endpoint, apiKey)
		break
	}
	return fetcher, err
}
