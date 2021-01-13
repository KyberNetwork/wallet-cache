package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/KyberNetwork/cache/ethereum"
	"github.com/KyberNetwork/cache/fetcher"
	"github.com/KyberNetwork/cache/http"
	"github.com/KyberNetwork/cache/node"
	persister "github.com/KyberNetwork/cache/persister"

	cli "gopkg.in/urfave/cli.v1"
)

type fetcherFunc func(memPersister persister.MemoryPersister, fetcher *fetcher.Fetcher)

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	//set log for server
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	app := cli.NewApp()
	app.Name = "Kyber Swap Cache"
	app.Usage = "Cache"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{}

	app.Commands = []cli.Command{}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Action = cmdMain

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func cmdMain(ctx *cli.Context) error {
	kyberENV := os.Getenv("KYBER_ENV")
	memPersister, err := persister.NewMemoryPersister("ram")
	if err != nil {
		log.Fatal(err)
	}
	diskPersister, err := persister.NewDiskPersister("leveldb")
	if err != nil {
		log.Fatal(err)
	}

	fetcherIns, err := fetcher.NewFetcher(kyberENV)
	if err != nil {
		log.Fatal(err)
	}
	nodeMiddleware, err := node.NewNodeMiddleware()
	if err != nil {
		log.Fatal(err)
	}

	err = fetcherIns.TryUpdateListToken()
	if err != nil {
		log.Println(err)
	}

	tickerUpdateToken := time.NewTicker(300 * time.Second)
	go func() {
		for {
			<-tickerUpdateToken.C
			fetcherIns.TryUpdateListToken()
		}
	}()

	//run fetch data
	runFetchData(memPersister, fetchKyberEnabled, fetcherIns, 10)
	runFetchData(memPersister, fetchMaxGasPrice, fetcherIns, 60)
	runFetchData(memPersister, fetchRateUSD, fetcherIns, 300)
	runFetchData(memPersister, fetchBlockNumber, fetcherIns, 10)

	go runFetchGasPrice(memPersister, diskPersister, fetcherIns)
	go fetchRate(memPersister, fetcherIns)

	server := http.NewHTTPServer(":3001", memPersister, diskPersister, fetcherIns, nodeMiddleware)
	server.Run(kyberENV)
	return nil
}

func runFetchData(memPersister persister.MemoryPersister, fn fetcherFunc, fertcherIns *fetcher.Fetcher, interval time.Duration) {
	ticker := time.NewTicker(interval * time.Second)
	go func() {
		for {
			fn(memPersister, fertcherIns)
			<-ticker.C
		}
	}()
}

// runFetchGasPrice save to memory each 30 seconds, save to disk each 1 hour
func runFetchGasPrice(memPersister persister.MemoryPersister, diskPersister persister.DiskPersister, fetcherIns *fetcher.Fetcher) {
	fastTicker := time.NewTicker(30 * time.Second)
	defer fastTicker.Stop()
	slowTicker := time.NewTicker(1 * time.Hour)
	defer slowTicker.Stop()

	var isNotFirstTime bool
	for {
		if !isNotFirstTime {
			gasPrice, err := fetcherIns.GetGasPrice()
			if err != nil {
				log.Println(err)
				memPersister.SetNewGasPrice(false)
				continue
			}
			memPersister.SaveGasPrice(gasPrice)

			if err := diskPersister.SaveGasPrice(*gasPrice); err != nil {
				log.Println(err)
				continue
			}
			isNotFirstTime = true
		}
		select {
		case <-fastTicker.C:
			gasPrice, err := fetcherIns.GetGasPrice()
			if err != nil {
				log.Println(err)
				memPersister.SetNewGasPrice(false)
				continue
			}

			memPersister.SaveGasPrice(gasPrice)
		case <-slowTicker.C:
			gasPrice := memPersister.GetGasPrice()
			if err := diskPersister.SaveGasPrice(*gasPrice); err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func fetchGasPrice(memPersister persister.MemoryPersister, fetcher *fetcher.Fetcher) {
	gasPrice, err := fetcher.GetGasPrice()
	if err != nil {
		log.Print(err)
		memPersister.SetNewGasPrice(false)
		return
	}
	memPersister.SaveGasPrice(gasPrice)
}

func fetchMaxGasPrice(memPersister persister.MemoryPersister, fetcher *fetcher.Fetcher) {
	gasPrice, err := fetcher.GetMaxGasPrice()
	if err != nil {
		log.Print(err)
		memPersister.SetNewMaxGasPrice(false)
		return
	}
	memPersister.SaveMaxGasPrice(gasPrice)
}

func fetchKyberEnabled(memPersister persister.MemoryPersister, fetcher *fetcher.Fetcher) {
	enabled, err := fetcher.CheckKyberEnable()
	if err != nil {
		log.Print(err)
		memPersister.SetNewKyberEnabled(false)
		return
	}
	memPersister.SaveKyberEnabled(enabled)
}

func fetchRateUSD(memPersister persister.MemoryPersister, fetcher *fetcher.Fetcher) {
	rateUSD, err := fetcher.GetRateUsdEther()
	if err != nil {
		log.Print(err)
		memPersister.SetNewRateUSD(false)
		return
	}

	if rateUSD == "" {
		memPersister.SetNewRateUSD(false)
		return
	}

	err = memPersister.SaveRateUSD(rateUSD)
	if err != nil {
		log.Print(err)
		memPersister.SetNewRateUSD(false)
		return
	}
}

func fetchBlockNumber(memPersister persister.MemoryPersister, fetcher *fetcher.Fetcher) {
	blockNum, err := fetcher.GetLatestBlock()
	if err != nil {
		log.Print(err)
		memPersister.SetNewLatestBlock(false)
		return
	}
	err = memPersister.SaveLatestBlock(blockNum)
	if err != nil {
		memPersister.SetNewLatestBlock(false)
		log.Print(err)
		return
	}
}

func fetchRate(memPersister persister.MemoryPersister, fetcher *fetcher.Fetcher) {
	ticker := time.NewTicker(15 * time.Second)
	for {
		var result []ethereum.Rate

		result, err := fetcher.FetchRate()
		if err != nil {
			log.Print(err)
			memPersister.SetIsNewRate(false)
			<-ticker.C
			continue
		}

		timeNow := time.Now().UTC().Unix()
		memPersister.SaveRate(result, timeNow)
		memPersister.SetIsNewRate(true)
		<-ticker.C
	}
}
