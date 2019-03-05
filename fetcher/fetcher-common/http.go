package fCommon

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func HTTPCall(url string) ([]byte, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Get(url)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if response.StatusCode != 200 {
		return []byte{}, errors.New("Status code is not 200")
	}

	defer (response.Body).Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return b, nil
}
