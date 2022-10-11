package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func genericRequest(client *http.Client, url string, values map[string]interface{}, bearer string) map[string]interface{} {
	jsonData, err := json.Marshal(values)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization", bearer)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	return res
}
