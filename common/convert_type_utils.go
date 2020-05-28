package common

import (
	"github.com/KyberNetwork/cache/ethereum"
)

// ArrTokenToMap convert array token to map with key is address
func ArrTokenToMap(listToken []ethereum.Token) map[string]ethereum.Token {
	m := make(map[string]ethereum.Token)
	for _, t := range listToken {
		m[t.TokenID] = t
	}
	return m
}
