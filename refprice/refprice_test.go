package refprice

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestGetRefPrice(t *testing.T) {
	os.Setenv("NODE_ENDPOINT", "https://ethereum.knstats.com/v1/mainnet/geth")
	ref := NewRefPrice()
	price, err := ref.GetRefPrice("ETH", "USDT")
	assert.NoError(t, err)
	assert.NotEqual(t, 0, price)

	log.Println(price)
}