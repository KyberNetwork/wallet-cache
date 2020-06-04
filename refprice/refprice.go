package refprice

import (
	"errors"
	"math/big"
	"sync"
	"time"
)

type CachePrice struct {
	Base      string
	Quote     string
	Price     *big.Float
	Timestamp int64
}

type RefPrice struct {
	storage *ContractStorage
	cache   map[string]CachePrice
	fetcher *RefFetcher
	mu      sync.RWMutex
}

func NewRefPrice() *RefPrice {
	return &RefPrice{
		storage: NewContractStorage(),
		cache:   make(map[string]CachePrice),
		fetcher: NewRefFetcher(),
		mu:      sync.RWMutex{},
	}
}

type RefPriceApi struct {
	Price     string  `json:"price"`
	ThresHold float64 `json:"threshold"`
}

func (r *RefPrice) GetRefPrice(base string, quote string) (RefPriceApi, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// fetch from chain link
	contract := r.storage.GetContract(base, quote)
	if contract.Address == "" {
		return RefPriceApi{}, errors.New("Cannot get chainlink contract")
	}

	if c, ok := r.cache[getKey(base, quote)]; ok {
		if time.Now().Unix()-c.Timestamp < 120 { //cache 2 minutes
			return RefPriceApi{
				Price:     c.Price.String(),
				ThresHold: contract.ThresHold,
			}, nil
		}
	}

	// get refprice from blockchain
	price, err := r.fetcher.GetRefPrice(contract.Address)
	if err != nil {
		return RefPriceApi{}, err
	}

	tokenPrice := getTokenPrice(price, contract.Multiply)

	r.cache[getKey(base, quote)] = CachePrice{
		Base: base, Quote: quote, Price: tokenPrice, Timestamp: time.Now().Unix(),
	}

	return RefPriceApi{
		Price:     tokenPrice.String(),
		ThresHold: contract.ThresHold,
	}, nil
}

func getTokenPrice(price *big.Int, multiplier *big.Int) *big.Float {
	priceF := new(big.Float).SetInt(price)
	mulF := new(big.Float).SetInt(multiplier)

	result := new(big.Float).Quo(priceF, mulF)
	return result
}

func getKey(base string, quote string) string {
	return base + quote
}
