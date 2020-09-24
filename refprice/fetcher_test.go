package refprice

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRefPrice(t *testing.T) {
	os.Setenv("NODE_ENDPOINT", "https://eth-mainnet.alchemyapi.io/v2/SzkoydhV1A0B-0s73tAANYoU71jriFkE")
	f := NewRefFetcher()
	price, err := f.GetRefPrice("0x656c0544eF4C98A6a98491833A89204Abb045d6b")	// use proxy address
	assert.Equal(t, nil, err)
	log.Print(price.String())
}
