package fetcher

import (
	"io"

	"github.com/KyberNetwork/server-go/ethereum"
)

type RateUSD struct {
	Symbol   string `json:"symbol"`
	PriceUsd string `json:"price_usd"`
}

type FetcherInterface interface {
	EthCall(string, string) (string, error)
	GetLatestBlock() (string, error)
	GetEvents(string, string, string, string) (*[]ethereum.EventRaw, error)

	GetRateUsd([]string) ([]io.ReadCloser, error)
	GetGasPrice() (*ethereum.GasPrice, error)

	GetTypeName() string

	GetRateUsdEther() (string, error)

	GetGeneralInfo(string) (*ethereum.TokenGeneralInfo, error)

	// get data from tracker
	GetTrackerData(trackerEndpoint string) (map[string]*ethereum.Rates, error)
}

//var transactionPersistent = models.NewTransactionPersister()

func NewFetcherIns(typeName string, endpoint string, apiKey string) (FetcherInterface, error) {
	var fetcher FetcherInterface
	var err error
	switch typeName {
	case "etherscan":
		fetcher, err = NewEtherScan(typeName, endpoint, apiKey)
		break
	case "node":
		fetcher, err = NewBlockchainFetcher(typeName, endpoint, apiKey)
		break
	}
	return fetcher, err
}
