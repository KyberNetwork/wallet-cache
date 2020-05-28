package refprice

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetContractChainLink(t *testing.T) {
	contracts, err := fetchContractList()
	assert.Equal(t, nil, err)
	log.Print(contracts)
}
