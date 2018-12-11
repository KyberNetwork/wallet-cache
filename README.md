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
 - /getLatestBlock: return latest block number of network
 - /getRateUSD: return USD price of token base on it's expectedRate
 - /getRate: return rate of token with eth (expectedRate and minRate)
 - /getKyberEnabled: get kyberEnabled from contract
 - /getMaxGasPrice: get max GasPrice from contract
 - /getGasPrice: return gasPrice get from https://ethgasstation.info/
 - /getRightMarketInfo: return market info (volume, marketcap, ...) from Coingecko
 - /getLast7D: ```params: listToken=KNC-DAI-...``` return last 7 days mid price (base on ETH) of token in listToken
  param listToken is created by linking tokens (token's symbol in uppercase) with "-"
 - /getRateETH: return USD price of ETH from Coingecko
 - /cacheVersion: return current cached version
