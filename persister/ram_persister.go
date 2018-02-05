package Persister

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"sync"

	"github.com/KyberNetwork/server-go/ethereum"
)

type RamPersister struct {
	mu          sync.RWMutex
	rates       *[]ethereum.Rate
	latestBlock string
	rateUSD     []RateUSD
	events      []ethereum.EventHistory

	isNewRate        bool
	isNewLatestBlock bool
	isNewRateUsd     bool
	isNewEvent       bool
}

func NewRamPersister() (*RamPersister, error) {
	var mu sync.RWMutex
	rates := make([]ethereum.Rate, 0)
	latestBlock := "0"
	rateUSD := make([]RateUSD, 0)
	events := make([]ethereum.EventHistory, 0)

	isNewRate := true
	isNewLatestBlock := true
	isNewRateUsd := true
	isNewEvent := true

	persister := &RamPersister{
		mu, &rates, latestBlock, rateUSD, events, isNewRate, isNewLatestBlock, isNewRateUsd, isNewEvent,
	}
	return persister, nil
}

func (self *RamPersister) GetRate() *[]ethereum.Rate {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.rates
}

func (self *RamPersister) GetEvent() []ethereum.EventHistory {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.events
}

func (self *RamPersister) GetLatestBlock() string {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.latestBlock
}

func (self *RamPersister) GetRateUSD() []RateUSD {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.rateUSD
}

func (self *RamPersister) GetIsNewRate() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewRate
}

func (self *RamPersister) GetIsNewLatestBlock() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewLatestBlock
}

func (self *RamPersister) GetIsNewRateUSD() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewRateUsd
}

func (self *RamPersister) GetIsNewEvent() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewEvent
}

func (self *RamPersister) SaveRate(rates *[]ethereum.Rate) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.rates = rates
	return nil
}

func (self *RamPersister) SaveEvent(events *[]ethereum.EventHistory) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.events = *events
	return nil
}

func (self *RamPersister) SaveLatestBlock(blockNumber string) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.latestBlock = blockNumber
	self.isNewLatestBlock = true
	return nil
}

func (self *RamPersister) SaveRateUSD(body []io.ReadCloser) error {
	rates := make([]RateUSD, 0)
	for _, item := range body {
		rateItem := make([]RateUSD, 0)
		defer (item).Close()
		b, err := ioutil.ReadAll(item)
		if err != nil {
			log.Print(err)
			return err
		}
		err = json.Unmarshal(b, &rateItem)
		if err != nil {
			log.Print(err)
			return err
		}
		rates = append(rates, rateItem[0])
	}
	self.mu.Lock()
	defer self.mu.Unlock()
	self.rateUSD = rates
	self.isNewRateUsd = true
	return nil
}

func (self *RamPersister) SetNewRate(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewRate = isNew
}

func (self *RamPersister) SetNewLatestBlock(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewLatestBlock = isNew
}

func (self *RamPersister) SetNewRateUSD(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewRateUsd = isNew
}

func (self *RamPersister) SetNewEvents(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewEvent = isNew
}
