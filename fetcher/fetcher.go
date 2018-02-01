package fetcher

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/KyberNetwork/server-go/ethereum"
)

type Token struct {
	name    string
	symbol  string
	address string
	decimal int
	usdId   string
}

type InfoData struct {
	apiUsd       string
	tokens       []Token
	serverLogUrl string
	serverLogKey string

	nodeEndpoint string

	network    string
	networkAbi string
	tradeTopic string

	wapper     string
	wrapperAbi string
	ethAdress  string
	ethSymbol  string

	averageBlockTime int64
}

type Fetcher struct {
	info     *InfoData
	ethereum *Ethereum
}

func NewFetcher() (*Fetcher, error) {
	tokens := make([]Token, 0)
	tokens = append(tokens, Token{
		name:    "Ethereum",
		symbol:  "ETH",
		address: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		decimal: 18,
		usdId:   "ethereum",
	})
	tokens = append(tokens, Token{
		name:    "KyberNetwork",
		symbol:  "KNC",
		address: "0xdd974D5C2e2928deA5F71b9825b8b646686BD200",
		decimal: 18,
		usdId:   "kyber-network",
	})

	tokens = append(tokens, Token{
		name:    "OmiseGO",
		symbol:  "OMG",
		address: "0xd26114cd6ee289accf82350c8d8487fedb8a0c07",
		decimal: 18,
		usdId:   "omisego",
	})

	tokens = append(tokens, Token{
		name:    "Eos",
		symbol:  "EOS",
		address: "0x86fa049857e0209aa7d9e616f7eb3b3b78ecfdb0",
		decimal: 18,
		usdId:   "eos",
	})

	tokens = append(tokens, Token{
		name:    "Salt",
		address: "0x4156d3342d5c385a87d264f90653733592000581",
		symbol:  "SALT",
		decimal: 8,
		usdId:   "salt",
	})

	tokens = append(tokens, Token{
		name:    "Status",
		address: "0x744d70fdbe2ba4cf95131626614a1763df805b9e",
		symbol:  "SNT",
		decimal: 18,
		usdId:   "status",
	})

	infoData := &InfoData{
		apiUsd:       "https://api.coinmarketcap.com/v1/ticker/",
		tokens:       tokens,
		serverLogUrl: "https://api.etherscan.io",
		serverLogKey: "D8YAEQ3V4THAPDA9YSB1YGA1QY9KAMHY6M",

		nodeEndpoint: "https://mainnet.infura.io/DtzEYY0Km2BA3YwyJcBG",

		network:    "0x6bc0e45a62e952171fbd7562e7b5f30c50e564ac",
		networkAbi: `[{"constant":false,"inputs":[{"name":"alerter","type":"address"}],"name":"removeAlerter","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"reserve","type":"address"},{"name":"src","type":"address"},{"name":"dest","type":"address"},{"name":"add","type":"bool"}],"name":"listPairForReserve","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"bytes32"}],"name":"perReserveListedPairs","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getReserves","outputs":[{"name":"","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"enabled","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"pendingAdmin","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getOperators","outputs":[{"name":"","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"token","type":"address"},{"name":"amount","type":"uint256"},{"name":"sendTo","type":"address"}],"name":"withdrawToken","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"maxGasPrice","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newAlerter","type":"address"}],"name":"addAlerter","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"negligibleRateDiff","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"feeBurnerContract","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"expectedRateContract","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"whiteListContract","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"user","type":"address"}],"name":"getUserCapInWei","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newAdmin","type":"address"}],"name":"transferAdmin","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_enable","type":"bool"}],"name":"setEnable","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"claimAdmin","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"isReserve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getAlerters","outputs":[{"name":"","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"src","type":"address"},{"name":"dest","type":"address"},{"name":"srcQty","type":"uint256"}],"name":"getExpectedRate","outputs":[{"name":"expectedRate","type":"uint256"},{"name":"slippageRate","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"reserves","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newOperator","type":"address"}],"name":"addOperator","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"reserve","type":"address"},{"name":"add","type":"bool"}],"name":"addReserve","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"operator","type":"address"}],"name":"removeOperator","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_whiteList","type":"address"},{"name":"_expectedRate","type":"address"},{"name":"_feeBurner","type":"address"},{"name":"_maxGasPrice","type":"uint256"},{"name":"_negligibleRateDiff","type":"uint256"}],"name":"setParams","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"src","type":"address"},{"name":"dest","type":"address"},{"name":"srcQty","type":"uint256"}],"name":"findBestRate","outputs":[{"name":"","type":"uint256"},{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"src","type":"address"},{"name":"srcAmount","type":"uint256"},{"name":"dest","type":"address"},{"name":"destAddress","type":"address"},{"name":"maxDestAmount","type":"uint256"},{"name":"minConversionRate","type":"uint256"},{"name":"walletId","type":"address"}],"name":"trade","outputs":[{"name":"","type":"uint256"}],"payable":true,"stateMutability":"payable","type":"function"},{"constant":false,"inputs":[{"name":"amount","type":"uint256"},{"name":"sendTo","type":"address"}],"name":"withdrawEther","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getNumReserves","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"token","type":"address"},{"name":"user","type":"address"}],"name":"getBalance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"admin","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_admin","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":true,"name":"sender","type":"address"},{"indexed":false,"name":"amount","type":"uint256"}],"name":"EtherReceival","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"source","type":"address"},{"indexed":false,"name":"dest","type":"address"},{"indexed":false,"name":"actualSrcAmount","type":"uint256"},{"indexed":false,"name":"actualDestAmount","type":"uint256"}],"name":"ExecuteTrade","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"reserve","type":"address"},{"indexed":false,"name":"add","type":"bool"}],"name":"AddReserveToNetwork","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"reserve","type":"address"},{"indexed":false,"name":"src","type":"address"},{"indexed":false,"name":"dest","type":"address"},{"indexed":false,"name":"add","type":"bool"}],"name":"ListReservePairs","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"token","type":"address"},{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"sendTo","type":"address"}],"name":"TokenWithdraw","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"sendTo","type":"address"}],"name":"EtherWithdraw","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"pendingAdmin","type":"address"}],"name":"TransferAdminPending","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"newAdmin","type":"address"},{"indexed":false,"name":"previousAdmin","type":"address"}],"name":"AdminClaimed","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"newAlerter","type":"address"},{"indexed":false,"name":"isAdd","type":"bool"}],"name":"AlerterAdded","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"newOperator","type":"address"},{"indexed":false,"name":"isAdd","type":"bool"}],"name":"OperatorAdded","type":"event"}]`,
		tradeTopic: "0x1849bd6a030a1bca28b83437fd3de96f3d27a5d172fa7e9c78e7b61468928a39",

		wapper:           "0x28e9221a3bb961bd89641ea55b4dfa4ae464e9cf",
		wrapperAbi:       `[{"constant":true,"inputs":[{"name":"x","type":"bytes14"},{"name":"byteInd","type":"uint256"}],"name":"getInt8FromByte","outputs":[{"name":"","type":"int8"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"reserve","type":"address"},{"name":"tokens","type":"address[]"}],"name":"getBalances","outputs":[{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"pricingContract","type":"address"},{"name":"tokenList","type":"address[]"}],"name":"getTokenIndicies","outputs":[{"name":"","type":"uint256[]"},{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"x","type":"bytes14"},{"name":"byteInd","type":"uint256"}],"name":"getByteFromBytes14","outputs":[{"name":"","type":"bytes1"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"network","type":"address"},{"name":"sources","type":"address[]"},{"name":"dests","type":"address[]"},{"name":"qty","type":"uint256[]"}],"name":"getExpectedRates","outputs":[{"name":"expectedRate","type":"uint256[]"},{"name":"slippageRate","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"pricingContract","type":"address"},{"name":"tokenList","type":"address[]"}],"name":"getTokenRates","outputs":[{"name":"","type":"uint256[]"},{"name":"","type":"uint256[]"},{"name":"","type":"int8[]"},{"name":"","type":"int8[]"},{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"}]`,
		ethAdress:        "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		ethSymbol:        "ETH",
		averageBlockTime: 16,
	}

	ethereum, err := NewEthereum(infoData.nodeEndpoint, infoData.network, infoData.networkAbi, infoData.tradeTopic, infoData.wapper, infoData.wrapperAbi)
	if err != nil {
		log.Print(err)
	}
	fetcher := &Fetcher{
		info:     infoData,
		ethereum: ethereum,
	}
	//reader info from json

	return fetcher, nil
}

func (self *Fetcher) GetRateUsd() ([]io.ReadCloser, error) {
	outPut := make([]io.ReadCloser, 0)
	for _, token := range self.info.tokens {
		usdId := token.usdId
		response, err := http.Get(self.info.apiUsd + usdId)
		if err != nil {
			return nil, err
		}
		if response.StatusCode != 200 {
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
	for _, token := range self.info.tokens {
		sourceAddr = append(sourceAddr, token.address)
		sourceSymbol = append(sourceSymbol, token.symbol)
		destAddr = append(destAddr, self.info.ethAdress)
		destSymbol = append(destSymbol, self.info.ethSymbol)
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

	url := self.info.serverLogUrl + "/api?module=proxy&action=eth_call&to=" + self.info.wapper + "&data=" + dataAbi + "&tag=latest&apikey=" + self.info.serverLogKey
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
	response, err := http.Get(self.info.serverLogUrl + "/api?module=proxy&action=eth_blockNumber")
	if err != nil {
		log.Print(err)
		return "", err
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

	url := self.info.serverLogUrl + "/api?module=logs&action=getLogs&fromBlock=" + fromBlock + "&toBlock=" + toBlock + "&address=" + self.info.network + "&topic0=" + self.info.tradeTopic + "&apikey=" + self.info.serverLogKey
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
