package bfetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KyberNetwork/cache/libs/limiter"
	"log"
	"time"

	"github.com/KyberNetwork/cache/ethereum"
	fCommon "github.com/KyberNetwork/cache/fetcher/fetcher-common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// api_key for tracker.kyber
const (
	etherscanRequestTimeout = 2 * time.Second
)

type Etherscan struct {
	url      string
	apiKey   string
	TypeName string
	limiter *limiter.RateLimiter
}

type ResultEvent struct {
	Result []ethereum.EventRaw `json:"result"`
}

func NewEtherScan(typeName string, url string, apiKey string) (*Etherscan, error) {
	l := limiter.NewRateLimiter(5, 1)

	etherscan := Etherscan{
		url: url,
		apiKey: apiKey,
		TypeName: typeName,
		limiter: l,
	}
	return &etherscan, nil
}

func (self *Etherscan) EthCall(to string, data string) (string, error) {
	url := self.url + "/api?module=proxy&action=eth_call&to=" +
		to + "&data=" + data + "&tag=latest&apikey=" + self.apiKey

	if err := self.limiter.WaitN(etherscanRequestTimeout, 1); err != nil {
		return "", errors.New("request timeout")
	}

	b, err := fCommon.HTTPCall(url)
	if err != nil {
		log.Print(err)
		return "", err
	}
	result := ethereum.ResultRpc{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return result.Result, nil

}

func (self *Etherscan) GetRate(to string, data string) (string, error) {
	return "", errors.New("not support this func")
}

func (self *Etherscan) GetLatestBlock() (string, error) {
	url := self.url + "/api?module=proxy&action=eth_blockNumber"
	if err := self.limiter.WaitN(etherscanRequestTimeout, 1); err != nil {
		return "", errors.New("request timeout")
	}
	b, err := fCommon.HTTPCall(url)
	if err != nil {
		log.Print(err)
		return "", err
	}
	blockNum := ethereum.ResultRpc{}
	err = json.Unmarshal(b, &blockNum)
	if err != nil {
		return "", err
	}
	num, err := hexutil.DecodeBig(blockNum.Result)
	if err != nil {
		return "", err
	}
	return num.String(), nil
}

func (self *Etherscan) GetTypeName() string {
	return self.TypeName
}

func (self *Etherscan) GetGasPrice() (*ethereum.GasPrice, error) {
	endpoint := fmt.Sprintf("%v/api?module=gastracker&action=gasoracle&apikey=%v", self.url, self.apiKey)
	if err := self.limiter.WaitN(etherscanRequestTimeout, 1); err != nil {
		return nil, errors.New("request timeout")
	}
	resultB, err := fCommon.HTTPCall(endpoint)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var response struct {
		Status string `json:"status"`
		Message string `json:"message"`
		Result struct {
			LastBlock string	`json:"LastBlock"`
			SafeGasPrice string 	`json:"SafeGasPrice"`
			ProposeGasPrice string 	`json:"ProposeGasPrice"`
			FastGasPrice string `json:"FastGasPrice"`
		}	`json:"result"`
	}
	if err := json.Unmarshal(resultB, &response); err != nil {
		log.Println(err)
		return nil, err
	}

	if response.Status != "1" {
		log.Println(err)
		return nil, errors.New(response.Message)
	}
	return &ethereum.GasPrice{
		Fast:     response.Result.FastGasPrice,
		Standard: response.Result.ProposeGasPrice,
		Low:      response.Result.SafeGasPrice,
		Default:  response.Result.ProposeGasPrice,
	}, nil
}
