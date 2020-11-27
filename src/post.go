package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

func postJson(url string, apiKey string, json []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	req.Header.Set("user-agent", "mendsail-cli/1.0")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, bodyErr := ioutil.ReadAll(resp.Body)

	if resp.Status[0] != '2' {
		response := ""
		if bodyErr == nil {
			response = "\n" + string(body)
		}
		return errors.New("Server returned error: " + resp.Status + response)
	}

	return nil
}
