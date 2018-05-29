package fetcher

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/KyberNetwork/server-go/ethereum"
)

// type Token struct {
// 	Name    string `json:"name"`
// 	Symbol  string `json:"symbol"`
// 	Address string `json:"address"`
// 	Decimal int    `json:"decimal"`
// 	UsdId   string `json:"usd_id"`
// }

type Connection struct {
	Endpoint string `json:"endPoint"`
	Type     string `json:"type"`
	Apikey   string `json:"api_key"`
}

type InfoData struct {
	ApiUsd string                    `json:"api_usd"`
	Tokens map[string]ethereum.Token `json:"tokens"`
	//ServerLog ServerLog        `json:"server_logs"`
	Connections []Connection `json:"connections"`

	//NodeEndpoint string `json:"node_endpoint"`

	Network    string `json:"network"`
	NetworkAbi string
	TradeTopic string `json:"trade_topic"`

	Wapper     string `json:"wrapper"`
	WrapperAbi string
	EthAdress  string
	EthSymbol  string

	AverageBlockTime int64 `json:"averageBlockTime"`
}

type ResultRpc struct {
	Result string `json:"result"`
}

type Fetcher struct {
	info     *InfoData
	ethereum *Ethereum
	fetIns   []FetcherInterface
}

func NewFetcher() (*Fetcher, error) {
	var file []byte
	var err error
	switch os.Getenv("KYBER_ENV") {
	case "internal_mainnet":
		file, err = ioutil.ReadFile("env/internal_mainnet.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
		break
	case "staging":
		file, err = ioutil.ReadFile("env/staging.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
		break
	case "production":
		file, err = ioutil.ReadFile("env/production.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
	case "kovan":
		file, err = ioutil.ReadFile("env/kovan.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
	case "ropsten":
		file, err = ioutil.ReadFile("env/ropsten.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
	case "production_test":
		file, err = ioutil.ReadFile("env/production_test.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
	default:
		file, err = ioutil.ReadFile("env/internal_mainnet.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
		break
	}

	infoData := InfoData{
		WrapperAbi: `[{"constant":true,"inputs":[{"name":"x","type":"bytes14"},{"name":"byteInd","type":"uint256"}],"name":"getInt8FromByte","outputs":[{"name":"","type":"int8"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"reserve","type":"address"},{"name":"tokens","type":"address[]"}],"name":"getBalances","outputs":[{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"pricingContract","type":"address"},{"name":"tokenList","type":"address[]"}],"name":"getTokenIndicies","outputs":[{"name":"","type":"uint256[]"},{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"x","type":"bytes14"},{"name":"byteInd","type":"uint256"}],"name":"getByteFromBytes14","outputs":[{"name":"","type":"bytes1"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"network","type":"address"},{"name":"sources","type":"address[]"},{"name":"dests","type":"address[]"},{"name":"qty","type":"uint256[]"}],"name":"getExpectedRates","outputs":[{"name":"expectedRate","type":"uint256[]"},{"name":"slippageRate","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"pricingContract","type":"address"},{"name":"tokenList","type":"address[]"}],"name":"getTokenRates","outputs":[{"name":"","type":"uint256[]"},{"name":"","type":"uint256[]"},{"name":"","type":"int8[]"},{"name":"","type":"int8[]"},{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"}]`,
		EthAdress:  "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		EthSymbol:  "ETH",
		NetworkAbi: `[{"constant":false,"inputs":[{"name":"alerter","type":"address"}],"name":"removeAlerter","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"reserve","type":"address"},{"name":"src","type":"address"},{"name":"dest","type":"address"},{"name":"add","type":"bool"}],"name":"listPairForReserve","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"bytes32"}],"name":"perReserveListedPairs","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getReserves","outputs":[{"name":"","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"enabled","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"pendingAdmin","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getOperators","outputs":[{"name":"","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"token","type":"address"},{"name":"amount","type":"uint256"},{"name":"sendTo","type":"address"}],"name":"withdrawToken","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"maxGasPrice","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newAlerter","type":"address"}],"name":"addAlerter","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"negligibleRateDiff","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"feeBurnerContract","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"expectedRateContract","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"whiteListContract","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"user","type":"address"}],"name":"getUserCapInWei","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newAdmin","type":"address"}],"name":"transferAdmin","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_enable","type":"bool"}],"name":"setEnable","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"claimAdmin","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"isReserve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getAlerters","outputs":[{"name":"","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"src","type":"address"},{"name":"dest","type":"address"},{"name":"srcQty","type":"uint256"}],"name":"getExpectedRate","outputs":[{"name":"expectedRate","type":"uint256"},{"name":"slippageRate","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"reserves","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newOperator","type":"address"}],"name":"addOperator","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"reserve","type":"address"},{"name":"add","type":"bool"}],"name":"addReserve","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"operator","type":"address"}],"name":"removeOperator","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_whiteList","type":"address"},{"name":"_expectedRate","type":"address"},{"name":"_feeBurner","type":"address"},{"name":"_maxGasPrice","type":"uint256"},{"name":"_negligibleRateDiff","type":"uint256"}],"name":"setParams","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"src","type":"address"},{"name":"dest","type":"address"},{"name":"srcQty","type":"uint256"}],"name":"findBestRate","outputs":[{"name":"","type":"uint256"},{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"src","type":"address"},{"name":"srcAmount","type":"uint256"},{"name":"dest","type":"address"},{"name":"destAddress","type":"address"},{"name":"maxDestAmount","type":"uint256"},{"name":"minConversionRate","type":"uint256"},{"name":"walletId","type":"address"}],"name":"trade","outputs":[{"name":"","type":"uint256"}],"payable":true,"stateMutability":"payable","type":"function"},{"constant":false,"inputs":[{"name":"amount","type":"uint256"},{"name":"sendTo","type":"address"}],"name":"withdrawEther","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getNumReserves","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"token","type":"address"},{"name":"user","type":"address"}],"name":"getBalance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"admin","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_admin","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":true,"name":"sender","type":"address"},{"indexed":false,"name":"amount","type":"uint256"}],"name":"EtherReceival","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"source","type":"address"},{"indexed":false,"name":"dest","type":"address"},{"indexed":false,"name":"actualSrcAmount","type":"uint256"},{"indexed":false,"name":"actualDestAmount","type":"uint256"}],"name":"ExecuteTrade","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"reserve","type":"address"},{"indexed":false,"name":"add","type":"bool"}],"name":"AddReserveToNetwork","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"reserve","type":"address"},{"indexed":false,"name":"src","type":"address"},{"indexed":false,"name":"dest","type":"address"},{"indexed":false,"name":"add","type":"bool"}],"name":"ListReservePairs","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"token","type":"address"},{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"sendTo","type":"address"}],"name":"TokenWithdraw","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"sendTo","type":"address"}],"name":"EtherWithdraw","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"pendingAdmin","type":"address"}],"name":"TransferAdminPending","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"newAdmin","type":"address"},{"indexed":false,"name":"previousAdmin","type":"address"}],"name":"AdminClaimed","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"newAlerter","type":"address"},{"indexed":false,"name":"isAdd","type":"bool"}],"name":"AlerterAdded","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"newOperator","type":"address"},{"indexed":false,"name":"isAdd","type":"bool"}],"name":"OperatorAdded","type":"event"}]`,
	}
	err = json.Unmarshal(file, &infoData)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	fetIns := make([]FetcherInterface, 0)
	for _, connection := range infoData.Connections {
		newFetcher, err := NewFetcherIns(connection.Type, connection.Endpoint, connection.Apikey)
		if err != nil {
			log.Print(err)
		} else {
			fetIns = append(fetIns, newFetcher)
		}
	}

	ethereum, err := NewEthereum(infoData.Network, infoData.NetworkAbi, infoData.TradeTopic,
		infoData.Wapper, infoData.WrapperAbi, infoData.AverageBlockTime)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	fetcher := &Fetcher{
		info:     &infoData,
		ethereum: ethereum,
		fetIns:   fetIns,
	}
	//reader info from json

	return fetcher, nil
}

func (self *Fetcher) GetListToken() map[string]ethereum.Token {
	return self.info.Tokens
}

func (self *Fetcher) GetRateUsd() ([]io.ReadCloser, error) {
	usdId := make([]string, 0)
	for _, token := range self.info.Tokens {
		if token.UsdId != "" {
			usdId = append(usdId, token.UsdId)
		}
	}
	for _, fetIns := range self.fetIns {
		result, err := fetIns.GetRateUsd(usdId)
		if err != nil {
			log.Print(err)
			continue
		}
		return result, nil
	}
	return nil, errors.New("Cannot get rate USD")
}

func (self *Fetcher) GetGeneralInfoTokens() map[string]*ethereum.TokenGeneralInfo {
	generalInfo := map[string]*ethereum.TokenGeneralInfo{}
	//	usdId := make([]string, 0)
	for _, token := range self.info.Tokens {
		if token.UsdId != "" {
			//usdId = append(usdId, token.UsdId)
			for _, fetIns := range self.fetIns {
				result, err := fetIns.GetGeneralInfo(token.UsdId)
				if err != nil {
					log.Print(err)
					continue
				}
				//return result, nil
				generalInfo[token.Symbol] = result
				break
			}

		}
	}

	return generalInfo
	// for _, fetIns := range self.fetIns {
	// 	result, err := fetIns.GetGeneralInfo(usdId)
	// 	if err != nil {
	// 		log.Print(err)
	// 		continue
	// 	}
	// 	return result, nil
	// }
	// return nil, errors.New("Cannot get rate USD")
}

func (self *Fetcher) GetRateUsdEther() (string, error) {
	//rateUsd, err := fetIns.GetRateUsdEther()

	// usdId := make([]string, 0)
	// for _, token := range self.info.Tokens {
	// 	usdId = append(usdId, token.UsdId)
	// }
	for _, fetIns := range self.fetIns {
		rateUsd, err := fetIns.GetRateUsdEther()
		//fmt.Print(rateUsd)
		if err != nil {
			log.Print(err)
			continue
		}
		return rateUsd, nil
	}
	return "", errors.New("Cannot get rate USD")
}

func (self *Fetcher) GetGasPrice() (*ethereum.GasPrice, error) {
	for _, fetIns := range self.fetIns {
		result, err := fetIns.GetGasPrice()
		if err != nil {
			log.Print(err)
			continue
		}
		return result, nil
	}
	return nil, errors.New("Cannot get gas price")
}

func (self *Fetcher) GetMaxGasPrice() (string, error) {
	dataAbi, err := self.ethereum.EncodeMaxGasPrice()
	if err != nil {
		log.Print(err)
		return "", err
	}
	for _, fetIns := range self.fetIns {
		result, err := fetIns.EthCall(self.info.Network, dataAbi)
		if err != nil {
			log.Print(err)
			continue
		}
		gasPrice, err := self.ethereum.ExtractMaxGasPrice(result)
		if err != nil {
			log.Print(err)
			continue
		}
		return gasPrice, nil
	}
	return "", errors.New("Cannot get gas price")
}

func (self *Fetcher) CheckKyberEnable() (bool, error) {
	dataAbi, err := self.ethereum.EncodeKyberEnable()
	if err != nil {
		log.Print(err)
		return false, err
	}
	for _, fetIns := range self.fetIns {
		result, err := fetIns.EthCall(self.info.Network, dataAbi)
		if err != nil {
			log.Print(err)
			continue
		}
		enabled, err := self.ethereum.ExtractEnabled(result)
		if err != nil {
			log.Print(err)
			continue
		}
		return enabled, nil
	}
	return false, errors.New("Cannot check kyber enable")
}

func (self *Fetcher) GetRate() (*[]ethereum.Rate, error) {
	//append rate
	sourceAddr := make([]string, 0)
	sourceSymbol := make([]string, 0)
	destAddr := make([]string, 0)
	destSymbol := make([]string, 0)
	amount := make([]int64, 0)
	for _, token := range self.info.Tokens {
		sourceAddr = append(sourceAddr, token.Address)
		sourceSymbol = append(sourceSymbol, token.Symbol)
		destAddr = append(destAddr, self.info.EthAdress)
		destSymbol = append(destSymbol, self.info.EthSymbol)
		amount = append(amount, 0)
	}
	sourceArr := append(sourceAddr, destAddr...)
	sourceSymbolArr := append(sourceSymbol, destSymbol...)
	destArr := append(destAddr, sourceAddr...)
	destSymbolArr := append(destSymbol, sourceSymbol...)
	amountArr := append(amount, amount...)

	dataAbi, err := self.ethereum.EncodeRateData(sourceArr, destArr, amountArr)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	for _, fetIns := range self.fetIns {
		result, err := fetIns.EthCall(self.info.Wapper, dataAbi)
		if err != nil {
			log.Print(err)
			continue
		}
		rates, err := self.ethereum.ExtractRateData(result, sourceSymbolArr, destSymbolArr)
		if err != nil {
			log.Print(err)
			continue
		}
		return rates, nil
	}
	return nil, errors.New("Cannot get rate")
}

func (self *Fetcher) GetLatestBlock() (string, error) {
	for _, fetIns := range self.fetIns {
		blockNo, err := fetIns.GetLatestBlock()
		if err != nil {
			log.Print(err)
			continue
		}
		return blockNo, nil
	}
	return "", errors.New("Cannot get block number")
}

func (self *Fetcher) GetEvents(blockNum string) (*[]ethereum.EventHistory, error) {
	toBlock := blockNum

	blockInt, err := strconv.Atoi(blockNum)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if blockInt > 5000 {
		blockInt = blockInt - 5000
	} else {
		blockInt = 0
	}
	fromBlock := strconv.Itoa(blockInt)

	for _, fetIns := range self.fetIns {
		eventRaw, err := fetIns.GetEvents(fromBlock, toBlock, self.info.Network, self.info.TradeTopic)
		if err != nil {
			log.Print(err)
			continue
		}
		var events *[]ethereum.EventHistory
		var errorEvent error
		if fetIns.GetTypeName() == "node" {
			latestBlock, err := self.GetLatestBlock()
			if err != nil {
				log.Print(err)
				continue
			}
			events, errorEvent = self.ethereum.ReadEventsWithBlockNumber(eventRaw, latestBlock)
		} else {
			events, errorEvent = self.ethereum.ReadEventsWithTimeStamp(eventRaw)
		}
		if errorEvent != nil {
			log.Print(errorEvent)
			continue
		}
		return events, nil
	}
	return nil, errors.New("Cannot get events")
}
