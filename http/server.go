package http

import (
	"net/http"

	"github.com/KyberNetwork/server-go/persistor"
	raven "github.com/getsentry/raven-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	persistor persistor.Persistor
	host      string
	r         *gin.Engine
}

func (self *HTTPServer) GetRate(c *gin.Context) {
	if !self.persistor.GetIsNewRate() {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false},
		)
		return
	}

	rates := self.persistor.GetRate()
	c.JSON(
		http.StatusOK,
		gin.H{"success": true, "data": rates},
	)
	return
}

func (self *HTTPServer) GetEvent(c *gin.Context) {
	if !self.persistor.GetIsNewEvent() {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false},
		)
		return
	}

	events := self.persistor.GetEvent()
	c.JSON(
		http.StatusOK,
		gin.H{"success": true, "data": events},
	)
}

func (self *HTTPServer) GetLatestBlock(c *gin.Context) {
	if !self.persistor.GetIsNewLatestBlock() {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false},
		)
		return
	}
	blockNum := self.persistor.GetLatestBlock()
	c.JSON(
		http.StatusOK,
		gin.H{"success": true, "data": blockNum},
	)
}

func (self *HTTPServer) GetRateUSD(c *gin.Context) {
	if !self.persistor.GetIsNewRateUSD() {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false},
		)
		return
	}

	rates := self.persistor.GetRateUSD()
	c.JSON(
		http.StatusOK,
		gin.H{"success": true, "data": rates},
	)
}

// func (self *HTTPServer) GetLanguagePack(c *gin.Context) {
// 	c.JSON(
// 		http.StatusOK,
// 		gin.H{"success": true, "data": "get language pack"},
// 	)
// 	return
// }

func (self *HTTPServer) Run() {
	self.r.GET("/getRate", self.GetRate)
	self.r.GET("/getHistoryOneColumn", self.GetEvent)
	self.r.GET("/getLatestBlock", self.GetLatestBlock)

	self.r.GET("/getRateUSD", self.GetRateUSD)
	//self.r.GET("/getLanguagePack", self.GetLanguagePack)

	self.r.Run(self.host)
}

func NewHTTPServer(host string, persistor persistor.Persistor) *HTTPServer {
	r := gin.Default()
	r.Use(sentry.Recovery(raven.DefaultClient, false))
	r.Use(cors.Default())

	return &HTTPServer{
		persistor, host, r,
	}
}
