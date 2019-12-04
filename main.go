package main

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/KyberNetwork/server-go/ethereum"
	"github.com/KyberNetwork/server-go/fetcher"
	"github.com/KyberNetwork/server-go/http"
	persister "github.com/KyberNetwork/server-go/persister"
)

type fetcherFunc func(persister persister.Persister, fetcher *fetcher.Fetcher)

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	//set log for server
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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

	//run fetch data
	runFetchData(persisterIns, fetchKyberEnabled, fertcherIns, 10)
	runFetchData(persisterIns, fetchMaxGasPrice, fertcherIns, 60)

	runFetchData(persisterIns, fetchGasPrice, fertcherIns, 30)

	runFetchData(persisterIns, fetchRateUSD, fertcherIns, 300)

	runFetchData(persisterIns, fetchBlockNumber, fertcherIns, 10)

	go fetchRate(persisterIns, fertcherIns)

	server := http.NewHTTPServer(":3001", persisterIns, fertcherIns)
	server.Run(kyberENV)
}

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
	ticker := time.NewTicker(15 * time.Second)
	for {
		var result []ethereum.Rate

		result, err := fetcher.FetchRate()
		if err != nil {
			log.Print(err)
			persister.SetIsNewRate(false)
			<-ticker.C
			continue
		}

		timeNow := time.Now().UTC().Unix()
		persister.SaveRate(result, timeNow)
		persister.SetIsNewRate(true)
		<-ticker.C
	}
}
