package refprice

import (
	"context"
	"errors"
	"github.com/KyberNetwork/cache/libs/contracts"
	etherCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
	"time"
)

type ChainlinkFetcher struct {
	storage *ChainlinkContractStorage
	client *ethclient.Client
}

func NewChainlinkFetcher() *ChainlinkFetcher {
	clientIns, err := ethclient.Dial(os.Getenv("NODE_ENDPOINT"))
	if err != nil {
		log.Fatal("missing `NODE_ENDPOINT` env")
	}
	return &ChainlinkFetcher{
		storage: NewChainlinkContractStorage(),
		client: clientIns,
	}
}

func (f *ChainlinkFetcher) GetRefPrice(base, quote string) (*big.Float, error) {
	contract := f.storage.GetContract(base, quote)
	if contract.Address == "" {
		return nil, errors.New("cannot get chainlink contract")
	}

	price, err := f.fetchPrice(contract.Address)
	if err != nil {
		return nil, err
	}

	tokenPrice := GetTokenPrice(price, contract.Multiply)
	return tokenPrice, nil
}

// fetchPrice fetch price from contract
func (f *ChainlinkFetcher) fetchPrice(address string) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	opts := &bind.CallOpts{
		Context: ctx,
	}

	aggrator, err := contracts.NewAggregator(etherCommon.HexToAddress(address), f.client)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return aggrator.LatestAnswer(opts)
}

func GetTokenPrice(price *big.Int, multiplier *big.Int) *big.Float {
	priceF := new(big.Float).SetInt(price)
	mulF := new(big.Float).SetInt(multiplier)

	result := new(big.Float).Quo(priceF, mulF)
	return result
}
