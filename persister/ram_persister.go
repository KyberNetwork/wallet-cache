package Persister

import (
	"errors"
	"fmt"
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
	mu      sync.RWMutex
	timeRun string

	kyberEnabled      bool
	isNewKyberEnabled bool

	rates     *[]ethereum.Rate
	isNewRate bool

	latestBlock      string
	isNewLatestBlock bool

	rateUSD      []RateUSD
	rateETH      string
	isNewRateUsd bool

	// rateUSDCG      []RateUSD
	// rateETHCG      string
	// isNewRateUsdCG bool

	events     []ethereum.EventHistory
	isNewEvent bool

	maxGasPrice      string
	isNewMaxGasPrice bool

	gasPrice      *ethereum.GasPrice
	isNewGasPrice bool

	// ethRate      string
	// isNewEthRate bool

	tokenInfo map[string]*ethereum.TokenGeneralInfo
	// tokenInfoCG map[string]*ethereum.TokenGeneralInfo

	//isNewTokenInfo bool

	marketInfo       map[string]*ethereum.MarketInfo
	last7D           map[string][]float64
	isNewTrackerData bool

	rightMarketInfo map[string]*ethereum.RightMarketInfo
	// rightMarketInfoCG map[string]*ethereum.RightMarketInfo

	isNewMarketInfo bool
	// isNewMarketInfoCG bool
}

func NewRamPersister() (*RamPersister, error) {
	var mu sync.RWMutex
	location, _ := time.LoadLocation("Asia/Bangkok")
	tNow := time.Now().In(location)
	timeRun := fmt.Sprintf("%02d:%02d:%02d %02d-%02d-%d", tNow.Hour(), tNow.Minute(), tNow.Second(), tNow.Day(), tNow.Month(), tNow.Year())

	kyberEnabled := true
	isNewKyberEnabled := true

	rates := make([]ethereum.Rate, 0)
	isNewRate := true

	latestBlock := "0"
	isNewLatestBlock := true

	rateUSD := make([]RateUSD, 0)
	rateETH := "0"
	isNewRateUsd := true
	// rateUSDCG := make([]RateUSD, 0)
	// rateETHCG := "0"
	// isNewRateUsdCG := true

	events := make([]ethereum.EventHistory, 0)
	isNewEvent := true

	maxGasPrice := "50"
	isNewMaxGasPrice := true

	gasPrice := ethereum.GasPrice{}
	isNewGasPrice := true

	// ethRate := "0"
	// isNewEthRate := true

	tokenInfo := map[string]*ethereum.TokenGeneralInfo{}
	// tokenInfoCG := map[string]*ethereum.TokenGeneralInfo{}
	//isNewTokenInfo := true

	marketInfo := map[string]*ethereum.MarketInfo{}
	last7D := map[string][]float64{}
	isNewTrackerData := true

	rightMarketInfo := map[string]*ethereum.RightMarketInfo{}
	// rightMarketInfoCG := map[string]*ethereum.RightMarketInfo{}

	isNewMarketInfo := true
	// isNewMarketInfoCG := true

	persister := &RamPersister{
		mu:                mu,
		timeRun:           timeRun,
		kyberEnabled:      kyberEnabled,
		isNewKyberEnabled: isNewKyberEnabled,
		rates:             &rates,
		isNewRate:         isNewRate,
		latestBlock:       latestBlock,
		isNewLatestBlock:  isNewLatestBlock,
		rateUSD:           rateUSD,
		rateETH:           rateETH,
		isNewRateUsd:      isNewRateUsd,
		// rateUSDCG:         rateUSDCG,
		// rateETHCG:         rateETHCG,
		// isNewRateUsdCG:    isNewRateUsdCG,
		events:           events,
		isNewEvent:       isNewEvent,
		maxGasPrice:      maxGasPrice,
		isNewMaxGasPrice: isNewMaxGasPrice,
		gasPrice:         &gasPrice,
		isNewGasPrice:    isNewGasPrice,
		tokenInfo:        tokenInfo,
		// tokenInfoCG:       tokenInfoCG,
		marketInfo:       marketInfo,
		last7D:           last7D,
		isNewTrackerData: isNewTrackerData,
		rightMarketInfo:  rightMarketInfo,
		// rightMarketInfoCG: rightMarketInfoCG,
		isNewMarketInfo: isNewMarketInfo,
		// isNewMarketInfoCG: isNewMarketInfoCG,
	}
	return persister, nil
}

func (self *RamPersister) SaveGeneralInfoTokens(generalInfo map[string]*ethereum.TokenGeneralInfo) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.tokenInfo = generalInfo
	// self.tokenInfoCG = generalInfoCG
}

func (self *RamPersister) GetTokenInfo() map[string]*ethereum.TokenGeneralInfo {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.tokenInfo
}

/////------------------------------
func (self *RamPersister) GetRate() *[]ethereum.Rate {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.rates
}

func (self *RamPersister) SetIsNewRate(isNewRate bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	// return self.rates
	self.isNewRate = isNewRate
}

func (self *RamPersister) GetIsNewRate() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewRate
}

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

// func (self *RamPersister) GetRateUSDCG() []RateUSD {
// 	self.mu.RLock()
// 	defer self.mu.RUnlock()
// 	return self.rateUSDCG
// }

func (self *RamPersister) GetRateETH() string {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.rateETH
}

// func (self *RamPersister) GetRateETHCG() string {
// 	self.mu.RLock()
// 	defer self.mu.RUnlock()
// 	return self.rateETHCG
// }

func (self *RamPersister) GetIsNewRateUSD() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewRateUsd
}

// func (self *RamPersister) GetIsNewRateUSDCG() bool {
// 	self.mu.RLock()
// 	defer self.mu.RUnlock()
// 	return self.isNewRateUsdCG
// }

func (self *RamPersister) SaveRateUSD(rateUSDEth string) error {
	self.mu.Lock()
	defer self.mu.Unlock()

	rates := make([]RateUSD, 0)
	// ratesCG := make([]RateUSD, 0)

	itemRateEth := RateUSD{Symbol: "ETH", PriceUsd: rateUSDEth}
	// itemRateEthCG := RateUSD{Symbol: "ETH", PriceUsd: rateUSDEthCG}
	rates = append(rates, itemRateEth)
	// ratesCG = append(ratesCG, itemRateEthCG)
	for _, item := range *(self.rates) {
		if item.Source != "ETH" {
			priceUsd, err := CalculateRateUSD(item.Rate, rateUSDEth)
			if err != nil {
				log.Print(err)
				self.isNewRateUsd = false
				return nil
			}
			// priceUsdCG, err := CalculateRateUSD(item.Rate, rateUSDEthCG)
			// if err != nil {
			// 	log.Print(err)
			// 	self.isNewRateUsdCG = false
			// 	return nil
			// }
			sourceSymbol := item.Source
			if sourceSymbol == "ETHOS" {
				sourceSymbol = "BQX"
			}
			itemRate := RateUSD{Symbol: sourceSymbol, PriceUsd: priceUsd}
			rates = append(rates, itemRate)
			// ratesCG = append(ratesCG, RateUSD{
			// 	Symbol:   sourceSymbol,
			// 	PriceUsd: priceUsdCG,
			// })
		}
	}

	self.rateUSD = rates
	self.rateETH = rateUSDEth
	self.isNewRateUsd = true

	// self.rateUSDCG = ratesCG
	// self.rateETHCG = rateUSDEthCG
	// self.isNewRateUsdCG = true

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

// func (self *RamPersister) SetNewRateUSDCG(isNew bool) {
// 	self.mu.Lock()
// 	defer self.mu.Unlock()
// 	self.isNewRateUsdCG = isNew
// }

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

// ----------------------------------------
// return data from kyber tracker

// use this api for 3 infomations change, marketcap, volume
func (self *RamPersister) GetRightMarketData() map[string]*ethereum.RightMarketInfo {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.rightMarketInfo
}

// func (self *RamPersister) GetRightMarketDataCG() map[string]*ethereum.RightMarketInfo {
// 	self.mu.Lock()
// 	defer self.mu.Unlock()
// 	return self.rightMarketInfoCG
// }

func (self *RamPersister) GetIsNewTrackerData() bool {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.isNewTrackerData
}

func (self *RamPersister) SetIsNewTrackerData(isNewTrackerData bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewTrackerData = isNewTrackerData
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
	// result := map[string]*ethereum.MarketInfo{}
	lastSevenDays := map[string][]float64{}
	newResult := map[string]*ethereum.RightMarketInfo{}
	// newResultCG := map[string]*ethereum.RightMarketInfo{}

	for symbol, _ := range tokens {
		// marketInfo := &ethereum.MarketInfo{}
		dataSevenDays := []float64{}
		rightMarketInfo := &ethereum.RightMarketInfo{}
		// rightMarketInfoCG := &ethereum.RightMarketInfo{}
		rateInfo := marketRate[symbol]
		if rateInfo != nil {
			// marketInfo.Rates = rateInfo
			dataSevenDays = rateInfo.P
			rightMarketInfo.Rate = &rateInfo.R
			// rightMarketInfoCG.Rate = &rateInfo.R
		}
		if tokenInfo := self.tokenInfo[symbol]; tokenInfo != nil {
			// marketInfo.Quotes = tokenInfo.Quotes
			rightMarketInfo.Quotes = tokenInfo.Quotes
		}

		// if tokenInfoCG := self.tokenInfoCG[symbol]; tokenInfoCG != nil {
		// marketInfo.Quotes = tokenInfo.Quotes
		// rightMarketInfoCG.Quotes = tokenInfoCG.Quotes
		// }

		if rateInfo == nil && rightMarketInfo.Quotes == nil {
			// newResult[symbol] = rightMarketInfo
			// lastSevenDays[symbol] = dataSevenDays
			continue
		}

		// // if rightMarketInfoCG.Rate != nil && rightMarketInfoCG.Quotes != nil {
		// newResultCG[symbol] = rightMarketInfoCG
		// lastSevenDays[symbol] = dataSevenDays
		// }

		// result[symbol] = marketInfo
		newResult[symbol] = rightMarketInfo
		lastSevenDays[symbol] = dataSevenDays
	}

	// self.marketInfo = result
	self.last7D = lastSevenDays
	self.rightMarketInfo = newResult
	// self.rightMarketInfoCG = newResultCG
}

func (self *RamPersister) SetIsNewMarketInfo(isNewMarketInfo bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewMarketInfo = isNewMarketInfo
}

// func (self *RamPersister) SetIsNewMarketInfoCG(isNewMarketInfo bool) {
// 	self.mu.Lock()
// 	defer self.mu.Unlock()
// 	self.isNewMarketInfoCG = isNewMarketInfo
// }

func (self *RamPersister) GetIsNewMarketInfo() bool {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.isNewMarketInfo
}

// func (self *RamPersister) GetIsNewMarketInfoCG() bool {
// 	self.mu.Lock()
// 	defer self.mu.Unlock()
// 	return self.isNewMarketInfoCG
// }

func (self *RamPersister) GetTimeVersion() string {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.timeRun
}
