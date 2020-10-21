package refprice

import (
	"math/big"
	"sync"
	"time"
)

type CachePrice struct {
	Base      string
	Quote     string
	Price     *big.Float
	Timestamp int64
	SourceData []string
}

type RefPrice struct {
	kyberFetcher *KyberFetcher
	chainlinkFetcher *ChainlinkFetcher
	bandchainFetcher *BandchainFetcher
	cache   map[string]CachePrice
	mu      sync.RWMutex
}

func NewRefPrice() *RefPrice {
	return &RefPrice{
		kyberFetcher: NewKyberFetcher(),
		chainlinkFetcher: NewChainlinkFetcher(),
		bandchainFetcher: NewBandchainFetcher(),
		cache:   make(map[string]CachePrice),
		mu:      sync.RWMutex{},
	}
}

// GetRefPrice get reference price from multiple sources data (ex: Kyber, Chainlink, Bandchain)
func (r *RefPrice) GetRefPrice(base string, quote string) (price string, sourceData []string, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.cache[getKey(base, quote)]; ok {
		if time.Now().Unix()-c.Timestamp < 120 { //cache 2 minutes
			return c.Price.String(), c.SourceData, nil
		}
	}
	sourceData = make([]string, 0)

	chainlinkPrice, err := r.chainlinkFetcher.GetRefPrice(base, quote)
	if err == nil {
		sourceData = append(sourceData, "ChainLink")
	}

	bandchainPrice, err := r.bandchainFetcher.GetRefPrice(base, quote)
	if err == nil {
		sourceData = append(sourceData, "BandChain")
	}

	kyberPrice, err := r.kyberFetcher.GetRefPrice(base, quote)
	if err == nil {
		sourceData = append(sourceData, "KyberNetwork")
	}

	result := getAvgPrice([]*big.Float{chainlinkPrice, bandchainPrice, kyberPrice})
	r.cache[getKey(base, quote)] = CachePrice{
		Base: base, Quote: quote, Price: result, Timestamp: time.Now().Unix(), SourceData: sourceData,
	}

	return result.String(), sourceData, nil
}

func getAvgPrice(prices []*big.Float) *big.Float {
	var (
		avgPrice = big.NewFloat(0)
		zero = big.NewFloat(0)
		counter float64
	)
	if len(prices) == 0 {
		return avgPrice
	}
	for _, p := range prices {
		if p != nil && p.Cmp(zero) != 0 {
			avgPrice = avgPrice.Add(avgPrice, p)
			counter += 1
		}
	}

	if counter == 0 {
		return avgPrice
	}
	return new(big.Float).Quo(avgPrice, big.NewFloat(counter))
}

func getKey(base string, quote string) string {
	return base + quote
}
