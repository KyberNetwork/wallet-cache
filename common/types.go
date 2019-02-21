package common

import (
	"math/big"
)

type UserInfo struct {
	Cap   *big.Int `json:"cap"`
	Kyced bool     `json:"kyced"`
	Rich  bool     `json:"rich"`
}
