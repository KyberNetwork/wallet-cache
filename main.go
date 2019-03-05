package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/KyberNetwork/server-go/common"
	"github.com/KyberNetwork/server-go/ethereum"
	"github.com/KyberNetwork/server-go/fetcher"
	"github.com/KyberNetwork/server-go/http"
	persister "github.com/KyberNetwork/server-go/persister"
)

type fetcherFunc func(persister persister.Persister, fetcher *fetcher.Fetcher)

func enableLogToFile() (*os.File, error) {
	const logFileName = "error.log"
	f, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	//clear error log file
	if err = f.Truncate(0); err != nil {
		log.Fatal(err)
	}

	log.SetOutput(f)
	return f, nil
}

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	//set log for server
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if os.Getenv("LOG_TO_STDOUT") != "true" {
		f, err := enableLogToFile()
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}

	kyberENV := os.Getenv("KYBER_ENV")
	persisterIns, _ := persister.NewPersister("ram")
	fertcherIns, err := fetcher.NewFetcher(kyberENV)
	if err != nil {
		log.Fatal(err)
	}

	err = fertcherIns.TryUpdateListToken()
	if err != nil {
		log.Println(err)
	}

	tickerUpdateToken := time.NewTicker(300 * time.Second)
	go func() {
		for {
			<-tickerUpdateToken.C
			fertcherIns.TryUpdateListToken()
		}
	}()
	var (
		initRate  []ethereum.Rate
		ethSymbol = common.ETHSymbol
	)
	for symbol := range fertcherIns.GetListToken() {
		if symbol == ethSymbol {
			ethRate := ethereum.Rate{
				Source:  ethSymbol,
				Dest:    ethSymbol,
				Rate:    "0",
				Minrate: "0",
			}
			initRate = append(initRate, ethRate, ethRate)
		} else {
			buyRate := ethereum.Rate{
				Source:  ethSymbol,
				Dest:    symbol,
				Rate:    "0",
				Minrate: "0",
			}
			sellRate := ethereum.Rate{
				Source:  symbol,
				Dest:    ethSymbol,
				Rate:    "0",
				Minrate: "0",
			}
			initRate = append(initRate, buyRate, sellRate)
		}
	}
	persisterIns.SaveRate(initRate, 0)
	tokenNum := fertcherIns.GetNumTokens()
	bonusTimeWait := 900
	if tokenNum > 200 {
		bonusTimeWait = 60
	}
	intervalFetchGeneralInfoTokens := time.Duration((tokenNum * 7) + bonusTimeWait)
	//	initRateToken(persisterIns, fertcherIns)

	//run fetch data
	runFetchData(persisterIns, fetchKyberEnabled, fertcherIns, 10)
	runFetchData(persisterIns, fetchMaxGasPrice, fertcherIns, 60)

	runFetchData(persisterIns, fetchGasPrice, fertcherIns, 30)

	runFetchData(persisterIns, fetchRateUSD, fertcherIns, 600)

	//runFetchData(persisterIns, fetchRateUSDEther, fertcherIns, 600)

	runFetchData(persisterIns, fetchGeneralInfoTokens, fertcherIns, intervalFetchGeneralInfoTokens)

	runFetchData(persisterIns, fetchBlockNumber, fertcherIns, 10)
	runFetchData(persisterIns, fetchRate, fertcherIns, 30)
	runFetchData(persisterIns, fetchRateWithFallback, fertcherIns, 600)
	// runFetchData(persisterIns, fetchEvent, fertcherIns, 30)
	//runFetchData(persisterIns, fetchKyberEnable, fertcherIns, 10)

	runFetchData(persisterIns, fetchRate7dData, fertcherIns, 300)

	//run server
	server := http.NewHTTPServer(":3001", persisterIns, fertcherIns)
	server.Run(kyberENV)

	//init fetch data

}

// func setLogServer() {
// 	log.SetFlags(log.LstdFlags | log.Lshortfile)
// 	f, err := os.OpenFile("error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()
// 	log.SetOutput(f)
// }

// func initRateToken(persister persister.Persister, fertcherIns *fetcher.Fetcher) {
// 	tokens := fertcherIns.GetListToken()
// 	persister.SetRateToken(tokens)
// }

func runFetchData(persister persister.Persister, fn fetcherFunc, fertcherIns *fetcher.Fetcher, interval time.Duration) {
	ticker := time.NewTicker(interval * time.Second)
	go func() {
		for {
			fn(persister, fertcherIns)
			<-ticker.C
		}
	}()
}

func fetchGasPrice(persister persister.Persister, fetcher *fetcher.Fetcher) {
	gasPrice, err := fetcher.GetGasPrice()
	if err != nil {
		log.Print(err)
		persister.SetNewGasPrice(false)
		return
	}
	persister.SaveGasPrice(gasPrice)
}

func fetchMaxGasPrice(persister persister.Persister, fetcher *fetcher.Fetcher) {
	gasPrice, err := fetcher.GetMaxGasPrice()
	if err != nil {
		log.Print(err)
		persister.SetNewMaxGasPrice(false)
		return
	}
	persister.SaveMaxGasPrice(gasPrice)
}

func fetchKyberEnabled(persister persister.Persister, fetcher *fetcher.Fetcher) {
	enabled, err := fetcher.CheckKyberEnable()
	if err != nil {
		log.Print(err)
		persister.SetNewKyberEnabled(false)
		return
	}
	persister.SaveKyberEnabled(enabled)
}

func fetchRateUSD(persister persister.Persister, fetcher *fetcher.Fetcher) {
	rateUSD, err := fetcher.GetRateUsdEther()
	if err != nil {
		log.Print(err)
		persister.SetNewRateUSD(false)
		return
	}

	// if rateUSDCG == "" {
	// 	persister.SetNewRateUSDCG(false)
	// 	return
	// }

	if rateUSD == "" {
		persister.SetNewRateUSD(false)
		return
	}

	err = persister.SaveRateUSD(rateUSD)
	if err != nil {
		log.Print(err)
		persister.SetNewRateUSD(false)
		return
	}
}

// func fetchRateUSDEther(persister persister.Persister, fetcher *fetcher.Fetcher) {
// 	rateUSD, err := fetcher.GetRateUsdEther()
// 	if err != nil {
// 		log.Print(err)
// 		persister.SaveNewRateUsdEther(false)
// 		return
// 	}
// 	persister.SaveRateUSDEther(rateUSD)
// }

func fetchBlockNumber(persister persister.Persister, fetcher *fetcher.Fetcher) {
	blockNum, err := fetcher.GetLatestBlock()
	if err != nil {
		log.Print(err)
		persister.SetNewLatestBlock(false)
		return
	}
	err = persister.SaveLatestBlock(blockNum)
	if err != nil {
		persister.SetNewLatestBlock(false)
		log.Print(err)
		return
	}
}

func makeMapRate(rates []ethereum.Rate) map[string]ethereum.Rate {
	mapRate := make(map[string]ethereum.Rate)
	for _, r := range rates {
		mapRate[fmt.Sprintf("%s_%s", r.Source, r.Dest)] = r
	}
	return mapRate
}

func fetchRate(persister persister.Persister, fetcher *fetcher.Fetcher) {
	var result []ethereum.Rate
	currentRate := persister.GetRate()
	tokenPriority := fetcher.GetListTokenPriority()
	rates, err := fetcher.GetRate(currentRate, persister.GetIsNewRate(), tokenPriority, false)
	if err != nil {
		log.Print(err)
		persister.SetIsNewRate(false)
		return
	}
	mapRate := makeMapRate(rates)
	for _, cr := range currentRate {
		keyRate := fmt.Sprintf("%s_%s", cr.Source, cr.Dest)
		if r, ok := mapRate[keyRate]; ok {
			result = append(result, r)
			delete(mapRate, keyRate)
		} else {
			result = append(result, cr)
		}
	}
	// add new token to current rate
	if len(mapRate) > 0 {
		for _, nr := range mapRate {
			result = append(result, nr)
		}
	}
	timeNow := time.Now().UTC().Unix()
	persister.SaveRate(result, timeNow)
	persister.SetIsNewRate(true)
}

func fetchRateWithFallback(persister persister.Persister, fetcher *fetcher.Fetcher) {
	var result []ethereum.Rate
	currentRate := persister.GetRate()
	listToken := fetcher.GetListToken()
	newList := make(map[string]ethereum.Token)
	for _, t := range listToken {
		if !t.Priority {
			newList[t.Symbol] = t
		}
	}
	rates, err := fetcher.GetRate(currentRate, persister.GetIsNewRate(), newList, true)
	if err != nil {
		log.Print(err)
		persister.SetIsNewRate(false)
		return
	}
	mapRate := makeMapRate(rates)
	for _, cr := range currentRate {
		keyRate := fmt.Sprintf("%s_%s", cr.Source, cr.Dest)
		if r, ok := mapRate[keyRate]; ok {
			result = append(result, r)
			if keyRate != "ETH_ETH" {
				delete(mapRate, keyRate)
			}
		} else {
			result = append(result, cr)
		}
	}
	// add new token to current rate
	if len(mapRate) > 1 {
		for _, nr := range mapRate {
			result = append(result, nr)
		}
	}
	persister.SaveRate(result, 0)
	// persister.SetIsNewRate(true)
}

func fetchGeneralInfoTokens(persister persister.Persister, fetcher *fetcher.Fetcher) {
	generalInfo := fetcher.GetGeneralInfoTokens()
	persister.SaveGeneralInfoTokens(generalInfo)
}

func fetchRate7dData(persister persister.Persister, fetcher *fetcher.Fetcher) {
	data, err := fetcher.FetchRate7dData()
	if err != nil {
		log.Print(err)
		if !persister.IsFailedToFetchTracker() {
			return
		}
		persister.SetIsNewTrackerData(false)
	} else {
		persister.SetIsNewTrackerData(true)
	}
	tokens := fetcher.GetListToken()
	persister.SaveMarketData(data, tokens)
	// persister.SetIsNewMarketInfo(true)
}
