package main

import (
	"github.com/KyberNetwork/cache/ethereum"
	"github.com/KyberNetwork/cache/fetcher"
	"github.com/KyberNetwork/cache/http"
	"github.com/KyberNetwork/cache/node"
	"github.com/KyberNetwork/cache/persister"
	"log"
	"os"
	"time"
)

const (
	snapshotGasPriceInterval = 1 * time.Hour
)

type CacheCore struct {
	memPersister  persister.MemoryPersister
	diskPersister persister.DiskPersister
	fetcherIns    *fetcher.Fetcher
	httpServer    *http.HTTPServer
}

func NewCacheCore() (*CacheCore, error) {
	memPersister, err := persister.NewMemoryPersister("ram")
	if err != nil {
		return nil, err
	}
	diskPersister, err := persister.NewDiskPersister("leveldb")
	if err != nil {
		return nil, err
	}

	fetcherIns, err := fetcher.NewFetcher(os.Getenv("KYBER_ENV"))
	if err != nil {
		return nil, err
	}
	nodeMiddleware, err := node.NewNodeMiddleware()
	if err != nil {
		return nil, err
	}
	err = fetcherIns.TryUpdateListToken()
	if err != nil {
		log.Println(err)
	}

	server := http.NewHTTPServer(":3001", memPersister, diskPersister, fetcherIns, nodeMiddleware)
	return &CacheCore{
		memPersister:  memPersister,
		diskPersister: diskPersister,
		fetcherIns:    fetcherIns,
		httpServer:    server,
	}, nil
}

func (c *CacheCore) Run() error {
	go c.FetchListToken()
	go c.runFetchData(fetchKyberEnabled, 10*time.Second)
	go c.runFetchData(fetchMaxGasPrice, 60*time.Second)
	go c.runFetchData(fetchRateUSD, 300*time.Second)
	go c.runFetchData(fetchBlockNumber, 10*time.Second)
	go c.runFetchData(fetchRate, 15*time.Second)
	go c.runFetchGasPrice()

	c.httpServer.Run(os.Getenv("KYBER_ENV"))
	return nil
}

func (c *CacheCore) FetchListToken() {
	ticker := time.NewTicker(300 * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		c.fetcherIns.TryUpdateListToken()
	}
}

func (c *CacheCore) runFetchGasPrice() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	var (
		snapshotGasPrices = make(map[int64]ethereum.GasPrice, 0)
		err               error
	)
	snapshotGasPrices, err = c.diskPersister.GetWeeklyGasPrice()
	if err != nil {
		log.Fatal(err)
	}
	// we take the first snapshot from disk and save to mem
	c.memPersister.SaveWeeklyGasPrice(snapshotGasPrices)

	var isFirstTime = true
	for {
		gasPrice, err := c.fetcherIns.GetGasPrice()
		if err != nil {
			log.Println(err)
			c.memPersister.SetNewGasPrice(false)
			continue
		}

		c.memPersister.SaveGasPrice(gasPrice)
		if isFirstTime {
			go c.snapshotGasPrice()
			isFirstTime = false
		}
		<-ticker.C
	}
}

func (c *CacheCore) snapshotGasPrice() {
	ticker := time.NewTicker(snapshotGasPriceInterval)
	defer ticker.Stop()
	for {
		c.memPersister.SnapshotWeeklyGasPrice()
		snapshotGasPrices := c.memPersister.GetWeeklyGasPrice()
		if err := c.diskPersister.SaveWeeklyGasPrice(snapshotGasPrices); err != nil {
			log.Println(err)
			<-ticker.C
			continue
		}

		<-ticker.C
	}
}

func (c *CacheCore) runFetchData(fn fetcherFunc, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		fn(c.memPersister, c.fetcherIns)
		<-ticker.C
	}
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

func fetchMaxGasPrice(memPersister persister.MemoryPersister, fetcher *fetcher.Fetcher) {
	gasPrice, err := fetcher.GetMaxGasPrice()
	if err != nil {
		log.Print(err)
		memPersister.SetNewMaxGasPrice(false)
		return
	}
	memPersister.SaveMaxGasPrice(gasPrice)
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
	var result []ethereum.Rate

	result, err := fetcher.FetchRate()
	if err != nil {
		log.Print(err)
		memPersister.SetIsNewRate(false)
		return
	}

	timeNow := time.Now().UTC().Unix()
	memPersister.SaveRate(result, timeNow)
	memPersister.SetIsNewRate(true)
}
