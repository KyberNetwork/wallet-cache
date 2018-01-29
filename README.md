Cached server for KYBER wallet

Build

docker build . -t kyber/server

Run

docker run -p 3002:3002 kyber/server


Api:
 - /getRate
 - /getHistoryOneColumn
 - /getLatestBlock
 - /getRateUSD
