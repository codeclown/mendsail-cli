package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func postJson(url string, json []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("X-Custom-Header", "myvalue")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return nil
}
