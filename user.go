package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var currentlyLoggedIn string

func getBankAccount(client *http.Client, bearer string, accountID int32) map[string]interface{} {
	url := "https://api.futurehack.virtusa.dev/smartbankaccounts/1.0/account/id"
	values := map[string]interface{}{
		"accountId": accountID,
	}

	return genericRequest(client, url, values, bearer)
}

func getBalance(emid string) float64 {
	accountID := accounts[emid].ID
	token := getBearerToken(http.DefaultClient)
	return getBankAccount(http.DefaultClient, token, accountID)["balance"].(float64)
}

func sendVDirhams(sender string, reciever string, amount float64) {
	if senderAcc, ok := accounts[sender]; ok && senderAcc.VDirhams >= amount {
		if recieverAcc, ok := accounts[reciever]; ok {
			fmt.Println("test")
			senderAcc.VDirhams -= amount
			accounts[sender] = senderAcc
			recieverAcc.VDirhams += amount
			accounts[reciever] = recieverAcc
			exportAccountAsJson()
		} else {
			log.Fatal("reciever account doesn't exist")
		}
	} else {
		log.Fatal("sender account doesn't exist, or doesn't have enough VDirhams")
	}
}

func buyVDirhams(emid string, amount float64) {
	// creaing a transaction
	// can't figure out reverse transactions cuz lack of docs
	values := map[string]interface{}{
		"accountId":                   accounts[emid].ID,
		"balance":                     getBalance(emid),
		"balanceCreditDebitIndicator": "Debit",
		"balanceType":                 "OpeningBooked",
		"bankId":                      1006,
		"bankLocation":                "string",
		"bookingDateTime":             "2022-09-30",
		"transactionType":             "WITHDRAWAL",
		"transactionAmount":           amount,
		"transactionName":             "Transaction 1000000005",
		"counterPartyAccountId":       1000009379,
		"counterPartyBankId":          1007,
		"counterPartyBankLocation":    "Brighouse",
		"currencyCode":                "AED",
		"makerDate":                   "2022-09-30",
		"makerId":                     "1",
		"paymentId":                   1000000147,
		"paymentRefId":                "PID965798595",
		"purpose":                     "Order descriptions",
		"status":                      "completed",
		"transactionReference":        "Withdrawal 1000000005",
	}
	resp := genericRequest(http.DefaultClient, "https://smartbank-account.services.futurehack.virtusa.dev/accounts/transaction", values, getBearerToken(http.DefaultClient))
	if resp != nil {
		tempstrc := accounts[emid]
		tempstrc.VDirhams += amount
		accounts[emid] = tempstrc
		fmt.Println(accounts)
	}
	exportAccountAsJson()
}

func login(emid string, password string) bool {
	if doesAccountExist(emid) {
		if accounts[emid].Password == password {
			return true
		}
	}
	return false
}

func displayUserInformation() {
	for {
		windowsClear()

		fmt.Println("Current VDirhams: ", accounts[currentlyLoggedIn].VDirhams)
		if accounts[currentlyLoggedIn].ID > 0 {
			fmt.Println("Current account balance: ", getBalance(currentlyLoggedIn))
			fmt.Println("1) buy VDirham")
		}
		fmt.Println("2) send VDirham")
		fmt.Println("0) Exit")

		reader := bufio.NewReader(os.Stdin)
		txt := printAndRead("\n", reader)
		if txt == "1" && accounts[currentlyLoggedIn].ID > 0 {
			amt := printAndRead("Enter Amount: ", reader)
			if s, err := strconv.ParseFloat(amt, 64); err == nil {
				buyVDirhams(currentlyLoggedIn, s)
			}
		} else if txt == "2" {
			amt := printAndRead("Enter Amount: ", reader)
			reciever := printAndRead("Enter reciver's emarites id: ", reader)
			if s, err := strconv.ParseFloat(amt, 64); err == nil {
				sendVDirhams(currentlyLoggedIn, reciever, s)
			}
		} else if txt == "0" {
			os.Exit(0)
		}
	}
}
