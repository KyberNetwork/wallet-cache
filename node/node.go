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
	client *http.Client
	url    string
}

var whiteListArr = []string{"kyberswap.com", "knstats.com"}
var banMethod = []string{"eth_getBlockByNumber"}

func NewNodeMiddleware() (*NodeMiddleware, error) {
	nodeEnpoint := os.Getenv("NODE_ENDPOINT")
	return &NodeMiddleware{
		client: &http.Client{},
		url:    nodeEnpoint,
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

	proxyReq, err := n.cloneRequest(req)
	if err != nil {
		log.Print(err)
		c.JSON(
			http.StatusBadGateway,
			gin.H{"err": err.Error()},
		)
		return
	}

	// We may want to filter some headers, otherwise we could just use a shallow copy
	resp, err := n.client.Do(proxyReq)
	if err != nil {
		log.Print(err)
		c.JSON(
			http.StatusBadGateway,
			gin.H{"err": err.Error()},
		)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			c.JSON(
				http.StatusBadGateway,
				gin.H{"err": err.Error()},
			)
			return
		}
		c.Writer.Write(bodyBytes)
		return
	}
	c.JSON(
		resp.StatusCode,
		gin.H{"err": "Status code is not 200"},
	)

}

func (n *NodeMiddleware) cloneRequest(req *http.Request) (*http.Request, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	proxyReq, err := http.NewRequest(req.Method, n.url, bytes.NewReader(body))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// proxyReq.Header = make(http.Header)
	// for h, val := range req.Header {
	// 	proxyReq.Header[h] = val
	// }
	// log.Print(req.Header)
	proxyReq.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")

	return proxyReq, nil
}

func filterRequest(req *http.Request) error {

	kyberENV := os.Getenv("KYBER_ENV")
	if kyberENV != "production" {
		return nil
	}

	// check origin
	origin := req.Header.Get("Origin")

	if !InListSubstring(origin, whiteListArr) {
		return errors.New("Domain is not in the whitelist")
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
