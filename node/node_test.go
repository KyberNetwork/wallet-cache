package node

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {

	bodyStr := `{"jsonrpc":"2.0","id":1018,"method":"eth_getTransactionCount","params":["0x3cf628d49ae46b49b210f0521fbd9f82b461a9e1","latest"]}`
	rbody := bytes.NewReader([]byte(bodyStr))

	request, err := http.NewRequest("POST", "https://eth-mainnet.alchemyapi.io/jsonrpc/hYvhNnvmAnUIkI7z1r2yTbvXPbMftTli", rbody)
	if err != nil {
		log.Print(err)
		return
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(request)
	if err != nil {
		log.Print(err)
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return
	}

	log.Print(string(bodyBytes))

}
