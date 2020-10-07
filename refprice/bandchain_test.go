package refprice

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestBandchainFetcher(t *testing.T) {
	fetcher := BandchainFetcher{}

	price, err := fetcher.GetRefPrice("ETH", "USDT")
	assert.NoError(t, err)
	assert.NotNil(t, price)
	log.Println(price)
}