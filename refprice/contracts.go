package refprice

import (
	"encoding/json"
	"log"
	"math/big"
	"sync"
	"time"
)

type Contract struct {
	Base     string
	Quote    string
	Multiply *big.Int
	Address  string
}

type ContractStorage struct {
	contracts []Contract

	mu sync.RWMutex
}

func NewContractStorage() *ContractStorage {
	s := &ContractStorage{
		contracts: make([]Contract, 0),
		mu:        sync.RWMutex{},
	}
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		for {
			contracts, err := fetchContractList()
			if err != nil {
				log.Print(err)
				<-ticker.C
				continue
			}
			s.saveContracts(contracts)
			<-ticker.C
		}
	}()
	return s
}

func (s *ContractStorage) saveContracts(contracts []Contract) {
	s.mu.Lock()
	defer s.mu.Unlock()
	newContracts := make([]Contract, 0)
	for _, c := range contracts {
		newContracts = append(newContracts, c)
	}
	s.contracts = newContracts
}

func (s *ContractStorage) GetContract(base string, quote string) Contract {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.contracts {
		if c.Base == base && c.Quote == quote {
			return c
		}
	}
	return Contract{}
}

type ContractRes struct {
	Address  string    `json:"contractAddress"`
	Multiply string    `json:"multiply"`
	Pair     [2]string `json:"pair"`
}

func fetchContractList() ([]Contract, error) {
	b, err := HTTPCall("https://weiwatchers.com/feeds.json")
	if err != nil {
		log.Print(err)
		return nil, err
	}
	var result []ContractRes
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	contracts := make([]Contract, 0)
	for _, c := range result {

		if n, ok := new(big.Int).SetString(c.Multiply, 10); ok {
			contracts = append(contracts, Contract{
				Base:     c.Pair[0],
				Quote:    c.Pair[1],
				Multiply: n,
				Address:  c.Address,
			})
		}
	}
	return contracts, nil
}
