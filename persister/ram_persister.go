package Persister

import (
	"errors"
	"log"
	"math/big"
	"sync"

	"github.com/KyberNetwork/server-go/ethereum"
)

type RamPersister struct {
	mu sync.RWMutex

	kyberEnabled      bool
	isNewKyberEnabled bool

	rates     *[]ethereum.Rate
	isNewRate bool

	latestBlock      string
	isNewLatestBlock bool

	rateUSD      []RateUSD
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

	kyberEnabled := true
	isNewKyberEnabled := true

	rates := make([]ethereum.Rate, 0)
	isNewRate := true

	latestBlock := "0"
	isNewLatestBlock := true

	rateUSD := make([]RateUSD, 0)
	isNewRateUsd := true

	events := make([]ethereum.EventHistory, 0)
	isNewEvent := true

	maxGasPrice := "50"
	isNewMaxGasPrice := true

	gasPrice := ethereum.GasPrice{}
	isNewGasPrice := true

	persister := &RamPersister{
		mu, kyberEnabled, isNewKyberEnabled, &rates, isNewRate, latestBlock, isNewLatestBlock, rateUSD, isNewRateUsd, events, isNewEvent, maxGasPrice, isNewMaxGasPrice,
		&gasPrice, isNewGasPrice,
	}
	return persister, nil
}

/////------------------------------
func (self *RamPersister) GetRate() *[]ethereum.Rate {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.rates
}
func (self *RamPersister) SaveRate(rates *[]ethereum.Rate) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.rates = rates
	return nil
}

func (self *RamPersister) GetIsNewRate() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewRate
}

func (self *RamPersister) SetNewRate(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewRate = isNew
}

/////-------------------------------

/////-----------------------------------
func (self *RamPersister) GetEvent() []ethereum.EventHistory {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.events
}

func (self *RamPersister) GetIsNewEvent() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewEvent
}

func (self *RamPersister) SaveEvent(events *[]ethereum.EventHistory) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.events = *events
	return nil
}

func (self *RamPersister) SetNewEvents(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewEvent = isNew
}

/////-------------------------------------------------

////--------------------------------------------------
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

////--------------------------------------------------

///----------------------------------------------------

func (self *RamPersister) GetRateUSD() []RateUSD {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.rateUSD
}

func (self *RamPersister) GetIsNewRateUSD() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewRateUsd
}

// func (self *RamPersister) SaveRateUSD(body []io.ReadCloser) error {
// 	rates := make([]RateUSD, 0)
// 	for _, item := range body {
// 		rateItem := make([]RateUSD, 0)
// 		defer (item).Close()
// 		b, err := ioutil.ReadAll(item)
// 		if err != nil {
// 			log.Print(err)
// 			return err
// 		}
// 		err = json.Unmarshal(b, &rateItem)
// 		if err != nil {
// 			log.Print(err)
// 			return err
// 		}
// 		if rateItem[0].Symbol == "ETHOS" {
// 			rateItem[0].Symbol = "BQX"
// 		}
// 		rates = append(rates, rateItem[0])
// 	}
// 	self.mu.Lock()
// 	defer self.mu.Unlock()
// 	self.rateUSD = rates
// 	self.isNewRateUsd = true
// 	return nil
// }

func (self *RamPersister) SaveRateUSD(rateUSDEth string) error {
	rates := make([]RateUSD, 0)

	itemRateEth := RateUSD{Symbol: "ETH", PriceUsd: rateUSDEth}
	rates = append(rates, itemRateEth)
	for _, item := range *(self.rates) {
		if item.Source != "ETH" {
			priceUsd, err := CalculateRateUSD(item.Rate, rateUSDEth)
			if err != nil {
				log.Print(err)
				return err
			}
			sourceSymbol := item.Source
			if sourceSymbol == "ETHOS" {
				sourceSymbol = "BQX"
			}
			itemRate := RateUSD{Symbol: sourceSymbol, PriceUsd: priceUsd}
			rates = append(rates, itemRate)
		}
	}

	self.mu.Lock()
	defer self.mu.Unlock()
	self.rateUSD = rates
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

///----------------------------------------------------

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
