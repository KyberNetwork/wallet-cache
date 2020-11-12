package gas

import "github.com/KyberNetwork/cache/ethereum"

type GasFetcher interface {
	GetGasPrice() (*ethereum.GasPrice, error)
}