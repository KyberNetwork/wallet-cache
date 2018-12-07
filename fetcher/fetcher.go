package fetcher

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"sync"

	// "strconv"
	"time"

	"github.com/KyberNetwork/server-go/ethereum"
)

const (
	ETH_TO_WEI = 1000000000000000000
	MIN_ETH    = 0.001
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
	mu            *sync.RWMutex
	ApiUsd        string                    `json:"api_usd"`
	Tokens        map[string]ethereum.Token `json:"tokens"`
	TokenSnapshot map[string]ethereum.Token
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

	TrackerEndpoint string `json:"tracker_endpoint"`
	ConfigEndpoint  string `json:"config_endpoint"`
}

func (self *InfoData) UpdateListToken(tokens map[string]ethereum.Token) {
	self.mu.Lock()
	defer self.mu.Unlock()
	currentListToken := self.Tokens
	finalListToken := make(map[string]ethereum.Token)
	for symbol, token := range tokens {
		if currentToken, ok := currentListToken[symbol]; ok {
			if token.CGId == "" {
				token.CGId = currentToken.CGId
			}
		}
		finalListToken[symbol] = token
	}
	self.Tokens = finalListToken
}

func (self *InfoData) UpdateByBackupToken() {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.Tokens = self.TokenSnapshot
}

func (self *InfoData) GetListToken() map[string]ethereum.Token {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.Tokens
}

type ResultRpc struct {
	Result string `json:"result"`
}

type Fetcher struct {
	info     *InfoData
	ethereum *Ethereum
	fetIns   []FetcherInterface
}

func (self *Fetcher) GetNumTokens() int {
	listTokens := self.GetListToken()
	return len(listTokens)
}

func NewFetcher() (*Fetcher, error) {
	var file []byte
	var err error
	kyberENV := os.Getenv("KYBER_ENV")
	switch kyberENV {
	case "semi_production":
		file, err = ioutil.ReadFile("env/semi_production.json")
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
		break
	case "ropsten":
		file, err = ioutil.ReadFile("env/ropsten.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
		break
	case "rinkeby":
		file, err = ioutil.ReadFile("env/rinkeby.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
		break
	default:
		file, err = ioutil.ReadFile("env/ropsten.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
		break
	}

	mu := &sync.RWMutex{}

	infoData := InfoData{
		mu:         mu,
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

	infoData.TokenSnapshot = infoData.Tokens

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

	return fetcher, nil
}

func (self *Fetcher) TryUpdateListToken() error {
	var err error
	for i := 0; i < 3; i++ {
		err = self.UpdateListToken()
		if err != nil {
			log.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		return nil
	}
	self.info.UpdateByBackupToken()
	return nil
}

func (self *Fetcher) UpdateListToken() error {
	var err error
	result := make(map[string]ethereum.Token)
	for _, fetIns := range self.fetIns {
		result, err = fetIns.GetListToken(self.info.ConfigEndpoint)
		if err != nil {
			log.Print(err)
			continue
		}
		break
	}
	if err == nil {
		self.info.UpdateListToken(result)
		// err = storeConfig(listToken)
		// if err != nil {
		// 	log.Println("fetch tokens success but write to file js failed")
		// }
	}
	return err
}

func (self *Fetcher) GetCurrentListToken() map[string]ethereum.Token {
	return self.GetListToken()
}

// store config to a file js
// func storeConfig(tokens map[string]ethereum.Token) error {
// 	fileJS, err := os.Create("config/tokens.js")
// 	if err != nil {
// 		return err
// 	}
// 	defer fileJS.Close()
// 	bytes, err := json.Marshal(tokens)
// 	if err != nil {
// 		return err
// 	}
// 	stringFile := fmt.Sprintf("var configTokens = %s;", bytes)
// 	stringReader := strings.NewReader(stringFile)
// 	_, err = stringReader.WriteTo(fileJS)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (self *Fetcher) GetListToken() map[string]ethereum.Token {
	return self.info.GetListToken()
}

func (self *Fetcher) GetRateUsd() ([]io.ReadCloser, error) {
	usdId := make([]string, 0)
	listTokens := self.GetListToken()
	for _, token := range listTokens {
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
	listTokens := self.GetListToken()
	for _, token := range listTokens {
		if token.UsdId != "" {
			//usdId = append(usdId, token.UsdId)
			for _, fetIns := range self.fetIns {
				result, err := fetIns.GetGeneralInfo(token.CGId)
				if err != nil {
					log.Print(err)
					continue
				}
				//return result, nil
				generalInfo[token.Symbol] = result
				break
			}
			time.Sleep(5 * time.Second)
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

func getAmountInWei(amount float64) *big.Int {
	amountFloat := big.NewFloat(amount)
	ethFloat := big.NewFloat(ETH_TO_WEI)
	weiFloat := big.NewFloat(0).Mul(amountFloat, ethFloat)
	amoutInt, _ := weiFloat.Int(nil)
	return amoutInt
}

func getAmountTokenWithMinETH(rate *big.Int, decimal int) *big.Int {
	rFloat := big.NewFloat(0).SetInt(rate)
	ethFloat := big.NewFloat(ETH_TO_WEI)
	amoutnToken1ETH := rFloat.Quo(rFloat, ethFloat)
	minAmountWithMinETH := amoutnToken1ETH.Mul(amoutnToken1ETH, big.NewFloat(MIN_ETH))
	decimalWei := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
	amountWithDecimal := big.NewFloat(0).Mul(minAmountWithMinETH, big.NewFloat(0).SetInt(decimalWei))
	amountInt, _ := amountWithDecimal.Int(nil)
	return amountInt
}

func (self *Fetcher) GetRate(rates *[]ethereum.Rate) (*[]ethereum.Rate, error) {
	//append rate
	listTokens := self.GetListToken()
	sourceAddr := make([]string, 0)
	sourceSymbol := make([]string, 0)
	destAddr := make([]string, 0)
	destSymbol := make([]string, 0)
	amount := make([]*big.Int, 0)
	amountETH := make([]*big.Int, 0)
	ethSymbol := self.info.EthSymbol
	ethAddr := self.info.EthAdress
	minAmountETH := getAmountInWei(MIN_ETH)

	if len(*rates) > 0 && len(*rates) == len(listTokens)*2 {
		for _, rate := range *rates {
			if rate.Source == "ETH" && rate.Dest == "ETH" {
				continue
			}
			if rate.Source == "ETH" {
				amountToken := big.NewInt(0)
				r := big.NewInt(0)
				r.SetString(rate.Rate, 10)
				destSym := rate.Dest
				if destToken, ok := listTokens[destSym]; ok {
					decimal := destToken.Decimal
					if decimal != 0 || r.Cmp(amountToken) != 0 {
						amountToken = getAmountTokenWithMinETH(r, decimal)
					}
					amount = append(amount, amountToken)
					tokenAddr := destToken.Address
					sourceAddr = append(sourceAddr, tokenAddr)
					sourceSymbol = append(sourceSymbol, destSym)
				}
			} else {
				if _, ok := listTokens[rate.Source]; ok {
					destAddr = append(destAddr, ethAddr)
					destSymbol = append(destSymbol, ethSymbol)
					amountETH = append(amountETH, minAmountETH)
				}
			}
		}
		sourceAddr = append(sourceAddr, ethAddr)
		destAddr = append(destAddr, ethAddr)
		sourceSymbol = append(sourceSymbol, ethSymbol)
		destSymbol = append(destSymbol, ethSymbol)
		amount = append(amount, minAmountETH)
		amountETH = append(amountETH, minAmountETH)
	} else {
		for _, token := range listTokens {
			sourceAddr = append(sourceAddr, token.Address)
			sourceSymbol = append(sourceSymbol, token.Symbol)
			destAddr = append(destAddr, ethAddr)
			destSymbol = append(destSymbol, ethSymbol)
			amount = append(amount, big.NewInt(0))
			amountETH = append(amountETH, minAmountETH)
		}
	}
	sourceArr := append(sourceAddr, destAddr...)
	sourceSymbolArr := append(sourceSymbol, destSymbol...)
	destArr := append(destAddr, sourceAddr...)
	destSymbolArr := append(destSymbol, sourceSymbol...)
	amountArr := append(amount, amountETH...)

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

// func (self *Fetcher) GetEvents(blockNum string) (*[]ethereum.EventHistory, error) {
// 	toBlock := blockNum

// 	blockInt, err := strconv.Atoi(blockNum)
// 	if err != nil {
// 		log.Print(err)
// 		return nil, err
// 	}
// 	if blockInt > 5000 {
// 		blockInt = blockInt - 5000
// 	} else {
// 		blockInt = 0
// 	}
// 	fromBlock := strconv.Itoa(blockInt)

// 	for _, fetIns := range self.fetIns {
// 		eventRaw, err := fetIns.GetEvents(fromBlock, toBlock, self.info.Network, self.info.TradeTopic)
// 		if err != nil {
// 			log.Print(err)
// 			continue
// 		}
// 		var events *[]ethereum.EventHistory
// 		var errorEvent error
// 		if fetIns.GetTypeName() == "node" {
// 			latestBlock, err := self.GetLatestBlock()
// 			if err != nil {
// 				log.Print(err)
// 				continue
// 			}
// 			events, errorEvent = self.ethereum.ReadEventsWithBlockNumber(eventRaw, latestBlock)
// 		} else {
// 			events, errorEvent = self.ethereum.ReadEventsWithTimeStamp(eventRaw)
// 		}
// 		if errorEvent != nil {
// 			log.Print(errorEvent)
// 			continue
// 		}
// 		return events, nil
// 	}
// 	return nil, errors.New("Cannot get events")
// }

func (self *Fetcher) FetchTrackerData() (map[string]*ethereum.Rates, error) {
	for _, fetIns := range self.fetIns {
		result, err := fetIns.GetTrackerData(self.info.TrackerEndpoint)
		if err != nil {
			log.Print(err)
			continue
		}
		return result, nil
	}
	return nil, errors.New("Cannot get data from tracker")
}
