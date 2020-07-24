package bfetcher

import (
	"log"
	"testing"
)

func TestGasPriceBlockchain(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	b, _ := NewBlockchainFetcher("blockchain", "https://eth-mainnet.alchemyapi.io/v2/SzkoydhV1A0B-0s73tAANYoU71jriFkE", "")
	gasPrice, _ := b.GetGasPrice()
	log.Print(gasPrice)
}
