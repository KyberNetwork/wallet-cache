package fetcher

import (
	"github.com/KyberNetwork/cache/ethereum"
	bFetcher "github.com/KyberNetwork/cache/fetcher/blockchain-fetcher"
)

type RateUSD struct {
	Symbol   string `json:"symbol"`
	PriceUsd string `json:"price_usd"`
}

type FetcherInterface interface {
	EthCall(string, string) (string, error)
	GetLatestBlock() (string, error)
	// GetEvents(string, string, string, string) (*[]ethereum.EventRaw, error)

	// GetRateUsd([]string) ([]io.ReadCloser, error)
	// GetGasPrice() (*ethereum.GasPrice, error)

	GetTypeName() string

	GetRate(string, string) (string, error)
	GetGasPrice() (*ethereum.GasPrice, error)
	// GetRateUsdEther() (string, error)

	// GetGeneralInfo(string) (*ethereum.TokenGeneralInfo, error)

	// get data from tracker
	// GetTrackerData(trackerEndpoint string) (map[string]*ethereum.Rates, error)
}

//var transactionPersistent = models.NewTransactionPersister()

func NewFetcherIns(typeName string, endpoint string, apiKey string) (FetcherInterface, error) {
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
