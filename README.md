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

## APIs (these APIs will be expired after Jan 20 2019)
 - /getLatestBlock: return latest block number of network
 - /getRateUSD: return USD price of token base on it's expectedRate
 - /getRate: return rate of token with eth (expectedRate and minRate)
 - /getKyberEnabled: get kyberEnabled from contract
 - /getMaxGasPrice: get max GasPrice from contract
 - /getGasPrice: return gasPrice get from https://ethgasstation.info/
 - /getRateETH: return USD price of ETH from Coingecko

## New APIs

 - /latestBlock: return latest block number of network
 - /rateUSD: return USD price of token base on it's expectedRate
 - /rate: return rate of token with eth (expectedRate and minRate)
 - /kyberEnabled: get kyberEnabled from contract
 - /maxGasPrice: get max GasPrice from contract
 - /gasPrice: return gasPrice get from https://ethgasstation.info/
 - /rateETH: return USD price of ETH from Coingecko
 - /users: ```params: address=0x2262d4f6312805851e3b27c40db2c7282e6e4a42``` return user stats info
 - /sourceAmount: ```params: ?source=TUSD&dest=ETH&destAmount=500``` calculate and return relative src amount when having dest amount
 
## Cache version
 - /cacheVersion: return current cache version
 
 ### 1. Get Latest Block
`/latestBlock`

(GET) Return latest block number of network

Response:
```javascript
{
    "data": "4790885",
    "success": true
}
```

### 2. Get Rate USD
`/rateUSD`

(GET) Return USD price of token base on it's expectedRate

Response:
```javascript
{
    "data": [
        {
            "symbol": "ETH",
            "price_usd": "150.110255"
        }
    ],
    "success": true
}
```

### 3. Get rate
`/rate`

(GET) Return rate of token with eth (expectedRate and minRate)

Response:
```javascript
{
    "data": [
        {
            "source": "POWR",
            "dest": "ETH",
            "rate": "580350000000000",
            "minRate": "562939500000000"
        },
        {
            "source": "REQ",
            "dest": "ETH",
            "rate": "251549999999999",
            "minRate": "244003499999999"
        }
    ],
    "success": true
    }
```

### 4. Get kyberEnabled
`/kyberEnabled`

(GET) Return kyberEnabled from contract

Response:
```javascript
{
    "data": true,
    "success": true
}
```

### 5. Get maxGasPrice
`/maxGasPrice`

(GET) Return max GasPrice from contract

Response:
```javascript
{
    "data": "50000000000",
    "success": true
}
```

### 6. Get gasPrice
`/gasPrice`

(GET) Return gasPrice get from https://ethgasstation.info/

Response:
```javascript
{
    "data": {
        "fast": "10",
        "standard": "5.55",
        "low": "1.1",
        "default": "5.55"
    },
    "success": true
}
```

### 9. Get gasPrice
`/rateETH`

(GET) Return USD price of ETH from Coingecko

Response:
```javascript
{
    "data": "150.480634",
    "success": true
}
```

### 10. Get cacheVersion
`/cacheVersion`

(GET) Return current cache version

Response:
```javascript
{
    "data": "14:40:42 09-01-2019",
    "success": true
}
```

### 11. Get UserInfo
`/users?address=0x2262d4f6312805851e3b27c40db2c7282e6e4a42`

(GET) Return User stats info

Response:
```javascript
{
    "cap": 40304044000000000000,
    "kyced": true,
    "rich": false
}
```

### 12. Get source amount from dest amount
`/sourceAmount?source=TUSD&dest=ETH&destAmount=500`

(GET) Return Source Amount

Response:
```javascript
{
  "success": true,
  "value": "129808.7692"
}
```
