package main

import (
	"log"
	"os"
	"time"

	"github.com/KyberNetwork/server-go/fetcher"
	"github.com/KyberNetwork/server-go/http"
	"github.com/KyberNetwork/server-go/persistor"
)

type fetcherFunc func(persistor persistor.Persistor, fetcher *fetcher.Fetcher)

func main() {

	//set log for server
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	f, err := os.OpenFile("error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//clear error log file
	err = f.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	persistorIns, _ := persistor.NewPersistor("ram")
	fertcherIns, _ := fetcher.NewFetcher()

	//run fetch data
	runFetchData(persistorIns, fetchRateUSD, fertcherIns, 60)
	runFetchData(persistorIns, fetchBlockNumber, fertcherIns, 10)
	runFetchData(persistorIns, fetchRate, fertcherIns, 10)
	runFetchData(persistorIns, fetchEvent, fertcherIns, 30)

	//run server
	server := http.NewHTTPServer(":3002", persistorIns)
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

func runFetchData(persistor persistor.Persistor, fn fetcherFunc, fertcherIns *fetcher.Fetcher, interval time.Duration) {
	fn(persistor, fertcherIns)
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fn(persistor, fertcherIns)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func fetchRateUSD(persistor persistor.Persistor, fetcher *fetcher.Fetcher) {
	body, err := fetcher.GetRateUsd()
	if err != nil {
		log.Print(err)
		persistor.SetNewRateUSD(false)
		return
	}
	err = persistor.SaveRateUSD(body)
	if err != nil {
		log.Print(err)
		persistor.SetNewRateUSD(false)
		return
	}
}

func fetchBlockNumber(persistor persistor.Persistor, fetcher *fetcher.Fetcher) {
	blockNum, err := fetcher.GetLatestBlock()
	if err != nil {
		log.Print(err)
		persistor.SetNewLatestBlock(false)
		return
	}
	err = persistor.SaveLatestBlock(blockNum)
	if err != nil {
		persistor.SetNewLatestBlock(false)
		log.Print(err)
		return
	}
}

func fetchRate(persistor persistor.Persistor, fetcher *fetcher.Fetcher) {
	rates, err := fetcher.GetRate()
	if err != nil {
		log.Print(err)
		persistor.SetNewRate(false)
		return
	}
	persistor.SaveRate(rates)
}

func fetchEvent(persistor persistor.Persistor, fetcher *fetcher.Fetcher) {
	if persistor.GetIsNewLatestBlock() {
		blockNum := persistor.GetLatestBlock()
		events, err := fetcher.GetEvents(blockNum)
		if err != nil {
			log.Print(err)
			persistor.SetNewEvents(false)
			return
		}
		persistor.SaveEvent(events)
	} else {
		persistor.SetNewEvents(false)
	}
}
