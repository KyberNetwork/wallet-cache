package ethereum

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
