package Persister

import (
	"github.com/KyberNetwork/server-go/ethereum"
)

// type Rate struct {
// 	Source   string `json:"source"`
// 	Dest     string `json:"dest"`
// 	Rate     string `json:"rate"`
// 	Minrate  string `json:"minRate"`
// }

// type EventHistory struct {
// 	ID               int    `json:"id"`
// 	ActualDestAmount string `json:"actualDestAmount"`
// 	ActualSrcAmount  string `json:"actualSrcAmount"`
// 	Dest             string `json:"dest"`
// 	Source           string `json:"source"`
// 	Sender           string `json:"sender"`
// 	Blocknumber      string `json:"blockNumber"`
// 	Txhash           string `json:"txHash"`
// 	Timestamp        string `json:"timestamp"`
// 	Status           string `json:"status"`
// }

type RateUSD struct {
	Symbol   string `json:"symbol"`
	PriceUsd string `json:"price_usd"`
}

type Persister interface {
	GetRate() []ethereum.Rate
	GetIsNewRate() bool
	SetIsNewRate(bool)
	GetTimeUpdateRate() int64

	SaveRate([]ethereum.Rate, int64)

	SaveGeneralInfoTokens(map[string]*ethereum.TokenGeneralInfo)
	GetTokenInfo() map[string]*ethereum.TokenGeneralInfo

	GetLatestBlock() string
	GetIsNewLatestBlock() bool
	SaveLatestBlock(string) error
	SetNewLatestBlock(bool)

	GetRateUSD() []RateUSD
	GetRateETH() string
	GetIsNewRateUSD() bool
	SaveRateUSD(string) error
	SetNewRateUSD(bool)

	// GetRateUSDCG() []RateUSD
	// GetRateETHCG() string
	// SetNewRateUSDCG(bool)
	// GetIsNewRateUSDCG() bool

	SaveKyberEnabled(bool)
	SetNewKyberEnabled(bool)
	GetKyberEnabled() bool
	GetNewKyberEnabled() bool

	SetNewMaxGasPrice(bool)
	SaveMaxGasPrice(string)
	GetMaxGasPrice() string
	GetNewMaxGasPrice() bool

	SaveGasPrice(*ethereum.GasPrice)
	SetNewGasPrice(bool)
	GetGasPrice() *ethereum.GasPrice
	GetNewGasPrice() bool

	SaveMarketData(rates map[string]*ethereum.Rates, tokens map[string]ethereum.Token)
	GetRightMarketData() map[string]*ethereum.RightMarketInfo
	// GetRightMarketDataCG() map[string]*ethereum.RightMarketInfo
	GetLast7D(listTokens string) map[string][]float64
	GetIsNewTrackerData() bool
	SetIsNewTrackerData(isNewTrackerData bool)
	SetIsNewMarketInfo(isNewMarketInfo bool)
	GetIsNewMarketInfo() bool
	// GetIsNewMarketInfoCG() bool
	GetTimeVersion() string

	IsFailedToFetchTracker() bool
}

//var transactionPersistent = models.NewTransactionPersister()

func NewPersister(name string) (Persister, error) {
	Persister, err := NewRamPersister()
	return Persister, err
}
