package Persister

import (
	"errors"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/KyberNetwork/server-go/ethereum"
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
	mu sync.RWMutex

	kyberEnabled          bool
	isNewKyberEnabled     bool
	timeUpdateKyberEnable int64

	rates          *[]ethereum.Rate
	isNewRate      bool
	timeUpdateRate int64

	latestBlock           string
	isNewLatestBlock      bool
	timeUpdateLatestBlock int64

	rateUSD           []RateUSD
	rateETH           string
	isNewRateUsd      bool
	timeUpdateRateUSD int64

	// events     []ethereum.EventHistory
	// isNewEvent bool

	maxGasPrice      string
	isNewMaxGasPrice bool
	timeUpdateMaxGas int64

	gasPrice           *ethereum.GasPrice
	isNewGasPrice      bool
	timeUpdateGasPrice int64

	// ethRate      string
	// isNewEthRate bool

	tokenInfo           map[string]*ethereum.TokenGeneralInfo
	timeUpdateTokenInfo int64

	//isNewTokenInfo bool

	marketInfo           map[string]*ethereum.MarketInfo
	last7D               map[string][]float64
	rightMarketInfo      map[string]*ethereum.RightMarketInfo
	isNewMarketInfo      bool
	timeUpdateMarketInfo int64
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
	rateETH := "0"
	isNewRateUsd := true

	// events := make([]ethereum.EventHistory, 0)
	// isNewEvent := true

	maxGasPrice := "50"
	isNewMaxGasPrice := true

	gasPrice := ethereum.GasPrice{}
	isNewGasPrice := true

	// ethRate := "0"
	// isNewEthRate := true

	tokenInfo := map[string]*ethereum.TokenGeneralInfo{}
	//isNewTokenInfo := true

	marketInfo := map[string]*ethereum.MarketInfo{}
	last7D := map[string][]float64{}
	rightMarketInfo := map[string]*ethereum.RightMarketInfo{}
	isNewMarketInfo := true

	persister := &RamPersister{
		mu, kyberEnabled, isNewKyberEnabled, 0, &rates, isNewRate, 0, latestBlock, isNewLatestBlock, 0, rateUSD, rateETH, isNewRateUsd, 0, maxGasPrice, isNewMaxGasPrice, 0,
		&gasPrice, isNewGasPrice, 0, tokenInfo, 0, marketInfo, last7D, rightMarketInfo, isNewMarketInfo, 0,
	}
	return persister, nil
}

func (self *RamPersister) SaveGeneralInfoTokens(generalInfo map[string]*ethereum.TokenGeneralInfo) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.tokenInfo = generalInfo
	self.timeUpdateTokenInfo = time.Now().Unix()
}

/////------------------------------ Rates
func (self *RamPersister) GetRate() *[]ethereum.Rate {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.rates
}

func (self *RamPersister) SaveRate(rates *[]ethereum.Rate) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.rates = rates
	self.isNewRate = true
	self.timeUpdateRate = time.Now().Unix()
}

func (self *RamPersister) SetIsNewRate(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewRate = isNew
}

func (self *RamPersister) GetIsNewRate() bool {
	self.mu.Lock()
	defer self.mu.Unlock()
	if (self.timeUpdateRate + INTERVAL_UPDATE_GET_RATE) < time.Now().Unix() {
		return false
	}
	return self.isNewRate
}

//--------------------------------------------------------
func (self *RamPersister) SaveKyberEnabled(enabled bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.kyberEnabled = enabled
	self.isNewKyberEnabled = true
	self.timeUpdateKyberEnable = time.Now().Unix()
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
	if (self.timeUpdateKyberEnable + INTERVAL_UPDATE_KYBER_ENABLE) < time.Now().Unix() {
		return false
	}
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
	self.timeUpdateMaxGas = time.Now().Unix()
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
	if (self.timeUpdateMaxGas + INTERVAL_UPDATE_MAX_GAS) < time.Now().Unix() {
		return false
	}
	return self.isNewMaxGasPrice
}

//--------------------------------------------------------

//--------------------------------------------------------------

func (self *RamPersister) SaveGasPrice(gasPrice *ethereum.GasPrice) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.gasPrice = gasPrice
	self.isNewGasPrice = true
	self.timeUpdateGasPrice = time.Now().Unix()
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
	if (self.timeUpdateGasPrice + INTERVAL_UPDATE_GAS) < time.Now().Unix() {
		return false
	}
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
	if (self.timeUpdateRateUSD + INTERVAL_UPDATE_RATE_USD) < time.Now().Unix() {
		return false
	}
	return self.isNewRateUsd
}

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
	self.rateETH = rateUSDEth
	self.isNewRateUsd = true
	self.timeUpdateRateUSD = time.Now().Unix()
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
	self.timeUpdateLatestBlock = time.Now().Unix()
	return nil
}

func (self *RamPersister) GetIsNewLatestBlock() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	if (self.timeUpdateLatestBlock + INTERVAL_UPDATE_GET_BLOCKNUM) < time.Now().Unix() {
		return false
	}
	return self.isNewLatestBlock
}

func (self *RamPersister) SetNewLatestBlock(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewLatestBlock = isNew
}

/////-----------------------------------
// func (self *RamPersister) GetEvent() []ethereum.EventHistory {
// 	self.mu.RLock()
// 	defer self.mu.RUnlock()
// 	return self.events
// }

// func (self *RamPersister) GetIsNewEvent() bool {
// 	self.mu.RLock()
// 	defer self.mu.RUnlock()
// 	return self.isNewEvent
// }

// func (self *RamPersister) SaveEvent(events *[]ethereum.EventHistory) error {
// 	self.mu.Lock()
// 	defer self.mu.Unlock()
// 	self.events = *events
// 	return nil
// }

// func (self *RamPersister) SetNewEvents(isNew bool) {
// 	self.mu.Lock()
// 	defer self.mu.Unlock()
// 	self.isNewEvent = isNew
// }

// ----------------------------------------
// return data from kyber tracker

// use this api for 3 infomations change, marketcap, volume
func (self *RamPersister) GetRightMarketData() map[string]*ethereum.RightMarketInfo {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.rightMarketInfo
}

func (self *RamPersister) GetLast7D(listTokens string) map[string][]float64 {
	self.mu.Lock()
	defer self.mu.Unlock()
	tokens := strings.Split(listTokens, "-")
	result := make(map[string][]float64)
	for _, symbol := range tokens {
		if self.last7D[symbol] != nil {
			result[symbol] = self.last7D[symbol]
		}
	}
	return result
}

func (self *RamPersister) SaveMarketData(marketRate map[string]*ethereum.Rates, tokens map[string]ethereum.Token) {
	self.mu.Lock()
	defer self.mu.Unlock()
	// resultMarketInfo := []map[string]*ethereum.MarketInfo{}
	result := map[string]*ethereum.MarketInfo{}
	lastSevenDays := map[string][]float64{}
	newResult := map[string]*ethereum.RightMarketInfo{}

	tokenInfo := self.tokenInfo
	if (self.timeUpdateTokenInfo + INTERVAL_UPDATE_GENERAL_TOKEN_INFO) > time.Now().Unix() {
		tokenInfo = map[string]*ethereum.TokenGeneralInfo{}
	}

	for symbol, _ := range tokens {
		marketInfo := &ethereum.MarketInfo{}
		dataSevenDays := []float64{}
		rightMarketInfo := &ethereum.RightMarketInfo{}
		if rateInfo := marketRate[symbol]; rateInfo != nil {
			marketInfo.Rates = rateInfo
			dataSevenDays = rateInfo.P
			rightMarketInfo.Rate = &rateInfo.R
		}
		if tokenInfo := tokenInfo[symbol]; tokenInfo != nil {
			marketInfo.Quotes = tokenInfo.Quotes
			rightMarketInfo.Quotes = tokenInfo.Quotes
		}
		if marketInfo.Rates == nil && marketInfo.Quotes == nil {
			continue
		}
		result[symbol] = marketInfo
		newResult[symbol] = rightMarketInfo
		lastSevenDays[symbol] = dataSevenDays
	}

	self.marketInfo = result
	self.last7D = lastSevenDays
	self.rightMarketInfo = newResult
	self.timeUpdateMarketInfo = time.Now().Unix()
}

func (self *RamPersister) SetIsNewMarketInfo(isNewMarketInfo bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewMarketInfo = isNewMarketInfo
}

func (self *RamPersister) GetIsNewMarketInfo() bool {
	self.mu.Lock()
	defer self.mu.Unlock()
	if (self.timeUpdateMarketInfo + INTERVAL_UPDATE_DATA_TRACKER) < time.Now().Unix() {
		return false
	}
	return self.isNewMarketInfo
}
