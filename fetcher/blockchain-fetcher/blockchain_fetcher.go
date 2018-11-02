package bfetcher

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"

	// "strconv"

	"github.com/KyberNetwork/server-go/ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

type BlockchainFetcher struct {
	client   *rpc.Client
	url      string
	TypeName string
}

func NewBlockchainFetcher(typeName string, endpoint string, apiKey string) (*BlockchainFetcher, error) {
	client, err := rpc.DialHTTP(endpoint)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	blockchain := BlockchainFetcher{
		client, endpoint, typeName,
	}
	return &blockchain, nil
}

func (self *BlockchainFetcher) EthCall(to string, data string) (string, error) {
	params := make(map[string]string)
	params["data"] = "0x" + data
	params["to"] = to

	var result string
	err := self.client.Call(&result, "eth_call", params, "latest")
	if err != nil {
		log.Print(err)
		return "", err
	}

	return result, nil

}

func (self *BlockchainFetcher) GetLatestBlock() (string, error) {
	var blockNum *hexutil.Big
	err := self.client.Call(&blockNum, "eth_blockNumber", "latest")
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

// func (self *BlockchainFetcher) GetEvents(fromBlock string, toBlock string, network string, tradeTopic string) (*[]ethereum.EventRaw, error) {
// 	fromBlockInt, err := strconv.ParseUint(fromBlock, 10, 64)
// 	if err != nil {
// 		log.Print(err)
// 		return nil, err
// 	}

// 	toBlockInt, err := strconv.ParseUint(toBlock, 10, 64)
// 	if err != nil {
// 		log.Print(err)
// 		return nil, err
// 	}

// 	fromBlockHex := hexutil.EncodeUint64(fromBlockInt)
// 	toBlockHex := hexutil.EncodeUint64(toBlockInt)

// 	param := TopicParam{
// 		fromBlockHex, toBlockHex, network, []string{tradeTopic},
// 	}

// 	var result []ethereum.EventRaw
// 	err = self.client.Call(&result, "eth_getLogs", param)
// 	if err != nil {
// 		log.Print(err)
// 		return nil, err
// 	}

// 	return &result, nil
// }

func (self *BlockchainFetcher) GetRateUsd(tickers []string) ([]io.ReadCloser, error) {
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

func (self *BlockchainFetcher) GetTypeName() string {
	return self.TypeName
}

func (self *BlockchainFetcher) GetGasPrice() (*ethereum.GasPrice, error) {
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

func (self *BlockchainFetcher) GetRateUsdEther() (string, error) {
	response, err := http.Get("https://api.coinmarketcap.com/v1/ticker/ethereum")
	if err != nil {
		log.Print(err)
		return "", err
	}
	defer (response.Body).Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return "", err
	}
	rateItem := make([]RateUSD, 0)
	err = json.Unmarshal(b, &rateItem)
	if err != nil {
		log.Print(err)
		return "", err
	}
	return rateItem[0].PriceUsd, nil
}

func (self *BlockchainFetcher) GetGeneralInfo(usdId string) (*ethereum.TokenGeneralInfo, error) {
	err := errors.New("Blockchain is not support this api")
	//log.Print(err)
	return nil, err
}

func (self *BlockchainFetcher) GetTrackerData(trackerEndpoint string) (map[string]*ethereum.Rates, error) {
	trackerData := map[string]*ethereum.Rates{}
	err := errors.New("Blockchain is not support this api")
	//log.Print(err)
	return trackerData, err
}
