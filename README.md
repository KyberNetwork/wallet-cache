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
 - /getRightMarketInfo: return market info (volume, marketcap, ...) from Coingecko
 - /getLast7D: ```params: listToken=KNC-DAI-...``` return last 7 days mid price (base on ETH) of token in listToken
  param listToken is created by linking tokens (token's symbol in uppercase) with "-"
 - /getRateETH: return USD price of ETH from Coingecko

## New APIs

 - /latestBlock: return latest block number of network
 - /rateUSD: return USD price of token base on it's expectedRate
 - /rate: return rate of token with eth (expectedRate and minRate)
 - /kyberEnabled: get kyberEnabled from contract
 - /maxGasPrice: get max GasPrice from contract
 - /gasPrice: return gasPrice get from https://ethgasstation.info/
 - /marketInfo: return market info (volume, marketcap, ...) from Coingecko
 - /last7D: ```params: listToken=KNC-DAI-...``` return last 7 days mid price (base on ETH) of token in listToken
  param listToken is created by linking tokens (token's symbol in uppercase) with "-"
 - /rateETH: return USD price of ETH from Coingecko

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
            "source": "0xa94758d328af7ef1815e73053e95b5f86588c16d",
            "dest": "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
            "rate": "0",
            "minRate": "0"
        },
        {
            "source": "0x7a17267576318efb728bc4a0833e489a46ba138f",
            "dest": "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
            "rate": "0",
            "minRate": "0"
        },
        {
            "source": "0x72fd6c7c1397040a66f33c2ecc83a0f71ee46d5c",
            "dest": "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
            "rate": "178480000000000",
            "minRate": "0"
        },
        {
            "source": "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
            "dest": "0x95cc8d8f29d0f7fcc425e8708893e759d1599c97",
            "rate": "287229960000000000000",
            "minRate": "278613061200000000000"
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

### 7. Get marketInfo
`/marketInfo`

(GET) Return market info (volume, marketcap, ...) from Coingecko

Response:
```javascript
{
    "data": {
        "0x1742c81075031b8f173d2327e3479d1fc3feaa76": {
            "rate": null,
            "change_24h": "3.71950684029737",
            "quotes": {
                "ETH": {
                    "market_cap": 52535.45868371555,
                    "volume_24h": 1715.0882823743132
                },
                "USD": {
                    "market_cap": 7904239.161183325,
                    "volume_24h": 258066.9832533837
                }
            }
        },
        "0x2799f05b55d56be756ca01af40bf7350787f48d4": {
            "rate": null,
            "change_24h": "2.7791374304334",
            "quotes": {
                "ETH": {
                    "market_cap": 32373.764520744804,
                    "volume_24h": 646.0999704525954
                },
                "USD": {
                    "market_cap": 4870803.400472581,
                    "volume_24h": 97209.1438766459
                }
            }
        },
        "0x499990db50b34687cdafb2c8dabae4e99d6f38a7": {
            "rate": null,
            "change_24h": "1.08168964231987",
            "quotes": {
                "ETH": {
                    "market_cap": 55691.72455294219,
                    "volume_24h": 2844.0081889865687
                },
                "USD": {
                    "market_cap": 8379113.314326793,
                    "volume_24h": 427896.01280416525
                }
            }
        }
    },
    "status": "latest",
    "success": true
}
```

### 8. Get last7D
`/last7D`

(GET) Return last 7 days mid price (base on ETH) of token in listToken param

Input Request Parameters
 
|Name | Type | Required | Description |
| ----------| ---------|------|-----------------------------|
|listToken|STRING|YES|The list token addresses, split by `-`|

ex: `/last7D?listToken=0x1742c81075031b8f173d2327e3479d1fc3feaa76-0x2799f05b55d56be756ca01af40bf7350787f48d4`

Response:
```javascript
{
    "data": {
        "0xad6d458402f60fd3bd25163575031acdce07538d": [
            0.006902771878817366,
            0.006618834209810034,
            0.006567853608262672,
            0.006509749112258536,
            0.006590922019290737,
            0.0067025246634046355,
            0.00676529933031774,
            0.006688278062621764,
            0.006569355373888682,
            0.0066587748879877555,
            0.0065973733552006425,
            0.006300725041627652,
            0.006440412028755606,
            0.006333154671241967,
            0.006344202073991879,
            0.0065451727836549995,
            0.006551426316158762,
            0.006474764775608181,
            0.006380682075884695,
            0.006455922292504784,
            0.006464179096982978,
            0.006477914623047843,
            0.006554514663776798,
            0.006668032217882891,
            0.006721637489388994,
            0.006583584985360085,
            0.006609396558562976,
            0.006564737023212862,
            0.006529161894668466
        ],
        "0x4e470dc7321e84ca96fcaedd0c8abcebbaeb68c6": [
            0.001106960802998027,
            0.0010872301461882037,
            0.001084494293012309,
            0.0010772279957433097,
            0.001068345398756188,
            0.0010697819251185591,
            0.0010541801424552954,
            0.0010418712420932937,
            0.001016859122973702,
            0.0010376282199194462,
            0.0010407530293247377,
            0.0009992837031336048,
            0.0010081489245161934,
            0.0010119447037499193,
            0.001015479738760005,
            0.001042729468247283,
            0.0010822533280792272,
            0.0011056124461407597,
            0.0010989667276303384,
            0.0010991963272327303,
            0.0010869376839133889,
            0.0010874043484461772,
            0.0010798410058743648,
            0.0010865361589820774,
            0.0010742438649890488,
            0.0010784629605376423,
            0.001082174642938953,
            0.0010816088282831681,
            0.001078535332811308
        ]
    },
    "status": "latest",
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

