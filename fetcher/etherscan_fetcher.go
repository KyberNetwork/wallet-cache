package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/KyberNetwork/server-go/ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// api_key for tracker.kyber
const (
	API_KEY_TRACKER = "jHGlaMKcGn5cCBxQCGwusS4VcnH0C6tN"
	TIME_TO_DELETE  = 18000
	COINGECKO_API   = "https://api.coingecko.com/api/v3"
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

	b, err := HTTPCall(url)
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
	url := self.url + "/api?module=proxy&action=eth_blockNumber"

	b, err := HTTPCall(url)
	if err != nil {
		log.Print(err)
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

// func (self *Etherscan) GetEvents(fromBlock string, toBlock string, network string, tradeTopic string) (*[]ethereum.EventRaw, error) {
// 	url := self.url + "/api?module=logs&action=getLogs&fromBlock=" +
// 		fromBlock + "&toBlock=" + toBlock + "&address=" + network + "&topic0=" +
// 		tradeTopic + "&apikey=" + self.apiKey
// 	response, err := http.Get(url)
// 	if err != nil {
// 		log.Print(err)
// 		return nil, err
// 	}
// 	if response.StatusCode != 200 {
// 		return nil, errors.New("Status code is 200")
// 	}

// 	defer (response.Body).Close()
// 	b, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		log.Print(err)
// 		return nil, err
// 	}
// 	result := ResultEvent{}
// 	err = json.Unmarshal(b, &result)
// 	if err != nil {
// 		log.Print(err)
// 		return nil, err
// 	}

// 	return &result.Result, nil
// }

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
	gasStationURL := "https://ethgasstation.info/json/ethgasAPI.json"

	b, err := HTTPCall(gasStationURL)
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

// get data from tracker.kyber

func (self *Etherscan) GetTrackerData(trackerEndpoint string) (map[string]*ethereum.Rates, error) {
	trackerAPI := trackerEndpoint + "/api/tokens/rates?api_key=" + API_KEY_TRACKER
	b, err := HTTPCall(trackerAPI)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	trackerData := map[string]*ethereum.Rates{}
	err = json.Unmarshal(b, &trackerData)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return trackerData, nil
}

func (self *Etherscan) GetListToken(configEndpoint string) (map[string]ethereum.Token, error) {
	b, err := HTTPCall(configEndpoint)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	var result ethereum.TokenConfig
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if result.Success == false {
		err = errors.New("Cannot get list token")
		return nil, err
	}
	listToken := make(map[string]ethereum.Token)
	for _, token := range result.Data {
		if token.DelistTime == 0 || uint64(time.Now().UTC().Unix()) <= TIME_TO_DELETE+token.DelistTime {
			listToken[token.Symbol] = token
		}
	}
	return listToken, nil
}

// func (self *Etherscan) GetRateUsdEther() (string, error) {
// 	response, err := http.Get("https://api.coinmarketcap.com/v1/ticker/ethereum")
// 	if err != nil {
// 		log.Print(err)
// 		return "", err
// 	}
// 	defer (response.Body).Close()
// 	b, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		log.Print(err)
// 		return "", err
// 	}
// 	rateItem := make([]RateUSD, 0)
// 	err = json.Unmarshal(b, &rateItem)
// 	if err != nil {
// 		log.Print(err)
// 		return "", err
// 	}
// 	return rateItem[0].PriceUsd, nil
// }

// func (self *Etherscan) GetGeneralInfo(usdId string) (*ethereum.TokenGeneralInfo, error) {
// 	response, err := http.Get("https://api.coinmarketcap.com/v2/ticker/" + usdId + "/?convert=ETH")
// 	if err != nil {
// 		log.Print(err)
// 		return nil, err
// 	}
// 	defer (response.Body).Close()
// 	b, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		log.Print(err)
// 		return nil, err
// 	}
// 	tokenItem := map[string]ethereum.TokenGeneralInfo{}
// 	err = json.Unmarshal(b, &tokenItem)
// 	if err != nil {
// 		log.Print(err)
// 		return nil, err
// 	}

// 	if data, ok := tokenItem["data"]; ok {
// 		data.MarketCap = data.Quotes["ETH"].MarketCap
// 		return &data, nil
// 	}
// 	err = errors.New("Cannot find data key in return quotes of ticker")
// 	log.Print(err)
// 	return nil, err
// }

func (self *Etherscan) GetRateUsdEther() (string, error) {
	url := COINGECKO_API + "/coins/ethereum"
	b, err := HTTPCall(url)
	if err != nil {
		log.Print(err)
		return "", err
	}
	rateItem := ethereum.RateUSDCG{}
	err = json.Unmarshal(b, &rateItem)
	if err != nil {
		log.Print(err)
		return "", err
	}
	rateString := fmt.Sprintf("%.6f", rateItem.MarketData.CurrentPrice.USD)
	return rateString, nil
}

func (self *Etherscan) GetGeneralInfo(coinID string) (*ethereum.TokenGeneralInfo, error) {
	url := COINGECKO_API + "/coins/" + coinID
	log.Println("url: ", url)
	b, err := HTTPCall(url)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	tokenItem := ethereum.TokenInfoCoinGecko{}
	err = json.Unmarshal(b, &tokenItem)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	tokenGenalInfo := tokenItem.ToTokenInfoCMC()
	log.Println("token info: ", tokenItem)
	return &tokenGenalInfo, nil
	err = errors.New("Cannot find data key in return quotes of ticker")
	log.Print(err)
	return nil, err
}
