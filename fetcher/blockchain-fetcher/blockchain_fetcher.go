package bfetcher

import (
	"context"
	"log"
	"math/big"
	"time"

	// "strconv"
	"github.com/KyberNetwork/cache/ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

type BlockchainFetcher struct {
	client   *rpc.Client
	url      string
	TypeName string

	timeout time.Duration
}

func NewBlockchainFetcher(typeName string, endpoint string, apiKey string) (*BlockchainFetcher, error) {
	client, err := rpc.DialHTTP(endpoint)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	timeout := 5 * time.Second
	blockchain := BlockchainFetcher{
		client:   client,
		url:      endpoint,
		TypeName: typeName,
		timeout:  timeout,
	}
	return &blockchain, nil
}

func (self *BlockchainFetcher) EthCall(to string, data string) (string, error) {
	params := make(map[string]string)
	params["data"] = "0x" + data
	params["to"] = to

	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	var result string
	err := self.client.CallContext(ctx, &result, "eth_call", params, "latest")
	if err != nil {
		log.Print(err)
		return "", err
	}

	return result, nil
}

func (self *BlockchainFetcher) GetRate(to string, data string) (string, error) {
	params := make(map[string]string)
	params["data"] = "0x" + data
	params["to"] = to

	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	var result string
	err := self.client.CallContext(ctx, &result, "eth_call", params, "latest")
	if err != nil {
		log.Print(err)
		return "", err
	}

	return result, nil
}

func (self *BlockchainFetcher) GetGasPrice() (*ethereum.GasPrice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	var gasPrice *hexutil.Big
	err := self.client.CallContext(ctx, &gasPrice, "eth_gasPrice")

	if err != nil {
		log.Print(err)
		return nil, err
	}

	gasPriceInt := new(big.Float).SetInt(gasPrice.ToInt())
	gasPriceWeight := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)
	gasPriceWeightF := new(big.Float).SetInt(gasPriceWeight)

	gasPriceGwei := new(big.Float).Quo(gasPriceInt, gasPriceWeightF)

	if gasPriceGwei.Cmp(big.NewFloat(1)) < 1 {
		return &ethereum.GasPrice{
			Fast:     "1",
			Standard: "1",
			Low:      "1",
			Default:  "1",
		}, nil
	}

	gasPriceFast := new(big.Float).Mul(gasPriceGwei, big.NewFloat(1.2))

	return &ethereum.GasPrice{
		Fast:     gasPriceFast.String(),
		Standard: gasPriceGwei.String(),
		Low:      gasPriceGwei.String(),
		Default:  gasPriceGwei.String(),
	}, nil
}

func (self *BlockchainFetcher) GetLatestBlock() (string, error) {
	var blockNum *hexutil.Big
	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	err := self.client.CallContext(ctx, &blockNum, "eth_blockNumber", "latest")
	if err != nil {
		return "", err
	}
	return blockNum.ToInt().String(), nil
}

type TopicParam struct {
	FromBlock string   `json:"fromBlock"`
	ToBlock   string   `json:"toBlock"`
	Address   string   `json:"address"`
	Topics    []string `json:"topics"`
}

func (self *BlockchainFetcher) GetTypeName() string {
	return self.TypeName
}
