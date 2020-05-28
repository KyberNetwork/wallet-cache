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
	price, err := f.GetRefPrice("0xd0e785973390fF8E77a83961efDb4F271E6B8152")
	assert.Equal(t, nil, err)
	log.Print(price.String())
}
