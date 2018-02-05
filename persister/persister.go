package Persister

import (
	"io"

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
	GetRate() *[]ethereum.Rate
	SaveRate(*[]ethereum.Rate) error
	SetNewRate(bool)
	GetIsNewRate() bool

	GetEvent() []ethereum.EventHistory
	SaveEvent(*[]ethereum.EventHistory) error
	GetIsNewEvent() bool
	SetNewEvents(bool)

	GetLatestBlock() string
	GetIsNewLatestBlock() bool
	SaveLatestBlock(string) error
	SetNewLatestBlock(bool)

	GetRateUSD() []RateUSD
	GetIsNewRateUSD() bool
	SaveRateUSD([]io.ReadCloser) error
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
}

//var transactionPersistent = models.NewTransactionPersister()

func NewPersister(name string) (Persister, error) {
	Persister, err := NewRamPersister()
	return Persister, err
}
