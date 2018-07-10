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
 - /getRate
 - /getHistoryOneColumn
 - /getLatestBlock
 - /getRateUSD
 - /getKyberEnabled
 - /getMaxGasPrice
 - /getGasPrice
 - /getMarketInfo