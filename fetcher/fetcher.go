package fetcher

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/KyberNetwork/server-go/ethereum"
)

type Token struct {
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Address string `json:"address"`
	Decimal int    `json:"decimal"`
	UsdId   string `json:"usd_id"`
}

type ServerLog struct {
	Url    string `json:"url"`
	ApiKey string `json:"api_key"`
}

type InfoData struct {
	ApiUsd    string           `json:"api_usd"`
	Tokens    map[string]Token `json:"tokens"`
	ServerLog ServerLog        `json:"server_logs"`

	NodeEndpoint string `json:"node_endpoint"`

	Network    string `json:"network"`
	NetworkAbi string
	TradeTopic string `json:"trade_topic"`

	Wapper     string `json:"wrapper"`
	WrapperAbi string
	EthAdress  string
	EthSymbol  string

	averageBlockTime int64 `json:"averageBlockTime"`
}

type Fetcher struct {
	info     *InfoData
	ethereum *Ethereum
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
	ethereum, err := NewEthereum(infoData.NodeEndpoint, infoData.Network, infoData.NetworkAbi, infoData.TradeTopic,
		infoData.Wapper, infoData.WrapperAbi)
	if err != nil {
		log.Print(err)
	}
	fetcher := &Fetcher{
		info:     &infoData,
		ethereum: ethereum,
	}
	//reader info from json

	return fetcher, nil
}

func (self *Fetcher) GetRateUsd() ([]io.ReadCloser, error) {
	outPut := make([]io.ReadCloser, 0)
	for _, token := range self.info.Tokens {
		usdId := token.UsdId
		response, err := http.Get(self.info.ApiUsd + "/v1/ticker/" + usdId)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		if response.StatusCode != 200 {
			log.Print(errors.New("Status code is 200"))
			return nil, errors.New("Status code is 200")
		}
		outPut = append(outPut, response.Body)
	}
	return outPut, nil
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

	// rates1, err := self.ethereum.GetRateFromNode(sourceSymbolArr, destSymbolArr, dataAbi)
	// if err != nil {
	// 	log.Print(err)
	// 	return nil, err
	// }
	// return rates1, nil

	url := self.info.ServerLog.Url + "/api?module=proxy&action=eth_call&to=" +
		self.info.Wapper + "&data=" + dataAbi + "&tag=latest&apikey=" + self.info.ServerLog.ApiKey
	response, err := http.Get(url)
	if err != nil {
		log.Print(err)
		//get rate from node
		rates, err := self.ethereum.GetRateFromNode(sourceSymbolArr, destSymbolArr, dataAbi)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		return rates, nil

	}
	if response.StatusCode != 200 {
		return nil, errors.New("Status code is 200")
	}
	rates, err := self.ethereum.ExactRateDataFromEtherscan(response.Body, sourceSymbolArr, destSymbolArr)
	if err != nil {
		log.Print(err)
		//get rate from node
		rates, err := self.ethereum.GetRateFromNode(sourceSymbolArr, destSymbolArr, dataAbi)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		return rates, nil
	}
	return rates, nil
}

func (self *Fetcher) GetLatestBlock() (string, error) {
	//get from etherscan
	response, err := http.Get(self.info.ServerLog.Url + "/api?module=proxy&action=eth_blockNumber")
	if err != nil {
		log.Print(err)
		blockNumber, err := self.ethereum.GetLatestBlock()
		if err != nil {
			log.Print(err)
			return "", err
		}
		return blockNumber, nil
	}
	//exact block number
	blockNumber, err := self.ethereum.ExactBlockNumber(response.Body)
	if err != nil {
		log.Print(err)
		//get from node
		blockNumber, err := self.ethereum.GetLatestBlock()
		if err != nil {
			log.Print(err)
			return "", err
		}
		return blockNumber, nil
	}
	return blockNumber, nil
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

	url := self.info.ServerLog.Url + "/api?module=logs&action=getLogs&fromBlock=" +
		fromBlock + "&toBlock=" + toBlock + "&address=" + self.info.Network + "&topic0=" +
		self.info.TradeTopic + "&apikey=" + self.info.ServerLog.ApiKey
	response, err := http.Get(url)
	if err != nil {
		log.Print(err)
		//get events from node
		events, err := self.ethereum.GetEventsFromNode(fromBlock, toBlock, self.info.averageBlockTime)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		return events, nil
	}
	if response.StatusCode != 200 {
		return nil, errors.New("Status code is 200")
	}
	events, err := self.ethereum.ExactEventFromEtherscan(response.Body)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return events, nil
}
