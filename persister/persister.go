package persister

import (
	"github.com/KyberNetwork/cache/ethereum"
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

type MemoryPersister interface {
	GetRate() []ethereum.Rate
	GetIsNewRate() bool
	SetIsNewRate(bool)
	GetTimeUpdateRate() int64

	SaveRate([]ethereum.Rate, int64)

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

	GetTimeVersion() string
}

func NewMemoryPersister(name string) (MemoryPersister, error) {
	return NewRamPersister()
}

type DiskPersister interface {
	SaveGasPrice(gasOracle ethereum.GasPrice) error
	GetWeeklyAverageGasPrice() (float64, error)
}

func NewDiskPersister(name string) (DiskPersister, error) {
	switch name {
	case "leveldb":
		return NewLeveldbPersister()
	default:
		return NewLeveldbPersister()
	}
}
