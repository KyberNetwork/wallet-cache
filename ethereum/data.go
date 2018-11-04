package ethereum

// const (
// 	TIME_TO_DELETE = 18000
// )

type EventRaw struct {
	Timestamp   string `json:"timestamp"`
	BlockNumber string `json:"blockNumber"`
	Txhash      string `json:"transactionHash"`
	Data        string `json:"data"`
}
type EventHistory struct {
	ActualDestAmount string `json:"actualDestAmount"`
	ActualSrcAmount  string `json:"actualSrcAmount"`
	Dest             string `json:"dest"`
	Source           string `json:"source"`

	BlockNumber string `json:"blockNumber"`
	Txhash      string `json:"txHash"`
	Timestamp   string `json:"timestamp"`
}

type Rate struct {
	Source  string `json:"source"`
	Dest    string `json:"dest"`
	Rate    string `json:"rate"`
	Minrate string `json:"minRate"`
}

type GasPrice struct {
	Fast     string `json:"fast"`
	Standard string `json:"standard"`
	Low      string `json:"low"`
	Default  string `json:"default"`
}

type Token struct {
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Address string `json:"address"`
	Decimal int    `json:"decimals"`
	UsdId   string `json:"cmc_id"`
	CGId    string `json:"cg_id"`
}

type TokenAPI struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Decimals    int    `json:"decimals"`
	UsdID       string `json:"cmc_id"`
	TimeListing uint64 `json:"listing_time,omitempty"`
	CGId        string `json:"cg_id"`
	// DelistTime  uint64 `json:"delist_time,omitempty"`
	ReserveSrc  []string `json:"reserves_src,omitempty"`
	ReserveDest []string `json:"reserves_dest,omitempty"`
}

func TokenAPIToToken(tokenAPI TokenAPI) Token {
	// if tokenAPI.DelistTime == 0 || uint64(time.Now().UTC().Unix()) <= TIME_TO_DELETE+tokenAPI.DelistTime {
	return Token{
		Name:    tokenAPI.Name,
		Symbol:  tokenAPI.Symbol,
		Address: tokenAPI.Address,
		Decimal: tokenAPI.Decimals,
		UsdId:   tokenAPI.UsdID,
		CGId:    tokenAPI.CGId,
	}
	// }
	// return nil
}

type QuoInfo struct {
	MarketCap float64 `json:"market_cap"`
	Volume24h float64 `json:"volume_24h"`
}

type TokenGeneralInfo struct {
	CirculatingSupply float64            `json:"circulating_supply"`
	TotalSupply       float64            `json:"total_supply"`
	MaxSupply         float64            `json:"max_supply"`
	MarketCap         float64            `json:"market_cap"`
	Quotes            map[string]QuoInfo `json:"quotes`
}

type CurrencyData struct {
	ETH float64 `json:"eth"`
	USD float64 `json:"usd"`
}

type TokenInfoCoinGecko struct {
	MarketData struct {
		MarketCap CurrencyData `json:"market_cap"`
		Volume24H CurrencyData `json:"total_volume"`
	} `json:"market_data"`
}

func TokenInfoCGToCMC(tokenInfo TokenInfoCoinGecko) TokenGeneralInfo {
	quotes := make(map[string]QuoInfo)
	quotes["ETH"] = QuoInfo{
		MarketCap: tokenInfo.MarketData.MarketCap.ETH,
		Volume24h: tokenInfo.MarketData.Volume24H.ETH,
	}
	quotes["USD"] = QuoInfo{
		MarketCap: tokenInfo.MarketData.MarketCap.USD,
		Volume24h: tokenInfo.MarketData.Volume24H.USD,
	}
	return TokenGeneralInfo{
		Quotes: quotes,
	}
}

type RateUSDCG struct {
	MarketData struct {
		CurrentPrice struct {
			USD float64 `json:"usd"`
		} `json:"current_price"`
	} `json:"market_data"`
}

type RateUSD struct {
	Symbol   string `json:"symbol"`
	PriceUsd string `json:"price_usd"`
}

type ResultRpc struct {
	Result string `json:"result"`
}

// type TokenInfoData struct {
// 	Data TokenGeneralInfo `data`
// }

type CandleTicker struct {
	SellPrice string `json:"sell_price"`
	BuyPrice  string `json:"buy_price"`
	UnixTime  uint64 `json:"unix_time"`
}

type TokenInfo struct {
	TokenSymbol       string                  `json:"symbol"`
	CirculatingSupply string                  `json:"circulating_supply"`
	TotalSupply       string                  `json:"total_supply"`
	MaxSupply         string                  `json:"max_supply"`
	MarketCap         string                  `json:"market_cap"`
	SellPrice         string                  `json:"sell_price"`
	BuyPrice          string                  `json:"buy_price"`
	Last7days         map[uint64]CandleTicker `json:"last_7d"`
}

type RateHistory struct {
	SellPrice string `json:"sell_price"`
	BuyPrice  string `json:"buy_price"`
}
type RateInfo struct {
	LastSell      string                 `json:"last_sell"`
	LastBuy       string                 `json:"last_buy"`
	HistoryRecord map[int64]*RateHistory `json:"history_record"`
}

type Rates struct {
	R float64   `json:"r"`
	P []float64 `json:"p"`
}

type MarketInfo struct {
	Rates  *Rates             `json:"rates"`
	Quotes map[string]QuoInfo `json:"quotes"`
}

type RightMarketInfo struct {
	Rate   *float64           `json:"rate"`
	Quotes map[string]QuoInfo `json:"quotes"`
}

func NewMarketInfo(quotes map[string]QuoInfo, rates *Rates) *MarketInfo {
	return &MarketInfo{
		Rates:  rates,
		Quotes: quotes,
	}
}

type TokenConfig struct {
	Success bool    `json:"success"`
	Data    []Token `json:"data"`
}
