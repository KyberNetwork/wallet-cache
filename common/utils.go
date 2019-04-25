package common

import "github.com/KyberNetwork/server-go/ethereum"

// IsDifferentMapToken compare two map tokens
func IsDifferentMapToken(mapTokenA, mapTokenB map[string]ethereum.Token) bool {
	if len(mapTokenA) != len(mapTokenB) {
		return true
	}
	for t := range mapTokenA {
		if _, ok := mapTokenB[t]; !ok {
			return true
		}
	}
	return false
}

// CopyMapToken make a copy of
func CopyMapToken(mapTokenDest, mapTokenSource map[string]ethereum.Token) {
	for k, v := range mapTokenSource {
		mapTokenDest[k] = v
	}
}
