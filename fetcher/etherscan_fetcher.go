package fetcher

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"

	"github.com/KyberNetwork/server-go/ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	response, err := http.Get(url)

	if err != nil {
		log.Print(err)
		return "", err
	}
	if response.StatusCode != 200 {
		return "", errors.New("Status code is 200")
	}

	defer (response.Body).Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return "", err
	}
	result := ResultRpc{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return result.Result, nil

}

func (self *Etherscan) GetLatestBlock() (string, error) {
	response, err := http.Get(self.url + "/api?module=proxy&action=eth_blockNumber")
	if err != nil {
		log.Print(err)
		return "", err
	}
	if response.StatusCode != 200 {
		return "", errors.New("Status code is 200")
	}
	defer (response.Body).Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	blockNum := ResultRpc{}
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

func (self *Etherscan) GetEvents(fromBlock string, toBlock string, network string, tradeTopic string) (*[]ethereum.EventRaw, error) {
	url := self.url + "/api?module=logs&action=getLogs&fromBlock=" +
		fromBlock + "&toBlock=" + toBlock + "&address=" + network + "&topic0=" +
		tradeTopic + "&apikey=" + self.apiKey
	response, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("Status code is 200")
	}

	defer (response.Body).Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	result := ResultEvent{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &result.Result, nil
}

func (self *Etherscan) GetRateUsd(tickers []string) ([]io.ReadCloser, error) {
	outPut := make([]io.ReadCloser, 0)
	for _, ticker := range tickers {
		response, err := http.Get("https://api.coinmarketcap.com/v1/ticker/" + ticker)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		outPut = append(outPut, response.Body)
	}

	return outPut, nil
}

func (self *Etherscan) GetTypeName() string {
	return self.TypeName
}

type GasStation struct {
	Fast     float64 `json:"fast"`
	Standard float64 `json:"average"`
	Low      float64 `json:"safeLow"`
}

func (self *Etherscan) GetGasPrice() (*ethereum.GasPrice, error) {
	response, err := http.Get("https://ethgasstation.info/json/ethgasAPI.json")
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("Status code is 200")
	}
	defer (response.Body).Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var gasPrice GasStation
	err = json.Unmarshal(b, &gasPrice)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	fast := big.NewFloat(gasPrice.Fast / 10)
	standard := big.NewFloat((gasPrice.Fast + gasPrice.Standard) / 20)
	low := big.NewFloat(gasPrice.Low / 10)
	defaultGas := standard

	return &ethereum.GasPrice{
		fast.String(), standard.String(), low.String(), defaultGas.String(),
	}, nil
}
