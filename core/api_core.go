package core

import (
	"github.com/KyberNetwork/server-go/fetcher"
	"github.com/KyberNetwork/server-go/persister"
	"github.com/KyberNetwork/server-go/common"
	"github.com/KyberNetwork/server-go/ethereum"
	"math/big"
	"log"
	"errors"
)

// Core is struct data include fetcher and storage
type Core struct {
	fetcher        *fetcher.Fetcher
	storage        persister.Persister	
}

// NewCore init core instance
func NewCore( fertcherIns *fetcher.Fetcher, persisterIns persister.Persister) (*Core, error) {
	return &Core{
		fetcher:        fertcherIns,
		storage:        persisterIns,		
	}, nil
}

func (self *Core) amountFromStepAmount(src string, dest string, destAmount string, stepRate []ethereum.StepRate)(string, error){
	// find 2 pair
	lowerDestAmount := big.NewInt(0)	
	lowerSourceAmount := big.NewInt(0)	
	higherDestAmount := big.NewInt(0)
	higherSourceAmount := big.NewInt(0)		

	zeroAmount := big.NewInt(0)
	amountBig := new(big.Int)
	amountBig, ok := amountBig.SetString(destAmount, 10)
	if !ok {
		return "", errors.New("cannot convert destAmount to big Int")
	}
	for _, rate := range stepRate {
		if rate.Dest == dest && rate.Source == src {			
			if rate.DestAmount.Cmp(amountBig) == 1 {
				if higherDestAmount.Cmp(zeroAmount) == 0 {
					higherDestAmount = rate.DestAmount
					higherSourceAmount = rate.SrcAmount
				}
				if higherDestAmount.Cmp(zeroAmount) != 0 && higherDestAmount.Cmp(rate.DestAmount) == 1{
					higherDestAmount = rate.DestAmount
					higherSourceAmount = rate.SrcAmount
				}
			}

			if rate.DestAmount.Cmp(amountBig) != 1 {
				if lowerDestAmount.Cmp(zeroAmount) == 0 {
					lowerDestAmount = rate.DestAmount
					lowerSourceAmount = rate.SrcAmount
				}
				if lowerDestAmount.Cmp(zeroAmount) != 0 && lowerDestAmount.Cmp(rate.DestAmount) != 1{
					lowerDestAmount = rate.DestAmount
					lowerSourceAmount = rate.SrcAmount
				}
			}
		}
	}
	// log.Println(higherDestAmount.String())
	// log.Println(higherSourceAmount.String())
	// log.Println(lowerDestAmount.String())
	// log.Println(lowerSourceAmount.String())
	// log.Println(amountBig.String())

	if lowerDestAmount.Cmp(zeroAmount) == 0 && higherDestAmount.Cmp(zeroAmount) == 0{ 
		// srcAmount := big.NewInt(0).Mul(higherSourceAmount, amountBig).Quo(higherDestAmount)
		err := errors.New("Canot calculate source amount")
		log.Println(err)
		return "", err
	}

	if higherDestAmount.Cmp(zeroAmount) == 0 && lowerDestAmount.Cmp(zeroAmount) != 0{ 
		srcAmount := big.NewInt(0).Mul(lowerSourceAmount, amountBig)
		srcAmount = big.NewInt(0).Quo(srcAmount, lowerDestAmount)
		return srcAmount.String(), nil
	}

	if lowerDestAmount.Cmp(zeroAmount) == 0 && higherDestAmount.Cmp(zeroAmount) != 0{ 
		srcAmount :=  big.NewInt(0).Mul(higherSourceAmount, amountBig)
		srcAmount = big.NewInt(0).Quo(srcAmount, higherDestAmount)
		return srcAmount.String(), nil
	}

	if (lowerDestAmount.Cmp(higherDestAmount) == 0){
		return lowerSourceAmount.String(), nil
	}


	factor1 := big.NewInt(0).Mul(higherSourceAmount, big.NewInt(0).Sub(amountBig, lowerDestAmount))
	factor2 := big.NewInt(0).Mul(lowerSourceAmount, big.NewInt(0).Sub(higherDestAmount, amountBig))
	factor3 := big.NewInt(0).Sub(higherDestAmount, lowerDestAmount)

	srcAmount := big.NewInt(0).Add(factor1, factor2)
	srcAmount = big.NewInt(0).Quo(srcAmount, factor3)
	return srcAmount.String(), nil
}

func (self *Core) amountFromRateInit(src string, dest string, destAmount string )(string, error){
	rates := self.storage.GetRate()
	destToken, err := self.fetcher.GetTokenBySymbol(dest)
	if err != nil {
		log.Println(err)
		return "", err
	}
	srcToken, err := self.fetcher.GetTokenBySymbol(src)
	if err != nil {
		log.Println(err)
		return "", err
	}
	for _, rate := range rates {
		if rate.Source == src && rate.Dest == dest {
			srcAmount, err := common.FromDestToSrc(destAmount, rate.Rate, srcToken.Decimal, destToken.Decimal)
			if err != nil {
				log.Println(err)
				return "", err
			}
			return srcAmount.String(), nil
		}
	}
		return "", errors.New("Cannot get src amount from rate init")
}

func (self *Core) fromTokenToEth(dest string, destAmount string)(string, error){
	// check step rate
	stepRate := self.storage.GetStepRate()
	check := false
	for _, rate := range stepRate{
		if rate.Dest == dest{
			check = true
			break
		}
	}
	if check {
		// calcualte rate from step rate
		srcAmount, err := self.amountFromStepAmount("ETH", dest, destAmount, stepRate)
		if err != nil {
			log.Println(err)
			return "", err
		}
		return srcAmount, nil
	}else{
		//calcualte rate from rate init		
		srcAmount, err := self.amountFromRateInit("ETH", dest, destAmount)
		if err != nil {
			log.Println(err)
			return "", err
		}
		return srcAmount, nil
	}
}

func (self *Core) fromEthToToken(src string, destAmount string)(string, error){
	stepRate := self.storage.GetStepRate()
	check := false
	for _, rate := range stepRate{
		if rate.Dest == src{
			check = true
			break
		}
	}

	if check {
		// calcualte rate from step rate
		srcAmount, err := self.amountFromStepAmount(src, "ETH", destAmount, stepRate)
		if err != nil {
			log.Println(err)
			return "", err
		}
		return srcAmount, nil
	}else{
		//calcualte rate from rate init		
		srcAmount, err := self.amountFromRateInit(src, "ETH", destAmount)
		if err != nil {
			log.Println(err)
			return "", err
		}
		return srcAmount, nil
	}

}

func (self *Core) GetSourceAmount(src string, dest string, destAmount string)( string, error){
	if src == dest {
		return destAmount, nil
	}
	// log.Println(src)
	// log.Println(dest)
	// log.Println(destAmount)
	_, err := self.fetcher.GetTokenBySymbol(src)
	if err != nil {
		log.Println(err)
		return "", err
	}
	_, err = self.fetcher.GetTokenBySymbol(dest)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// if dest != ETH get rate from token to ETH, calculate amount ETH
	srcEth := destAmount
	if dest != "ETH"{
		srcEth, err = self.fromTokenToEth(dest, destAmount)
		if err != nil {
			log.Println(err)
			return "", err
		}
	}

	srcAmount := srcEth
	if src != "ETH" {
		srcAmount, err = self.fromEthToToken(src, srcEth)
		if err != nil {
			log.Println(err)
			return "", err
		}
	}

	return srcAmount, nil
	
}

