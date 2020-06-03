package fetcher

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"sync"

	"time"

	"github.com/KyberNetwork/cache/common"
	"github.com/KyberNetwork/cache/ethereum"
)

const (
	ETH_TO_WEI = 1000000000000000000
	MIN_ETH    = 0.1
	KEY        = "kybersecret"

	timeW8Req         = 500
	timeW8CheckStatus = 3 * time.Second
)

type Connection struct {
	Endpoint string `json:"endPoint"`
	Type     string `json:"type"`
	Apikey   string `json:"api_key"`
}

type InfoData struct {
	mu             *sync.RWMutex
	ApiUsd         string              `json:"api_usd"`
	CoinMarket     []string            `json:"coin_market"`
	TokenAPI       []ethereum.TokenAPI `json:"tokens"`
	CanDeleteToken []string            `json:"can_delete"`

	Tokens       map[string]ethereum.Token
	BackupTokens map[string]ethereum.Token
	arrToken     []ethereum.Token

	Connections []Connection `json:"connections"`

	Network    string `json:"network"`
	NetworkAbi string
	TradeTopic string `json:"trade_topic"`

	Wapper     string `json:"wrapper"`
	WrapperAbi string
	EthAdress  string
	EthSymbol  string

	AverageBlockTime int64 `json:"averageBlockTime"`

	GasStationEndpoint string `json:"gasstation_endpoint"`
	APIEndpoint        string `json:"api_endpoint"`
	ConfigEndpoint     string `json:"config_endpoint"`
	UserStatsEndpoint  string `json:"user_stats_endpoint"`
}

func (self *InfoData) GetListToken() map[string]ethereum.Token {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.Tokens
}

func (self *InfoData) GetArrToken() []ethereum.Token {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.arrToken
}

func (self *InfoData) UpdateByBackupToken() {
	self.mu.Lock()
	defer self.mu.Unlock()
	mapBackUp := make(map[string]ethereum.Token)
	for k, v := range self.BackupTokens {
		mapBackUp[k] = v
	}
	self.Tokens = mapBackUp
}

func (self *InfoData) UpdateListToken(tokens map[string]ethereum.Token, arrToken []ethereum.Token) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.Tokens = tokens
	self.arrToken = arrToken
}

func (self *InfoData) GetTokenAPI() []ethereum.TokenAPI {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.TokenAPI
}

type Fetcher struct {
	info        *InfoData
	ethereum    *Ethereum
	fetIns      []FetcherInterface
	httpFetcher *HTTPFetcher
}

func (self *Fetcher) GetNumTokens() int {
	listTokens := self.GetListToken()
	return len(listTokens)
}

func NewFetcher(kyberENV string) (*Fetcher, error) {
	var file []byte
	var err error

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

	listToken := make(map[string]ethereum.Token)
	listBackup := make(map[string]ethereum.Token)
	infoData.Tokens = listToken
	for _, t := range infoData.TokenAPI {
		listBackup[t.Symbol] = ethereum.TokenAPIToToken(t)
	}
	infoData.BackupTokens = listBackup

	fetIns := make([]FetcherInterface, 0)
	for _, connection := range infoData.Connections {
		newFetcher, err := NewFetcherIns(connection.Type, connection.Endpoint, connection.Apikey)
		if err != nil {
			log.Print(err)
		} else {
			fetIns = append(fetIns, newFetcher)
		}
	}

	httpFetcher := NewHTTPFetcher(infoData.ConfigEndpoint, infoData.GasStationEndpoint, infoData.APIEndpoint)

	ethereum, err := NewEthereum(infoData.Network, infoData.NetworkAbi, infoData.TradeTopic,
		infoData.Wapper, infoData.WrapperAbi, infoData.AverageBlockTime)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	fetcher := &Fetcher{
		info:        &infoData,
		ethereum:    ethereum,
		fetIns:      fetIns,
		httpFetcher: httpFetcher,
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
	var (
		err    error
		result []ethereum.Token
	)
	result, err = self.httpFetcher.GetListToken()
	if err != nil {
		log.Println(err)
		return err
	}
	listToken := make(map[string]ethereum.Token)
	for _, token := range result {
		if token.DelistTime == 0 || uint64(time.Now().UTC().Unix()) <= TIME_TO_DELETE+token.DelistTime {
			tokenID := token.Symbol
			if token.TokenID != "" {
				tokenID = token.TokenID
			}
			newToken := token
			newToken.TokenID = tokenID
			listToken[tokenID] = newToken
		}
	}
	if common.IsDifferentMapToken(listToken, self.GetListToken()) {
		self.info.UpdateListToken(listToken, result)
	}
	return nil
}

func (self *Fetcher) GetArrToken() []ethereum.Token {
	return self.info.GetArrToken()
}

// api to get config token
func (self *Fetcher) GetListTokenAPI() []ethereum.TokenAPI {
	return self.info.GetTokenAPI()
}

// GetListToken return map token with key is token ID
func (self *Fetcher) GetListToken() map[string]ethereum.Token {
	return self.info.GetListToken()
}

func (self *Fetcher) GetRateUsdEther() (string, error) {
	// rateUsd, err := self.marketFetIns.GetRateUsdEther()
	rateUsd, err := self.httpFetcher.GetRateUsdEther()
	if err != nil {
		log.Print(err)
		return "", err
	}
	return rateUsd, nil
}

func (self *Fetcher) GetGasPrice() (*ethereum.GasPrice, error) {
	result, err := self.httpFetcher.GetGasPrice()
	if err != nil {
		log.Print(err)
		return nil, errors.New("Cannot get gas price")
	}
	return result, nil
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

func (self *Fetcher) FetchUserInfo(address string) (*common.UserInfo, error) {
	userInfo, err := self.httpFetcher.GetUserInfo(self.info.UserStatsEndpoint + "users?address=" + address)
	if err != nil {
		return nil, errors.New("Cannot get user info")
	}
	return userInfo, nil
}

func (self *Fetcher) GetTokenBySymbol(symbol string) (*ethereum.Token, error) {
	self.info.mu.RLock()
	defer self.info.mu.RUnlock()
	for _, token := range self.info.Tokens {
		if token.Symbol == symbol {
			return &token, nil
		}
	}
	err := errors.New("Token is not existed")
	log.Println(err)
	return nil, err
}

func (self *Fetcher) FetchRate() ([]ethereum.Rate, error) {
	return self.httpFetcher.GetRate()
}

func (self *Fetcher) GetSourceAmount(src string, dest string, destAmount string) (string, error) {
	srcToken, err := self.GetTokenBySymbol(src)
	if err != nil {
		log.Println(err)
		return "", err
	}
	destToken, err := self.GetTokenBySymbol(dest)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return self.httpFetcher.GetQuoteAmount(srcToken.Address, destToken.Address, destAmount, "buy")
}
