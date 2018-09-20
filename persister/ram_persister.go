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

	kyberEnabled      bool
	isNewKyberEnabled bool

	rates     *[]ethereum.Rate
	isNewRate bool

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

	// ethRate      string
	// isNewEthRate bool

	tokenInfo map[string]*ethereum.TokenGeneralInfo

	//isNewTokenInfo bool

	marketInfo      map[string]*ethereum.MarketInfo
	last7D          map[string][]float64
	rightMarketInfo map[string]*ethereum.RightMarketInfo
	isNewMarketInfo bool
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

	events := make([]ethereum.EventHistory, 0)
	isNewEvent := true

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
		mu, kyberEnabled, isNewKyberEnabled, &rates, isNewRate, latestBlock, isNewLatestBlock, rateUSD, rateETH, isNewRateUsd, events, isNewEvent, maxGasPrice, isNewMaxGasPrice,
		&gasPrice, isNewGasPrice, tokenInfo, marketInfo, last7D, rightMarketInfo, isNewMarketInfo,
	}
	return persister, nil
}

// func (self *RamPersister) SetRateToken(tokens map[string]ethereum.Token) {
// 	self.muRate.Lock()
// 	defer self.muRate.Unlock()
// 	for key, _ := range tokens {
// 		if key == "ETH" {
// 			continue
// 		}
// 		historyMap := map[int64]*ethereum.RateHistory{}
// 		rateInfo := ethereum.RateInfo{
// 			HistoryRecord: historyMap,
// 		}
// 		self.tokenRates[key] = &rateInfo
// 	}
// }

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
// func (self *RamPersister) SaveRateUSDEther(rate string) {
// 	self.mu.Lock()
// 	defer self.mu.Unlock()
// 	self.ethRate = rate
// }
// func (self *RamPersister) SaveNewRateUsdEther(isNew bool) {
// 	self.mu.Lock()
// 	defer self.mu.Unlock()
// 	self.isNewEthRate = isNew
// }

// func (self *RamPersister) GetIsNewRateUsdEther() bool {
// 	self.mu.RLock()
// 	defer self.mu.RUnlock()
// 	return self.isNewEthRate
// }
// func (self *RamPersister) GetRateUSDEther() string {
// 	self.mu.RLock()
// 	defer self.mu.RUnlock()
// 	return self.ethRate
// }

/////------------------------------
func (self *RamPersister) GetRate() *[]ethereum.Rate {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.rates
}

// func (self *RamPersister) ReStructureRate(rates *[]ethereum.Rate) map[string]*ethereum.RateHistory {
// 	refinedRate := map[string]*ethereum.RateHistory{}
// 	for key, _ := range self.tokenRates {
// 		refinedRate[key] = &ethereum.RateHistory{}
// 	}

// 	for _, rate := range *rates {
// 		if rate.Source != "ETH" {
// 			refinedRate[rate.Source].SellPrice = rate.Rate
// 		}
// 		if rate.Dest != "ETH" {
// 			refinedRate[rate.Dest].BuyPrice = rate.Rate
// 		}
// 	}
// 	return refinedRate
// }

func (self *RamPersister) SaveRate(rates *[]ethereum.Rate) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.rates = rates
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
		if tokenInfo := self.tokenInfo[symbol]; tokenInfo != nil {
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
}

func (self *RamPersister) SetIsNewMarketInfo(isNewMarketInfo bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewMarketInfo = isNewMarketInfo
}

func (self *RamPersister) GetIsNewMarketInfo() bool {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.isNewMarketInfo
}
