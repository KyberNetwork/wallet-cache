package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/KyberNetwork/cache/fetcher"
	"github.com/KyberNetwork/cache/node"
	persister "github.com/KyberNetwork/cache/persister"
	"github.com/KyberNetwork/cache/refprice"
	raven "github.com/getsentry/raven-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/gin"
)

const (
	MAX_PAGE_SIZE = 50
	DEFAULT_PAGE  = 1
)

type HTTPServer struct {
	node      *node.NodeMiddleware
	fetcher   *fetcher.Fetcher
	persister persister.Persister
	host      string
	r         *gin.Engine
	refPrice  *refprice.RefPrice
}

func (self *HTTPServer) GetRate(c *gin.Context) {
	isNewRate := self.persister.GetIsNewRate()
	if isNewRate != true {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "data": nil},
		)
		return
	}

	rates := self.persister.GetRate()
	updateAt := self.persister.GetTimeUpdateRate()
	c.JSON(
		http.StatusOK,
		gin.H{"success": true, "updateAt": updateAt, "data": rates},
	)
}

func (self *HTTPServer) GetLatestBlock(c *gin.Context) {
	if !self.persister.GetIsNewLatestBlock() {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false},
		)
		return
	}
	blockNum := self.persister.GetLatestBlock()
	c.JSON(
		http.StatusOK,
		gin.H{"success": true, "data": blockNum},
	)
}

func (self *HTTPServer) GetRateUSD(c *gin.Context) {
	if !self.persister.GetIsNewRateUSD() {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false},
		)
		return
	}

	rates := self.persister.GetRateUSD()
	c.JSON(
		http.StatusOK,
		gin.H{"success": true, "data": rates},
	)
}

func (self *HTTPServer) GetRateETH(c *gin.Context) {
	if !self.persister.GetIsNewRateUSD() {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false},
		)
		return
	}

	ethRate := self.persister.GetRateETH()
	c.JSON(
		http.StatusOK,
		gin.H{"success": true, "data": ethRate},
	)
}

func (self *HTTPServer) GetKyberEnabled(c *gin.Context) {
	if !self.persister.GetNewKyberEnabled() {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false},
		)
		return
	}

	enabled := self.persister.GetKyberEnabled()
	c.JSON(
		http.StatusOK,
		gin.H{"success": true, "data": enabled},
	)
}

func (self *HTTPServer) GetMaxGasPrice(c *gin.Context) {
	if !self.persister.GetNewMaxGasPrice() {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false},
		)
		return
	}

	gasPrice := self.persister.GetMaxGasPrice()
	c.JSON(
		http.StatusOK,
		gin.H{"success": true, "data": gasPrice},
	)
}

func (self *HTTPServer) GetGasPrice(c *gin.Context) {
	if !self.persister.GetNewGasPrice() {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false},
		)
		return
	}

	gasPrice := self.persister.GetGasPrice()
	c.JSON(
		http.StatusOK,
		gin.H{"success": true, "data": gasPrice},
	)
}

func (self *HTTPServer) GetErrorLog(c *gin.Context) {
	dat, err := ioutil.ReadFile("error.log")
	if err != nil {
		log.Print(err)
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "data": err},
		)
	}
	c.JSON(
		http.StatusOK,
		gin.H{"success": true, "data": string(dat[:])},
	)
}

func (self *HTTPServer) getCacheVersion(c *gin.Context) {
	timeRun := self.persister.GetTimeVersion()
	c.JSON(
		http.StatusOK,
		gin.H{"success": true, "data": timeRun},
	)
}

func (self *HTTPServer) GetUserInfo(c *gin.Context) {
	address := c.Query("address")
	userInfo, err := self.fetcher.FetchUserInfo(address)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		userInfo,
	)
}

func (self *HTTPServer) GetSourceAmount(c *gin.Context) {
	src := c.Query("source")
	dest := c.Query("dest")
	destAmount := c.Query("destAmount")
	srcAmount, err := self.fetcher.GetSourceAmount(src, dest, destAmount)

	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"success": true, "value": srcAmount},
	)
}

func (self *HTTPServer) GetRefprice(c *gin.Context) {
	base := c.Query("base")
	quote := c.Query("quote")

	price, err := self.refPrice.GetRefPrice(strings.ToUpper(base), strings.ToUpper(quote))
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{"success": true, "value": price},
	)
}

func (self *HTTPServer) PostNodeRequest(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, GET, OPTIONS, PATCH, POST, PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "accept, accept-encoding, authorization, content-type, dnt, origin, user-agent, x-csrftoken, x-requested-with, alchemy-web3-version")

	self.node.HandleNodeRequest(c)
}

func (self *HTTPServer) Run(kyberENV string) {
	self.r.GET("/getLatestBlock", self.GetLatestBlock)
	self.r.GET("/latestBlock", self.GetLatestBlock)

	self.r.GET("/getRateUSD", self.GetRateUSD)
	self.r.GET("/rateUSD", self.GetRateUSD)

	self.r.GET("/getRate", self.GetRate)
	self.r.GET("/rate", self.GetRate)

	self.r.GET("/getKyberEnabled", self.GetKyberEnabled)
	self.r.GET("/kyberEnabled", self.GetKyberEnabled)

	self.r.GET("/getMaxGasPrice", self.GetMaxGasPrice)
	self.r.GET("/maxGasPrice", self.GetMaxGasPrice)

	self.r.GET("/getGasPrice", self.GetGasPrice)
	self.r.GET("/gasPrice", self.GetGasPrice)

	self.r.GET("/getRateETH", self.GetRateETH)
	self.r.GET("/rateETH", self.GetRateETH)

	self.r.GET("/cacheVersion", self.getCacheVersion)

	self.r.GET("/users", self.GetUserInfo)

	self.r.GET("/sourceAmount", self.GetSourceAmount)

	self.r.GET("/refprice", self.GetRefprice)

	self.r.POST("/node", self.PostNodeRequest)

	// if kyberENV != "production" {
	// 	self.r.GET("/9d74529bc6c25401a2f984ccc9b0b2b3", self.GetErrorLog)
	// }

	self.r.Run(self.host)
}

func NewHTTPServer(host string, persister persister.Persister, fetcher *fetcher.Fetcher, node *node.NodeMiddleware) *HTTPServer {
	r := gin.Default()
	r.Use(sentry.Recovery(raven.DefaultClient, false))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"DELETE", "GET", "OPTIONS", "PATCH", "POST", "PUT"}
	corsConfig.AllowHeaders = []string{"accept", "accept-encoding", "authorization", "content-type", "dnt", "origin", "user-agent", "x-csrftoken", "x-requested-with", "alchemy-web3-version"}
	corsConfig.AllowCredentials = true

	corsConfig.MaxAge = 5 * time.Minute

	r.Use(cors.New(corsConfig))

	refPrice := refprice.NewRefPrice()

	return &HTTPServer{
		node, fetcher, persister, host, r, refPrice,
	}
}
