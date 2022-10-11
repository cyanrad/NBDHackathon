package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var accounts = make(map[string]account)

type account struct {
	ID       int32
	EM_ID    string
	Password string
	VDirhams float64
}

func createBankAccount(client *http.Client, bearer string) int32 {
	url := "https://api.futurehack.virtusa.dev/smartbankaccounts/1.0/accounts"
	values := map[string]interface{}{
		"accountCurrency":       "AED",
		"accountName":           "Reduwun",
		"nickname":              "Reduwun",
		"accountOpeningDate":    "2022-09-29",
		"accountClosingDate":    "2023-09-29",
		"accountTypeId":         1000000011,
		"productId":             1000000012,
		"accountrefnumber":      "64491989194984616561687",
		"balance":               10000,
		"bankId":                1006,
		"branchId":              6000000,
		"cardFacility":          "Y",
		"checkerDate":           "2022-09-29",
		"checkerId":             "1",
		"chequebookFacility":    "Y",
		"creditDebitIndicator":  "Debit",
		"creditLineAmount":      100,
		"creditLineIncluded":    "Y",
		"creditLineType":        "Pre-Agreed",
		"frozen":                "N",
		"isjointaccount":        "N",
		"isonlineaccessenabled": "N",
		"makerDate":             "2022-09-29",
		"makerId":               "1",
		"modifiedDate":          "2022-09-29",
		"noCredit":              "N",
		"noDebit":               "N",
		"nomineeAddress":        "UAE",
		"nomineeDob":            "2000-09-29",
		"nomineeName":           "John",
		"nomineePhoneNo":        "00000000",
		"nomineeRelatonship":    "Brother",
		"schemeName":            "BBAN",
		"status":                "Active",
		"stealthBookFacility":   "Y",
		"typeOfBalance":         "Information",
		"usage":                 "Y",
	}
	return (int32)(genericRequest(client, url, values, bearer)["accountId"].(float64))
}

func doesAccountExist(ID string) bool {
	_, ok := accounts[ID]
	return ok
}

func createAccount(emid string, password string, hasBankAccount bool) bool {
	if !doesAccountExist(emid) {
		var newID int32 = -1
		if hasBankAccount {
			token := getBearerToken(http.DefaultClient)
			newID = createBankAccount(http.DefaultClient, token)
		}

		// >> supposed to emulate a database
		newAcc := account{ID: newID, EM_ID: emid, Password: password}
		accounts[emid] = newAcc
		exportAccountAsJson()
		return true
	}
	return false
}

func exportAccountAsJson() {
	// clearing the file
	if err := os.Truncate("accounts.json", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	// writing accounts to file
	out, err := json.MarshalIndent(accounts, "", " ")
	if err != nil {
		panic(err)
	}

	if err = ioutil.WriteFile("accounts.json", out, 0644); err != nil {
		panic(err)
	}
}

func loadAccountsFromJson() {
	jsonFile, err := os.Open("accounts.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &accounts)
}
