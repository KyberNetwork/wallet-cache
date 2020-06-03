package persister

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/KyberNetwork/cache/ethereum"
)

const (
	STEP_SAVE_RATE      = 10 //1 minute
	MAXIMUM_SAVE_RECORD = 60 //60 records

	INTERVAL_UPDATE_KYBER_ENABLE       = 20
	INTERVAL_UPDATE_MAX_GAS            = 70
	INTERVAL_UPDATE_GAS                = 40
	INTERVAL_UPDATE_RATE_USD           = 610
	INTERVAL_UPDATE_GENERAL_TOKEN_INFO = 3600
	INTERVAL_UPDATE_GET_BLOCKNUM       = 20
	INTERVAL_UPDATE_GET_RATE           = 30
	INTERVAL_UPDATE_DATA_TRACKER       = 310
)

type RamPersister struct {
	mu      sync.RWMutex
	timeRun string

	kyberEnabled      bool
	isNewKyberEnabled bool

	rates     []ethereum.Rate
	isNewRate bool
	updatedAt int64

	latestBlock      string
	isNewLatestBlock bool

	rateUSD      []RateUSD
	rateETH      string
	isNewRateUsd bool

	events     []ethereum.EventHistory
	isNewEvent bool

	maxGasPrice      string
	isNewMaxGasPrice bool

	gasPrice      *ethereum.GasPrice
	isNewGasPrice bool
}

func NewRamPersister() (*RamPersister, error) {
	var mu sync.RWMutex
	location, _ := time.LoadLocation("Asia/Bangkok")
	tNow := time.Now().In(location)
	timeRun := fmt.Sprintf("%02d:%02d:%02d %02d-%02d-%d", tNow.Hour(), tNow.Minute(), tNow.Second(), tNow.Day(), tNow.Month(), tNow.Year())

	kyberEnabled := true
	isNewKyberEnabled := true

	rates := []ethereum.Rate{}
	isNewRate := false

	latestBlock := "0"
	isNewLatestBlock := true

	rateUSD := make([]RateUSD, 0)
	rateETH := "0"
	isNewRateUsd := true

	events := make([]ethereum.EventHistory, 0)
	isNewEvent := true

	maxGasPrice := "50"
	isNewMaxGasPrice := true

	gasPrice := ethereum.GasPrice{}
	isNewGasPrice := true

	persister := &RamPersister{
		mu:                mu,
		timeRun:           timeRun,
		kyberEnabled:      kyberEnabled,
		isNewKyberEnabled: isNewKyberEnabled,
		rates:             rates,
		isNewRate:         isNewRate,
		updatedAt:         0,
		latestBlock:       latestBlock,
		isNewLatestBlock:  isNewLatestBlock,
		rateUSD:           rateUSD,
		rateETH:           rateETH,
		isNewRateUsd:      isNewRateUsd,
		events:            events,
		isNewEvent:        isNewEvent,
		maxGasPrice:       maxGasPrice,
		isNewMaxGasPrice:  isNewMaxGasPrice,
		gasPrice:          &gasPrice,
		isNewGasPrice:     isNewGasPrice,
	}
	return persister, nil
}

/////------------------------------
func (self *RamPersister) GetRate() []ethereum.Rate {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.rates
}

func (self *RamPersister) GetTimeUpdateRate() int64 {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.updatedAt
}

func (self *RamPersister) SetIsNewRate(isNewRate bool) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	// return self.rates
	self.isNewRate = isNewRate
}

func (self *RamPersister) GetIsNewRate() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewRate
}

func (self *RamPersister) SaveRate(rates []ethereum.Rate, timestamp int64) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.rates = rates
	if timestamp != 0 {
		self.updatedAt = timestamp
	}
}

//--------------------------------------------------------
func (self *RamPersister) SaveKyberEnabled(enabled bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.kyberEnabled = enabled
	self.isNewKyberEnabled = true
}

func (self *RamPersister) SetNewKyberEnabled(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewKyberEnabled = isNew
}

func (self *RamPersister) GetKyberEnabled() bool {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.kyberEnabled
}

func (self *RamPersister) GetNewKyberEnabled() bool {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.isNewKyberEnabled
}

//--------------------------------------------------------

//--------------------------------------------------------

func (self *RamPersister) SetNewMaxGasPrice(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewMaxGasPrice = isNew
	return
}

func (self *RamPersister) SaveMaxGasPrice(maxGasPrice string) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.maxGasPrice = maxGasPrice
	self.isNewMaxGasPrice = true
	return
}
func (self *RamPersister) GetMaxGasPrice() string {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.maxGasPrice
}
func (self *RamPersister) GetNewMaxGasPrice() bool {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.isNewMaxGasPrice
}

//--------------------------------------------------------

//--------------------------------------------------------------

func (self *RamPersister) SaveGasPrice(gasPrice *ethereum.GasPrice) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.gasPrice = gasPrice
	self.isNewGasPrice = true
}
func (self *RamPersister) SetNewGasPrice(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewGasPrice = isNew
}
func (self *RamPersister) GetGasPrice() *ethereum.GasPrice {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.gasPrice
}
func (self *RamPersister) GetNewGasPrice() bool {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.isNewGasPrice
}

//-----------------------------------------------------------

func (self *RamPersister) GetRateUSD() []RateUSD {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.rateUSD
}

func (self *RamPersister) GetRateETH() string {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.rateETH
}

func (self *RamPersister) GetIsNewRateUSD() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewRateUsd
}

func (self *RamPersister) SaveRateUSD(rateUSDEth string) error {
	self.mu.Lock()
	defer self.mu.Unlock()

	rates := make([]RateUSD, 0)

	itemRateEth := RateUSD{Symbol: "ETH", PriceUsd: rateUSDEth}
	rates = append(rates, itemRateEth)
	for _, item := range self.rates {
		if item.Source != "ETH" {
			priceUsd, err := CalculateRateUSD(item.Rate, rateUSDEth)
			if err != nil {
				log.Print(err)
				self.isNewRateUsd = false
				return nil
			}
			sourceSymbol := item.Source
			if sourceSymbol == "ETHOS" {
				sourceSymbol = "BQX"
			}
			itemRate := RateUSD{Symbol: sourceSymbol, PriceUsd: priceUsd}
			rates = append(rates, itemRate)
		}
	}

	self.rateUSD = rates
	self.rateETH = rateUSDEth
	self.isNewRateUsd = true

	return nil
}

func CalculateRateUSD(rateEther string, rateUSD string) (string, error) {
	//func (z *Int) SetString(s string, base int) (*Int, bool)

	bigRateUSD, ok := new(big.Float).SetString(rateUSD)
	if !ok {
		err := errors.New("Cannot convert rate usd of ether to big float")
		return "", err
	}
	bigRateEth, ok := new(big.Float).SetString(rateEther)
	if !ok {
		err := errors.New("Cannot convert rate token-eth to big float")
		return "", err
	}
	i, e := big.NewInt(10), big.NewInt(18)
	i.Exp(i, e, nil)
	weight := new(big.Float).SetInt(i)

	rateUSDBig := new(big.Float).Mul(bigRateUSD, bigRateEth)
	rateUSDNormal := new(big.Float).Quo(rateUSDBig, weight)
	return rateUSDNormal.String(), nil
}

func (self *RamPersister) SetNewRateUSD(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewRateUsd = isNew
}

func (self *RamPersister) GetLatestBlock() string {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.latestBlock
}

func (self *RamPersister) SaveLatestBlock(blockNumber string) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.latestBlock = blockNumber
	self.isNewLatestBlock = true
	return nil
}

func (self *RamPersister) GetIsNewLatestBlock() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewLatestBlock
}

func (self *RamPersister) SetNewLatestBlock(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewLatestBlock = isNew
}

func (self *RamPersister) GetTimeVersion() string {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.timeRun
}
