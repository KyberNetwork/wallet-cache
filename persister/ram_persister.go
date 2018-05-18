package Persister

import (
	"sync"
	"time"

	"github.com/KyberNetwork/server-go/ethereum"
)

const STEP_SAVE_RATE = 10      //1 minute
const MAXIMUM_SAVE_RECORD = 60 //60 records

type RamPersister struct {
	mu     sync.RWMutex
	muRate sync.RWMutex

	kyberEnabled      bool
	isNewKyberEnabled bool

	rates     *[]ethereum.Rate
	isNewRate bool

	// latestBlock      string
	// isNewLatestBlock bool

	// rateUSD      []RateUSD
	// isNewRateUsd bool

	// events     []ethereum.EventHistory
	// isNewEvent bool

	maxGasPrice      string
	isNewMaxGasPrice bool

	gasPrice      *ethereum.GasPrice
	isNewGasPrice bool

	ethRate      string
	isNewEthRate bool

	tokenInfo map[string]*ethereum.TokenGeneralInfo

	tokenRates        map[string]*ethereum.RateInfo
	lastimeUpdateRate int64
	timeStampArr      []int64
	//isNewTokenInfo bool
}

func NewRamPersister() (*RamPersister, error) {
	var mu sync.RWMutex
	var muRate sync.RWMutex

	kyberEnabled := true
	isNewKyberEnabled := true

	rates := make([]ethereum.Rate, 0)
	isNewRate := true

	// latestBlock := "0"
	// isNewLatestBlock := true

	// rateUSD := make([]RateUSD, 0)
	// isNewRateUsd := true

	// events := make([]ethereum.EventHistory, 0)
	// isNewEvent := true

	maxGasPrice := "50"
	isNewMaxGasPrice := true

	gasPrice := ethereum.GasPrice{}
	isNewGasPrice := true

	ethRate := "0"
	isNewEthRate := true

	tokenInfo := map[string]*ethereum.TokenGeneralInfo{}
	//isNewTokenInfo := true

	tokenRates := map[string]*ethereum.RateInfo{}
	lastimeUpdateRate := time.Now().Unix()
	timeStampArr := make([]int64, 0)

	persister := &RamPersister{
		mu, muRate, kyberEnabled, isNewKyberEnabled, &rates, isNewRate, maxGasPrice, isNewMaxGasPrice,
		&gasPrice, isNewGasPrice, ethRate, isNewEthRate, tokenInfo, tokenRates, lastimeUpdateRate, timeStampArr,
	}
	return persister, nil
}

func (self *RamPersister) SetRateToken(tokens map[string]ethereum.Token) {
	self.muRate.Lock()
	defer self.muRate.Unlock()
	for key, _ := range tokens {
		if key == "ETH" {
			continue
		}
		historyMap := map[int64]*ethereum.RateHistory{}
		rateInfo := ethereum.RateInfo{
			HistoryRecord: historyMap,
		}
		self.tokenRates[key] = &rateInfo
	}
}

func (self *RamPersister) SaveGeneralInfoTokens(generalInfo map[string]*ethereum.TokenGeneralInfo) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.tokenInfo = generalInfo
}

func (self *RamPersister) GetTokenInfo() map[string]*ethereum.TokenGeneralInfo {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.tokenInfo
}

//----------------------------
func (self *RamPersister) SaveRateUSDEther(rate string) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.ethRate = rate
}
func (self *RamPersister) SaveNewRateUsdEther(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewEthRate = isNew
}

func (self *RamPersister) GetIsNewRateUsdEther() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewEthRate
}
func (self *RamPersister) GetRateUSDEther() string {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.ethRate
}

/////------------------------------
func (self *RamPersister) GetRate() map[string]*ethereum.RateInfo {
	self.muRate.RLock()
	defer self.muRate.RUnlock()
	return self.tokenRates
}

func (self *RamPersister) ReStructureRate(rates *[]ethereum.Rate) map[string]*ethereum.RateHistory {
	refinedRate := map[string]*ethereum.RateHistory{}
	for key, _ := range self.tokenRates {
		refinedRate[key] = &ethereum.RateHistory{}
	}

	for _, rate := range *rates {
		if rate.Source != "ETH" {
			refinedRate[rate.Source].SellPrice = rate.Rate
		}
		if rate.Dest != "ETH" {
			refinedRate[rate.Dest].BuyPrice = rate.Rate
		}
	}
	return refinedRate
}

func (self *RamPersister) SaveRate(rates *[]ethereum.Rate) {
	self.muRate.Lock()
	defer self.muRate.Unlock()
	self.rates = rates

	refinedRate := self.ReStructureRate(rates)
	//save rate token
	isSaveRateRecord := false
	currentTime := time.Now().Unix()
	if currentTime-self.lastimeUpdateRate > STEP_SAVE_RATE {
		isSaveRateRecord = true
		self.timeStampArr = append(self.timeStampArr, currentTime)
	}
	for key, val := range refinedRate {
		self.tokenRates[key].LastBuy = val.BuyPrice
		self.tokenRates[key].LastSell = val.SellPrice
		if isSaveRateRecord {
			self.tokenRates[key].HistoryRecord[currentTime] = val

			if len(self.tokenRates[key].HistoryRecord) > MAXIMUM_SAVE_RECORD {
				if _, ok := self.tokenRates[key].HistoryRecord[self.timeStampArr[0]]; ok {
					historyRecord := self.tokenRates[key].HistoryRecord
					rateInfo := ethereum.RateInfo{
						LastSell:      val.SellPrice,
						LastBuy:       val.BuyPrice,
						HistoryRecord: historyRecord,
					}
					self.tokenRates[key] = &rateInfo
				}

			}

		}
	}
	if len(self.timeStampArr) > MAXIMUM_SAVE_RECORD {
		self.timeStampArr = self.timeStampArr[1:]
	}

	if isSaveRateRecord {
		self.lastimeUpdateRate = currentTime
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
