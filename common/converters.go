package common 

import (
	"math/big"
	"errors"
	"log"
	"strconv"
)

func GetOrAmount(amount *big.Int) *big.Int {
	orNumber := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(255), nil)
	orAmount := big.NewInt(0).Or(amount, orNumber)
	return orAmount
}

func StringToWei(amount string, decimals int) (*big.Int, error) {
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ToWei(amountFloat, decimals), nil
}

func ToWei(amount float64, decimals int) *big.Int {
	amountFloat := big.NewFloat(amount)
	weight := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	weightFloat := new(big.Float).SetInt(weight)

	amountBig := big.NewFloat(0).Mul(amountFloat, weightFloat)

	result := new(big.Int)
	amountBig.Int(result) // store converted number in result
	
	return result
}

func ToToken(amount string, decimals int) (*big.Float, error) {
	amountBig := new(big.Int)
	amountBig, ok := amountBig.SetString(amount, 10)
	if !ok {
		return nil, errors.New("cannot convert rate to big Int")
	}
	amountBigFloat := new(big.Float).SetInt(amountBig)

	weight := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	weightBigFloat := new(big.Float).SetInt(weight)


	token := new(big.Float).Quo(amountBigFloat, weightBigFloat)
	return token, nil
}


func GetAmountEnableFirstBit(amount float64,  decimals int) *big.Int{
	amountBig := ToWei(amount, decimals)
	return GetOrAmount(amountBig)
}


func FromSrcToDest(srcAmount string, rate string, srcDecimals int, destDecimal int ) (*big.Int, error){
	srcBig, err := ToToken(srcAmount, srcDecimals)
	if err != nil {
		log.Println(err)
		return nil, err
	}	
	rateBig, err := ToToken(rate, 18)
	if err != nil {
		log.Println(err)
		return nil, err
	}	

	// calculate dest
	destAmount := new(big.Float).Mul(srcBig, rateBig)
	
	destAmountFloat, _ := destAmount.Float64()

	return ToWei(destAmountFloat, destDecimal), nil

}



func FromDestToSrc(destAmount string, rate string, srcDecimals int, destDecimal int ) (*big.Int, error){
	destBig, err := ToToken(destAmount, destDecimal)
	if err != nil {
		log.Println(err)
		return nil, err
	}	
	rateBig, err := ToToken(rate, 18)
	if err != nil {
		log.Println(err)
		return nil, err
	}	
	
	if rateBig.Cmp(big.NewFloat(0)) == 0 {
		err = errors.New("Cannot calculate src amount")
		log.Println(err)
		return nil, err
	}

	// calculate dest
	srcAmount := new(big.Float).Quo(destBig, rateBig)
	
	srcAmountFloat, _ := srcAmount.Float64()

	return ToWei(srcAmountFloat, srcDecimals), nil

}



