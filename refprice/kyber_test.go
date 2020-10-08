package refprice

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestKyberGetRefPrice(t *testing.T) {
	fetcher := NewKyberFetcher()
	price, err := fetcher.GetRefPrice("KNC", "USDT")
	assert.NoError(t, err)
	assert.NotNil(t, price)
	log.Println(price)

	price, err = fetcher.GetRefPrice("KNCOIN", "USDT")
	assert.Error(t, err)
	assert.Nil(t, price)

	price, err = fetcher.GetRefPrice("KNC", "USD")
	assert.NoError(t, err)
	assert.NotNil(t, price)
	log.Println(price)

	price, err = fetcher.GetRefPrice("USD", "KNC")
	assert.Error(t, err)
	assert.Nil(t, price)
	log.Println(price)
}
