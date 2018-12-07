package main

import (
	"log"
	"os"
	"runtime"
	"time"

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

	persisterIns, _ := persister.NewPersister("ram")
	fertcherIns, err := fetcher.NewFetcher()
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

	tokenNum := fertcherIns.GetNumTokens()
	bonusTimeWait := 900
	if tokenNum >= 200 {
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
	runFetchData(persisterIns, fetchRate, fertcherIns, 20)
	// runFetchData(persisterIns, fetchEvent, fertcherIns, 30)
	//runFetchData(persisterIns, fetchKyberEnable, fertcherIns, 10)

	runFetchData(persisterIns, fetchTrackerData, fertcherIns, 300)

	//run server
	server := http.NewHTTPServer(":3001", persisterIns, fertcherIns)
	server.Run()

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

func fetchRate(persister persister.Persister, fetcher *fetcher.Fetcher) {
	currentRate := persister.GetRate()
	rates, err := fetcher.GetRate(currentRate)
	if err != nil {
		log.Print(err)
		//persister.SetNewRate(false)
		return
	}
	persister.SaveRate(rates)
	//	persister.SetNewRate(true)
}

// func fetchEvent(persister persister.Persister, fetcher *fetcher.Fetcher) {
// 	if persister.GetIsNewLatestBlock() {
// 		blockNum := persister.GetLatestBlock()
// 		events, err := fetcher.GetEvents(blockNum)
// 		if err != nil {
// 			log.Print(err)
// 			persister.SetNewEvents(false)
// 			return
// 		}
// 		persister.SaveEvent(events)
// 		persister.SetNewEvents(true)
// 	} else {
// 		persister.SetNewEvents(false)
// 	}
// }

func fetchGeneralInfoTokens(persister persister.Persister, fetcher *fetcher.Fetcher) {
	generalInfo := fetcher.GetGeneralInfoTokens()
	persister.SaveGeneralInfoTokens(generalInfo)

	// if err != nil {
	// 	log.Print(err)
	// 	persister.SetIsNewGeneralInfoTokens(false)
	// 	return
	// }
	// persister.SaveGeneralInfoTokens(generalInfo)
	// persister.SetIsNewGeneralInfoTokens(true)
}

// func fetchKyberEnable(persister persister.Persister, fetcher *fetcher.Fetcher) {
// 	enable, err := fetcher.GetKyberEnable()
// 	if err != nil {
// 		log.Print(err)
// 		persister.SetNewKyberEnable(false)
// 		return
// 	}
// 	persister.SaveKyberEnable(enable)
// 	persister.SetNewKyberEnable(true)
// }

func fetchTrackerData(persister persister.Persister, fetcher *fetcher.Fetcher) {
	data, err := fetcher.FetchTrackerData()
	if err != nil {
		log.Print(err)
		persister.SetIsNewTrackerData(false)
		// return
	}
	tokens := fetcher.GetListToken()
	persister.SaveMarketData(data, tokens)
	persister.SetIsNewMarketInfo(true)
}
