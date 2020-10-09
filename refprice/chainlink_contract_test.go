package refprice

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestFetContractChainLink(t *testing.T) {
	contracts, err := fetchContractList()
	assert.Equal(t, nil, err)
	log.Print(contracts)
}

