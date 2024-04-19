package payment

import (
	"encoding/csv"
	"fmt"
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
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Card        Card    `json:"card"`
}

const (
	OmisePublicKey = "pkey_test_5zfn4aryqab1ylxwkvi"
	OmiseSecretKey = "skey_test_5zfn4asysrg03zmw3f5"
	Currency       = "THB"
)

func ProcessPayments(decryptedFilename string) (map[string]float64, float64, float64, float64, float64, error) {
	// Open the decrypted CSV file
	csvFile, err := os.Open(decryptedFilename)
	if err != nil {
		return nil, 0, 0, 0, 0, (fmt.Errorf("Error opening the decrypted CSV file: %v", err))
	}
	defer csvFile.Close()

	// Create a CSV reader
	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = 6 // Name,AmountSubunits,CCNumber,CVV,ExpMonth,ExpYear

	// Read all lines from the CSV
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, 0, 0, 0, 0, (fmt.Errorf("Error reading the CSV file: %v", err))
	}

	// Initialize variables to track donation statistics
	var totalReceived, successfullyDonated, faultyDonation float64
	donors := make(map[string]float64)

	// Process each line in the CSV
	for _, line := range lines {
		name := line[0]
		amountSubunits, err := strconv.Atoi(line[1])
		if err != nil {
			faultyDonation += float64(amountSubunits)
			continue
		}

		amount := float64(amountSubunits)
		ccNumber := line[2]
		cvv := line[3]
		expMonth, _ := strconv.Atoi(line[4])
		expYear, _ := strconv.Atoi(line[5])

		// Create time object for checking expiration date
		expDate := time.Date(expYear, time.Month(expMonth), 1, 0, 0, 0, 0, time.UTC)

		// Check if expiration date is in the past
		if expDate.Before(time.Now()) {
			fmt.Printf("(400/invalid_card) expiration date cannot be in the past\n")
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
		client, err := initializeOmiseClient()
		if err != nil {
			return nil, 0, 0, 0, 0, (fmt.Errorf("Error initializing Omise client: %v", err))
		}

		token, err := createToken(client, chargeRequest)
		if err != nil {
			return nil, 0, 0, 0, 0, (fmt.Errorf("Error creating token: %v", err))
		}
		chargeResponse, err := createCharge(client, chargeRequest, token)
		if err != nil {
			return nil, 0, 0, 0, 0, (fmt.Errorf("Error creating charge: %v", err))
		}
		successfullyDonated += float64(chargeResponse.Amount)
		donors[line[0]] += float64(amount)
	}
	totalReceived = successfullyDonated + faultyDonation
	averagePerPerson := (totalReceived) / float64(len(lines)-1)
	return donors, totalReceived, successfullyDonated, faultyDonation, averagePerPerson, nil
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

// initialize omise client
func initializeOmiseClient() (*omise.Client, error) {
	client, err := omise.NewClient(OmisePublicKey, OmiseSecretKey)
	if err != nil {
		fmt.Printf("error in calling an omise client %v\n", err)
		return nil, err
	}
	return client, nil
}

// Create a token
func createToken(client *omise.Client, chargeRequest ChargeRequest) (*omise.Token, error) {

	tokenRequest := &operations.CreateToken{
		Name:            chargeRequest.Card.Name,
		Number:          chargeRequest.Card.Number,
		ExpirationMonth: time.Month(chargeRequest.Card.ExpirationMonth),
		ExpirationYear:  chargeRequest.Card.ExpirationYear,
		SecurityCode:    chargeRequest.Card.SecurityCode,
	}

	token := &omise.Token{}
	if err := client.Do(token, tokenRequest); err != nil {
		fmt.Printf("error in generating token %v\n", err)
		return nil, err
	}

	return token, nil
}

// Create charge
func createCharge(client *omise.Client, chargeRequest ChargeRequest, token *omise.Token) (*omise.Charge, error) {

	createCharge := &operations.CreateCharge{
		Amount:   int64(chargeRequest.Amount),
		Currency: chargeRequest.Currency,
		Card:     token.ID,
	}

	chargeResponse := &omise.Charge{}
	err := client.Do(chargeResponse, createCharge)
	if err != nil {
		fmt.Printf("error in creating charge %v\n", err)
		return nil, err
	}

	return chargeResponse, nil
}

func PrintSummary(donors map[string]float64, totalReceived, successfullyDonated, faultyDonation, averagePerPerson float64) {
	fmt.Println("\n\nperforming donations...")
	fmt.Println("done.")
	fmt.Printf("\nTotal received: THB %.2f\n", float64(totalReceived))
	fmt.Printf("Successfully donated: THB %.2f\n", float64(successfullyDonated))
	fmt.Printf("Faulty donation: THB %.2f\n", float64(faultyDonation))
	fmt.Printf("\n\taverage per person: THB %.2f\n", averagePerPerson)
	fmt.Println("\ttop donors:")
	printTopDonors(donors)
}
