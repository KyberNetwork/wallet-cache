package fCommon

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

func HTTPCall(url string) ([]byte, error) {
	response, err := http.Get(url)
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
