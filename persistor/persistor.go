package persistor

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

type Persistor interface {
	GetRate() *[]ethereum.Rate
	GetEvent() []ethereum.EventHistory
	GetLatestBlock() string
	GetRateUSD() []RateUSD

	GetIsNewRate() bool
	GetIsNewLatestBlock() bool
	GetIsNewRateUSD() bool
	GetIsNewEvent() bool

	SaveRate(*[]ethereum.Rate) error
	SaveEvent(*[]ethereum.EventHistory) error
	SaveLatestBlock(string) error
	SaveRateUSD([]*io.ReadCloser) error

	SetNewRate(bool)
	SetNewLatestBlock(bool)
	SetNewRateUSD(bool)
	SetNewEvents(bool)
}

//var transactionPersistent = models.NewTransactionPersister()

func NewPersistor(name string) (Persistor, error) {
	persistor, err := NewRamPersistor()
	return persistor, err
}
