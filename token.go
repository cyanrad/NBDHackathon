package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func getBearerToken(client *http.Client) string {
	TokenValues := map[string]interface{}{"accessToken": "", "userName": "rjabraouti@outlook.com", "password": ""}
	jsonData, err := json.Marshal(TokenValues)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://services.futurehack.virtusa.dev/api-sandbox/application/token", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	return "bearer " + res["access_token"].(string)
}
