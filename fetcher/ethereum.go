package fetcher

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/KyberNetwork/server-go/ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

type ResultEvent struct {
	Result []ethereum.EventRaw `json:"result"`
}

type ResultRpc struct {
	Result string `json:"result"`
}

type RateWapper struct {
	ExpectedPrice []*big.Int
	SlippagePrice []*big.Int
}

type Ethereum struct {
	nodeEndpoint string
	network      string
	networkAbi   abi.ABI
	tradeTopic   string
	wrapper      string
	wrapperAbi   abi.ABI
	client       *rpc.Client
}

func NewEthereum(nodeEndpoint string, network string, networkAbiStr string, tradeTopic string, wrapper string, wrapperAbiStr string) (*Ethereum, error) {
	client, err := rpc.DialHTTP(nodeEndpoint)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	networkAbi, err := abi.JSON(strings.NewReader(networkAbiStr))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	wrapperAbi, err := abi.JSON(strings.NewReader(wrapperAbiStr))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	ethereum := &Ethereum{
		nodeEndpoint, network, networkAbi, tradeTopic, wrapper, wrapperAbi, client,
	}

	return ethereum, nil
}

func (self *Ethereum) ExactBlockNumber(body io.ReadCloser) (string, error) {
	defer (body).Close()
	b, err := ioutil.ReadAll(body)
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

func (self *Ethereum) GetLatestBlock() (string, error) {
	var blockNum *hexutil.Big
	err := self.client.Call(&blockNum, "eth_blockNumber", "latest")
	if err != nil {
		return "", err
	}
	return blockNum.ToInt().String(), nil
}

func (self *Ethereum) EncodeRateData(source []string, dest []string, quantity []int64) (string, error) {
	sourceList := make([]common.Address, 0)
	for _, sourceItem := range source {
		sourceList = append(sourceList, common.HexToAddress(sourceItem))
	}
	destList := make([]common.Address, 0)
	for _, destItem := range dest {
		destList = append(destList, common.HexToAddress(destItem))
	}

	quantityList := make([]*big.Int, 0)
	for _, quanItem := range quantity {
		quantityList = append(quantityList, big.NewInt(quanItem))
	}

	encodedData, err := self.wrapperAbi.Pack("getExpectedRates", common.HexToAddress(self.network), sourceList, destList, quantityList)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return common.Bytes2Hex(encodedData), nil
}

type RateWrapper struct {
	ExpectedRate []*big.Int `json:"expectedRate"`
	SlippageRate []*big.Int `json:"slippageRate"`
}

func (self *Ethereum) ExactRateDataFromEtherscan(body io.ReadCloser, sourceArr []string, destAddr []string) (*[]ethereum.Rate, error) {
	defer (body).Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	result := ResultRpc{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	rates, err := self.ReadRate(result.Result, sourceArr, destAddr)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return rates, nil
}

// type RateParam struct {
// 	Network common.Address   `json:"network"`
// 	Sources []common.Address `json:"sources"`
// 	Dests   []common.Address `json:"dests"`
// 	Quanity []*big.Int  `json:"qty"`
// }

func (self *Ethereum) GetRateFromNode(sourceSymbolArr []string, destSymbolArr []string, dataAbi string) (*[]ethereum.Rate, error) {

	params := make(map[string]string)
	params["data"] = "0x" + dataAbi
	params["to"] = self.wrapper

	//fmt.Print(params)
	var result string
	err := self.client.Call(&result, "eth_call", params, "latest")
	if err != nil {
		log.Print(err)
		return nil, err
	}

	rates, err := self.ReadRate(result, sourceSymbolArr, destSymbolArr)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return rates, nil
}

func (self *Ethereum) ReadRate(result string, sourceArr []string, destArr []string) (*[]ethereum.Rate, error) {
	rateByte, err := hexutil.Decode(result)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	var rateWapper RateWrapper
	err = self.wrapperAbi.Unpack(&rateWapper, "getExpectedRates", rateByte)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var lenArr = len(sourceArr)
	if (len(rateWapper.ExpectedRate) != lenArr) || (len(rateWapper.SlippageRate) != lenArr) {
		errorLength := errors.New("Length of expected for slippage rate is not enough")
		log.Print(errorLength)
		return nil, errorLength
	}

	rateReturn := make([]ethereum.Rate, 0)
	for i := 0; i < lenArr; i++ {
		source := sourceArr[i]
		dest := destArr[i]
		rate := rateWapper.ExpectedRate[i]
		minRate := rateWapper.SlippageRate[i]
		rateReturn = append(rateReturn, ethereum.Rate{
			source, dest, rate.String(), minRate.String(),
		})
	}
	return &rateReturn, nil
}

type LogData struct {
	Source           common.Address `json:"source"`
	Dest             common.Address `json:"dest"`
	ActualSrcAmount  *big.Int       `json:"actualSrcAmount"`
	ActualDestAmount *big.Int       `json:"actualDestAmount"`
}

func (self *Ethereum) ExactEventFromEtherscan(body io.ReadCloser) (*[]ethereum.EventHistory, error) {
	defer (body).Close()
	b, err := ioutil.ReadAll(body)
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

	listEvent := result.Result
	events, err := self.ReadEvents(&listEvent, "etherscan", nil, 0)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return events, nil
}

type TopicParam struct {
	FromBlock string   `json:"fromBlock"`
	ToBlock   string   `json:"toBlock"`
	Address   string   `json:"address"`
	Topics    []string `json:"topics"`
}

func (self *Ethereum) GetEventsFromNode(fromBlock string, toBlock string, averageBlockTime int64) (*[]ethereum.EventHistory, error) {
	fromBlockInt, err := strconv.ParseUint(fromBlock, 10, 64)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	toBlockInt, err := strconv.ParseUint(toBlock, 10, 64)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	fromBlockHex := hexutil.EncodeUint64(fromBlockInt)
	toBlockHex := hexutil.EncodeUint64(toBlockInt)

	param := TopicParam{
		fromBlockHex, toBlockHex, self.network, []string{self.tradeTopic},
	}
	var result []ethereum.EventRaw
	err = self.client.Call(&result, "eth_getLogs", param)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	//get latestBlock to calculate timestamp
	latestBlock, err := self.GetLatestBlock()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	events, err := self.ReadEvents(&result, "node", &latestBlock, averageBlockTime)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return events, nil
}

func (self *Ethereum) ReadEvents(listEventAddr *[]ethereum.EventRaw, typeFetch string, latestBlock *string, averageBlockTime int64) (*[]ethereum.EventHistory, error) {
	listEvent := *listEventAddr
	endIndex := len(listEvent) - 1
	var beginIndex = 0
	if endIndex > 4 {
		beginIndex = endIndex - 4
	}

	events := make([]ethereum.EventHistory, 0)
	for i := endIndex; i >= beginIndex; i-- {
		txHash := listEvent[i].Txhash

		blockNumber, err := hexutil.DecodeBig(listEvent[i].BlockNumber)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		var timestamp string
		if typeFetch == "etherscan" {
			timestampHex, err := hexutil.DecodeBig(listEvent[i].Timestamp)
			if err != nil {
				log.Print(err)
				return nil, err
			}
			timestamp = timestampHex.String()
		} else {
			timestamp, err = self.Gettimestamp(blockNumber.String(), *latestBlock, averageBlockTime)
			if err != nil {
				log.Print(err)
				return nil, err
			}
		}

		var logData LogData
		data, err := hexutil.Decode(listEvent[i].Data)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		//fmt.Print(listEvent[i].Data)
		err = self.networkAbi.Unpack(&logData, "ExecuteTrade", data)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		actualDestAmount := logData.ActualDestAmount.String()
		actualSrcAmount := logData.ActualSrcAmount.String()
		dest := logData.Dest.String()
		source := logData.Source.String()

		events = append(events, ethereum.EventHistory{
			actualDestAmount, actualSrcAmount, dest, source, blockNumber.String(), txHash, timestamp,
		})
	}

	return &events, nil
}

func (self *Ethereum) Gettimestamp(block string, latestBlock string, averageBlockTime int64) (string, error) {
	fromBlock, err := strconv.ParseInt(block, 10, 64)
	if err != nil {
		log.Print(err)
		return "", err
	}
	toBlock, err := strconv.ParseInt(block, 10, 64)
	if err != nil {
		log.Print(err)
		return "", err
	}
	timeNow := time.Now().Unix()
	timeStamp := timeNow - averageBlockTime*(toBlock-fromBlock)

	timeStampBig := big.NewInt(timeStamp)

	return timeStampBig.String(), nil

}
