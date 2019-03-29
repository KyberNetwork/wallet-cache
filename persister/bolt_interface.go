package persister

import "github.com/KyberNetwork/server-go/ethereum"

type BoltInterface interface {
	StoreGeneralInfo(map[string]*ethereum.TokenGeneralInfo) error
	GetGeneralInfo(map[string]ethereum.Token) (map[string]*ethereum.TokenGeneralInfo, error)
}
