# Cached server for KYBER wallet

## Build

```
docker build . -t cached
```

## Run
docker run -p 3001:3001 cached
```


## One command to build and run with docker-compose

docker-compose -f docker-compose-staging.yml up --build
```

## Api:
 - /currencies: return list supported tokens	 - /getRate
 - /latestBlock: return latest block number of network	 - /getHistoryOneColumn
 - /rateUSD: return USD price of token base on it's expectedRate	 - /getLatestBlock
 - /rate: return rate of token with eth (expectedRate and minRate)	 - /getRateUSD
 - /kyberEnabled: get kyberEnabled from contract	 - /getKyberEnabled
 - /maxGasPrice: get max GasPrice from contract	 - /getMaxGasPrice
 - /gasPrice: return gasPrice get from https://ethgasstation.info/	 - /getGasPrice
 - /marketInfo: return market info (volume, marketcap, ...) from CMC	 - /getMarketInfo
 - /last7D: ```params: listToken=KNC-DAI-...``` return last 7 days mid price (base on ETH) of token in listToken	
  param listToken is created by linking tokens (token's symbol in uppercase) with "-"	
 - /rateETH: return USD price of ETH from CMC	
 - /cacheVersion: return current cached version	
