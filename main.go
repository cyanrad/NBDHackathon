package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func printAndRead(text string, reader *bufio.Reader) string {
	fmt.Print(text)
	txt, _ := reader.ReadString('\n')
	return txt[:len(txt)-2]
}

func windowsClear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	loadAccountsFromJson()
	reader := bufio.NewReader(os.Stdin)
	defer exportAccountAsJson()

	// logging in users
	fmt.Printf("Hello and Welcome to QwikPay\n(please know that this is a demo, and not\n representative of the actual flow of the program\n note that in the accounts.json ID=-1 means bankless account)\nPlease select one of the options:\n1) Login\n2) register\n")
	txt, _ := reader.ReadString('\n')
	txt = txt[:len(txt)-2]
	if txt == "1" {
		for {
			emid := printAndRead("Enter your emarits ID: ", reader)
			pass := printAndRead("Enter your password: ", reader)
			if login(emid, pass) {
				fmt.Println("Logging in :)")
				currentlyLoggedIn = emid
				break
			} else {
				fmt.Println("incorrect try again")
			}
		}
	} else if txt == "2" {
		emid := printAndRead("Enter your emarits ID: ", reader)
		pass := printAndRead("Enter your password: ", reader)
		hasBankString := printAndRead("Do you have a bank(y/n): ", reader)
		hasBank := false
		if hasBankString == "y" {
			hasBank = true
		}

		if createAccount(emid, pass, hasBank) {
			fmt.Println("Congratulations, logging in right now")
			currentlyLoggedIn = emid
		} else {
			log.Fatal("Account already exists or incorrect information")
		}
	} else {
		log.Fatal("Unknown input")
	}

	displayUserInformation()
}
