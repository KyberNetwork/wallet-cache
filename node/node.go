package node

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type NodeMiddleware struct {
	client    *http.Client
	nodeCache *NodeCache
}

var whiteListArr = []string{"kyberswap.com", "knstats.com"}

// var banMethod = []string{"eth_getBlockByNumber"}
var banMethod = []string{}

func NewNodeMiddleware() (*NodeMiddleware, error) {
	return &NodeMiddleware{
		client:    &http.Client{},
		nodeCache: NewNodeCache(),
	}, nil
}

func (n *NodeMiddleware) HandleNodeRequest(c *gin.Context) {
	req := c.Request

	if err := filterRequest(req); err != nil {
		log.Print(err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"err": err.Error()},
		)
		return
	}

	respBytes, err := n.nodeCache.HandleRequest(req)
	if err != nil {
		log.Print(err)
		c.JSON(
			http.StatusBadGateway,
			gin.H{"err": err.Error()},
		)
		return
	}

	c.Writer.Write(respBytes)
}

func filterRequest(req *http.Request) error {

	kyberENV := os.Getenv("KYBER_ENV")
	if kyberENV != "production" {
		return nil
	}

	// check origin
	origin := req.Header.Get("Origin")

	if !InListSubstring(origin, whiteListArr) {
		return errors.New("Domain is not in the whitelist: " + origin)
	}

	// check method
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print(err)
		return err
	}

	var requestRpc RequestRPC
	err = json.Unmarshal(body, &requestRpc)
	if err != nil {
		log.Print(err)
		return err
	}

	if InList(requestRpc.Method, banMethod) {
		return errors.New("Method is not allowed")
	}

	// reassign again
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	return nil
}
