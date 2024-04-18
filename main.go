package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type Card struct {
	Name            string `json:"name"`
	Number          string `json:"number"`
	SecurityCode    string `json:"security_code"`
	ExpirationMonth int    `json:"expiration_month"`
	ExpirationYear  int    `json:"expiration_year"`
}

type ChargeRequest struct {
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	Card        Card   `json:"card"`
}

const (
	OmisePublicKey = "pkey_test_5zfn4aryqab1ylxwkvi"
	OmiseSecretKey = "skey_test_5zfn4asysrg03zmw3f5"
	Currency       = "THB"
)

func main() {
	fileName := "./data/fng.1000.csv.rot128"
	decryptCSVFile(fileName)
	decryptedFileName := "./data/decrypted_fng.1000.csv.rot128"

	csvFile, err := os.Open(decryptedFileName)
	if err != nil {
		fmt.Println("Error opening the decrypted CSV file:", err)
		return
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = 6 // Name,AmountSubunits,CCNumber,CVV,ExpMonth,ExpYear

	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading the CSV file:", err)
		return
	}

	var totalReceived, successfullyDonated, faultyDonation float64
	donors := make(map[string]float64)

	for _, line := range lines {
		name := line[0]
		amountSubunits, err := strconv.Atoi(line[1])
		if err != nil {
			faultyDonation += float64(amountSubunits)
			continue
		}

		amount := (amountSubunits)
		ccNumber := line[2]
		cvv := line[3]
		expMonth, _ := strconv.Atoi(line[4])
		expYear, _ := strconv.Atoi(line[5])

		// Create time object for checking expiration date
		expDate := time.Date(expYear, time.Month(expMonth), 1, 0, 0, 0, 0, time.UTC)

		// Check if expiration date is in the past
		if expDate.Before(time.Now()) {
			fmt.Printf("(400/invalid_card) expiration date cannot be in the past")
			continue
		}

		chargeRequest := ChargeRequest{
			Description: fmt.Sprintf("Donation from %s", name),
			Amount:      amount,
			Currency:    "THB",
			Card: Card{
				Name:            name,
				Number:          ccNumber,
				SecurityCode:    cvv,
				ExpirationMonth: expMonth,
				ExpirationYear:  expYear,
			},
		}

		// Initialize Omise client
		client, e := omise.NewClient(OmisePublicKey, OmiseSecretKey)
		if e != nil {
			fmt.Printf("error in calling a omise client %v", e)
		}

		//Create a token to use in create Charge API
		token, create := &omise.Token{}, &operations.CreateToken{
			Name:            chargeRequest.Card.Name,
			Number:          chargeRequest.Card.Number,
			ExpirationMonth: time.Month(expMonth),
			ExpirationYear:  chargeRequest.Card.ExpirationYear,
			SecurityCode:    chargeRequest.Card.SecurityCode,
		}
		if e := client.Do(token, create); e != nil {
			fmt.Printf("error in prcessing payment %v", e)
			panic(e)
		}

		//Create charge
		createCharge := &operations.CreateCharge{
			Amount:   int64(amount), // ฿ 1,000.00
			Currency: "THB",
			Card:     token.ID,
		}

		chargeResponse := &omise.Charge{}
		err = client.Do(chargeResponse, createCharge)
		if err != nil {
			faultyDonation += float64(amount)
			fmt.Printf("error in creating charge %v", err)
			continue
		}

		//fmt.Printf("Successfully donated THB %.2f for %s\n", float64(chargeResponse.Amount)/100, name)
		successfullyDonated += float64(chargeResponse.Amount)
		donors[line[0]] += float64(amount)

	}

	totalReceived = successfullyDonated + faultyDonation
	averagePerPerson := totalReceived / float64(len(lines))

	printSummary(donors, totalReceived, successfullyDonated, faultyDonation, averagePerPerson)

}

func printTopDonors(donors map[string]float64) {
	type donor struct {
		name   string
		amount float64
	}

	var topDonors []donor
	for name, amount := range donors {
		topDonors = append(topDonors, donor{name, amount})
	}

	sort.Slice(topDonors, func(i, j int) bool {
		return topDonors[i].amount > topDonors[j].amount
	})

	for _, d := range topDonors[:int(math.Min(3, float64(len(topDonors))))] {
		fmt.Printf("\t\t%s\n", d.name)
	}
}

func rot128Decrypt(data []byte) []byte {
	for index := range data {
		data[index] -= 128
	}
	return data
}

func decryptCSVFile(fileName string) {

	//open encrypted file for reading
	encryptedFile, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("error in opening a file %v", err)
	}
	defer encryptedFile.Close()

	//read an encryptedFile
	encryptedData, err2 := ioutil.ReadAll(encryptedFile)

	if err2 != nil {
		fmt.Println("Error in reading file", err)
		return
	}
	//decrypt data
	decryptedData := rot128Decrypt(encryptedData)

	//write the decrypted data into a new file
	decryptedFileName := "./data/decrypted_fng.1000.csv.rot128"
	os.WriteFile(decryptedFileName, decryptedData, 0644)
	if err != nil {
		fmt.Printf("error in decrpting file %v", err)
	}

	//fmt.Printf("decryption is completed %s", decryptedFileName)

}

func printSummary(donors map[string]float64, totalReceived, successfullyDonated, faultyDonation, averagePerPerson float64) {
	fmt.Println("\n\nperforming donations...")
	fmt.Println("done.")
	fmt.Printf("\nTotal received: THB %.2f\n", float64(totalReceived)/100)
	fmt.Printf("Successfully donated: THB %.2f\n", float64(successfullyDonated)/100)
	fmt.Printf("Faulty donation: THB %.2f\n", float64(faultyDonation)/100)
	fmt.Printf("\n\taverage per person: THB %.2f\n", averagePerPerson)
	fmt.Println("\ttop donors:")
	printTopDonors(donors)
}
