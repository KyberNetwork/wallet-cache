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
	//	SetRateToken(map[string]ethereum.Token)

	GetRate() *[]ethereum.Rate

	SaveRate(*[]ethereum.Rate)
	//	SaveNewRate(bool)
	//GetIsNewRate() bool

	SaveGeneralInfoTokens(map[string]*ethereum.TokenGeneralInfo)
	GetTokenInfo() map[string]*ethereum.TokenGeneralInfo
	//SetIsNewGeneralInfoTokens(bool)

	// SaveNewRateUsdEther(bool)
	// SaveRateUSDEther(string)
	// GetIsNewRateUsdEther() bool
	// GetRateUSDEther() string

	// GetEvent() []ethereum.EventHistory
	// SaveEvent(*[]ethereum.EventHistory) error
	// GetIsNewEvent() bool
	// SetNewEvents(bool)

	GetLatestBlock() string
	GetIsNewLatestBlock() bool
	SaveLatestBlock(string) error
	SetNewLatestBlock(bool)

	GetRateUSD() []RateUSD
	GetRateETH() string
	GetIsNewRateUSD() bool
	SaveRateUSD(string) error
	SetNewRateUSD(bool)

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

	SaveMarketData(map[string]*ethereum.Rates)
	GetMarketData(page, pageSize uint64) map[string]*ethereum.MarketInfo
	SetIsNewMarketInfo(isNewMarketInfo bool)
	GetIsNewMarketInfo() bool
}

//var transactionPersistent = models.NewTransactionPersister()

func NewPersister(name string) (Persister, error) {
	Persister, err := NewRamPersister()
	return Persister, err
}
