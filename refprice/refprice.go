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

func (r *RefPrice) GetRefPrice(base string, quote string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.cache[getKey(base, quote)]; ok {
		if time.Now().Unix()-c.Timestamp < 120 { //cache 2 minutes
			return c.Price.String(), nil
		}
	}

	// fetch from chain link
	contract := r.storage.GetContract(base, quote)
	if contract.Address == "" {
		return "", errors.New("Cannot get chainlink contract")
	}

	// get refprice from blockchain
	price, err := r.fetcher.GetRefPrice(contract.Address)
	if err != nil {
		return "", err
	}

	tokenPrice := getTokenPrice(price, contract.Multiply)

	r.cache[getKey(base, quote)] = CachePrice{
		Base: base, Quote: quote, Price: tokenPrice, Timestamp: time.Now().Unix(),
	}

	return tokenPrice.String(), nil
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
