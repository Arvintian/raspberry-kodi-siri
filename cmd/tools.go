package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func httppost(url string, data []byte) ([]byte, error) {
	rsp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(rsp.Body)
}
