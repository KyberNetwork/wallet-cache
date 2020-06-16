package bfetcher

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/KyberNetwork/cache/ethereum"
	fCommon "github.com/KyberNetwork/cache/fetcher/fetcher-common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// api_key for tracker.kyber
const (
	TIME_TO_DELETE = 18000
)

type Etherscan struct {
	url      string
	apiKey   string
	TypeName string
}

type ResultEvent struct {
	Result []ethereum.EventRaw `json:"result"`
}

func NewEtherScan(typeName string, url string, apiKey string) (*Etherscan, error) {
	etherscan := Etherscan{
		url, apiKey, typeName,
	}
	return &etherscan, nil
}

func (self *Etherscan) EthCall(to string, data string) (string, error) {
	url := self.url + "/api?module=proxy&action=eth_call&to=" +
		to + "&data=" + data + "&tag=latest&apikey=" + self.apiKey
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
