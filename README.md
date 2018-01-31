# Cached server for KYBER wallet

## Build

```
docker build . -t cached
```

## Run
```
docker run -p 3001:3001 cached
```

## One command to build and run with docker-compose
```
docker-compose up
```

## Api:
 - /getRate
 - /getHistoryOneColumn
 - /getLatestBlock
 - /getRateUSD
