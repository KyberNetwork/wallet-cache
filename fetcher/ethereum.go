package fetcher

import (
	"errors"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/KyberNetwork/server-go/ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type RateWrapper struct {
	ExpectedRate []*big.Int `json:"expectedRate"`
	SlippageRate []*big.Int `json:"slippageRate"`
}

type Ethereum struct {
	network          string
	networkAbi       abi.ABI
	tradeTopic       string
	wrapper          string
	wrapperAbi       abi.ABI
	averageBlockTime int64
}

func NewEthereum(network string, networkAbiStr string, tradeTopic string, wrapper string, wrapperAbiStr string, averageBlockTime int64) (*Ethereum, error) {

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
		network, networkAbi, tradeTopic, wrapper, wrapperAbi, averageBlockTime,
	}

	return ethereum, nil
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

func (self *Ethereum) EncodeKyberEnable() (string, error) {
	encodedData, err := self.networkAbi.Pack("enabled")
	if err != nil {
		log.Print(err)
		return "", err
	}
	return common.Bytes2Hex(encodedData), nil
}

func (self *Ethereum) ExtractEnabled(result string) (bool, error) {
	enabledByte, err := hexutil.Decode(result)
	if err != nil {
		log.Print(err)
		return false, err
	}
	var enabled bool
	err = self.networkAbi.Unpack(&enabled, "enabled", enabledByte)
	if err != nil {
		log.Print(err)
		return false, err
	}
	return enabled, nil
}

func (self *Ethereum) EncodeMaxGasPrice() (string, error) {
	encodedData, err := self.networkAbi.Pack("maxGasPrice")
	if err != nil {
		log.Print(err)
		return "", err
	}
	return common.Bytes2Hex(encodedData), nil
}

func (self *Ethereum) ExtractMaxGasPrice(result string) (string, error) {
	gasByte, err := hexutil.Decode(result)
	if err != nil {
		log.Print(err)
		return "", err
	}
	var gasPrice *big.Int
	err = self.networkAbi.Unpack(&gasPrice, "maxGasPrice", gasByte)
	if err != nil {
		log.Print(err)
		return "", err
	}
	return gasPrice.String(), nil
}

func (self *Ethereum) ExtractRateData(result string, sourceArr []string, destAddr []string) (*[]ethereum.Rate, error) {
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
		dest := destAddr[i]
		rate := rateWapper.ExpectedRate[i]
		minRate := rateWapper.SlippageRate[i]
		rateReturn = append(rateReturn, ethereum.Rate{
			source, dest, rate.String(), minRate.String(),
		})
	}
	return &rateReturn, nil
}

func (self *Ethereum) ReadEventsWithBlockNumber(eventRaw *[]ethereum.EventRaw, latestBlock string) (*[]ethereum.EventHistory, error) {
	//get latestBlock to calculate timestamp
	events, err := self.ReadEvents(eventRaw, "node", latestBlock)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return events, nil
}

func (self *Ethereum) ReadEventsWithTimeStamp(eventRaw *[]ethereum.EventRaw) (*[]ethereum.EventHistory, error) {
	//get latestBlock to calculate timestamp
	events, err := self.ReadEvents(eventRaw, "etherscan", "0")
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return events, nil
}

type LogData struct {
	Source           common.Address `json:"source"`
	Dest             common.Address `json:"dest"`
	ActualSrcAmount  *big.Int       `json:"actualSrcAmount"`
	ActualDestAmount *big.Int       `json:"actualDestAmount"`
}

func (self *Ethereum) ReadEvents(listEventAddr *[]ethereum.EventRaw, typeFetch string, latestBlock string) (*[]ethereum.EventHistory, error) {
	listEvent := *listEventAddr
	endIndex := len(listEvent) - 1
	// var beginIndex = 0
	// if endIndex > 4 {
	// 	beginIndex = endIndex - 4
	// }

	index := 0
	events := make([]ethereum.EventHistory, 0)
	for i := endIndex; i >= 0; i-- {
		if index >= 5 {
			break
		}
		//filter amount
		isSmallAmount, err := self.IsSmallAmount(listEvent[i])
		if err != nil {
			log.Print(err)
			continue
		}
		if isSmallAmount {
			continue
		}

		txHash := listEvent[i].Txhash

		blockNumber, err := hexutil.DecodeBig(listEvent[i].BlockNumber)
		if err != nil {
			log.Print(err)
			continue
		}

		var timestamp string
		if typeFetch == "etherscan" {
			timestampHex, err := hexutil.DecodeBig(listEvent[i].Timestamp)
			if err != nil {
				log.Print(err)
				continue
			}
			timestamp = timestampHex.String()
			//fmt.Println(timestamp)
		} else {
			timestamp, err = self.Gettimestamp(blockNumber.String(), latestBlock, self.averageBlockTime)
			if err != nil {
				log.Print(err)
				continue
			}
		}

		var logData LogData
		data, err := hexutil.Decode(listEvent[i].Data)
		if err != nil {
			log.Print(err)
			continue
		}
		//fmt.Print(listEvent[i].Data)
		err = self.networkAbi.Unpack(&logData, "ExecuteTrade", data)
		if err != nil {
			log.Print(err)
			continue
		}

		actualDestAmount := logData.ActualDestAmount.String()
		actualSrcAmount := logData.ActualSrcAmount.String()
		dest := logData.Dest.String()
		source := logData.Source.String()

		events = append(events, ethereum.EventHistory{
			actualDestAmount, actualSrcAmount, dest, source, blockNumber.String(), txHash, timestamp,
		})
		index++
	}
	return &events, nil
}

func (self *Ethereum) IsSmallAmount(eventRaw ethereum.EventRaw) (bool, error) {
	data, err := hexutil.Decode(eventRaw.Data)
	if err != nil {
		log.Print(err)
		return true, err
	}
	var logData LogData
	err = self.networkAbi.Unpack(&logData, "ExecuteTrade", data)
	if err != nil {
		log.Print(err)
		return true, err
	}

	source := logData.Source
	var amount *big.Int
	if strings.ToLower(source.String()) == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		amount = logData.ActualSrcAmount
	} else {
		amount = logData.ActualDestAmount
	}

	//fmt.Print(amount.String())
	//check amount is greater than episol
	// amountBig, ok := new(big.Int).SetString(amount)
	// if !ok {
	// 	errorAmount := errors.New("Cannot read amount as number")
	// 	log.Print(errorAmount)
	// 	return false, errorAmount
	// }
	var episol, weight = big.NewInt(10), big.NewInt(15)
	episol.Exp(episol, weight, nil)
	//fmt.Print(episol.String())

	//fmt.Print(amount.Cmp(episol))
	if amount.Cmp(episol) == -1 {
		return true, nil
	}
	return false, nil
}

func (self *Ethereum) Gettimestamp(block string, latestBlock string, averageBlockTime int64) (string, error) {
	fromBlock, err := strconv.ParseInt(block, 10, 64)
	if err != nil {
		log.Print(err)
		return "", err
	}
	toBlock, err := strconv.ParseInt(latestBlock, 10, 64)
	if err != nil {
		log.Print(err)
		return "", err
	}
	timeNow := time.Now().Unix()
	timeStamp := timeNow - averageBlockTime*(toBlock-fromBlock)/1000

	timeStampBig := big.NewInt(timeStamp)
	return timeStampBig.String(), nil
}
