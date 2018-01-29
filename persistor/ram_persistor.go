package persistor

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"sync"

	"github.com/KyberNetwork/server-go/ethereum"
)

type RamPersistor struct {
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

func NewRamPersistor() (*RamPersistor, error) {
	var mu sync.RWMutex
	rates := make([]ethereum.Rate, 0)
	latestBlock := "0"
	rateUSD := make([]RateUSD, 0)
	events := make([]ethereum.EventHistory, 0)

	isNewRate := true
	isNewLatestBlock := true
	isNewRateUsd := true
	isNewEvent := true

	persistor := &RamPersistor{
		mu, &rates, latestBlock, rateUSD, events, isNewRate, isNewLatestBlock, isNewRateUsd, isNewEvent,
	}
	return persistor, nil
}

func (self *RamPersistor) GetRate() *[]ethereum.Rate {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.rates
}

func (self *RamPersistor) GetEvent() []ethereum.EventHistory {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.events
}

func (self *RamPersistor) GetLatestBlock() string {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.latestBlock
}

func (self *RamPersistor) GetRateUSD() []RateUSD {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.rateUSD
}

func (self *RamPersistor) GetIsNewRate() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewRate
}

func (self *RamPersistor) GetIsNewLatestBlock() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewLatestBlock
}

func (self *RamPersistor) GetIsNewRateUSD() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewRateUsd
}

func (self *RamPersistor) GetIsNewEvent() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isNewEvent
}

func (self *RamPersistor) SaveRate(rates *[]ethereum.Rate) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.rates = rates
	return nil
}

func (self *RamPersistor) SaveEvent(events *[]ethereum.EventHistory) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.events = *events
	return nil
}

func (self *RamPersistor) SaveLatestBlock(blockNumber string) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.latestBlock = blockNumber
	self.isNewLatestBlock = true
	return nil
}

func (self *RamPersistor) SaveRateUSD(body []*io.ReadCloser) error {
	rates := make([]RateUSD, 0)
	for _, item := range body {
		rateItem := make([]RateUSD, 0)
		b, err := ioutil.ReadAll(*item)
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

func (self *RamPersistor) SetNewRate(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewRate = isNew
}

func (self *RamPersistor) SetNewLatestBlock(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewLatestBlock = isNew
}

func (self *RamPersistor) SetNewRateUSD(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewRateUsd = isNew
}

func (self *RamPersistor) SetNewEvents(isNew bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewEvent = isNew
}
