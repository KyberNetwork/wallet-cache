package refprice

import (
	"context"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/KyberNetwork/cache/libs/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RefFetcher struct {
	client *ethclient.Client
}

func NewRefFetcher() *RefFetcher {
	clientIns, err := ethclient.Dial(os.Getenv("NODE_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}

	return &RefFetcher{
		client: clientIns,
	}
}

func (f *RefFetcher) GetRefPrice(address string) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Context: ctx,
	}

	aggrator, err := contracts.NewAggregator(ethereum.HexToAddress(address), f.client)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return aggrator.LatestAnswer(opts)
}
