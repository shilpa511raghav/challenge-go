package main

import (
	"fmt"

	"challenge-go/decrypt"
	"challenge-go/payment"
)

func main() {
	fileName := "./data/fng.1000.csv.rot128"
	decrypt.DecryptCSVFile(fileName)
	decryptedFileName := "./data/decrypted_fng.1000.csv.rot128"
	donors, totalReceived, successfullyDonated, faultyDonation, averagePerPerson, err := payment.ProcessPayments(decryptedFileName)
	if err != nil {
		fmt.Println("Error processing payments:", err)
		return
	}
	payment.PrintSummary(donors, totalReceived, successfullyDonated, faultyDonation, averagePerPerson)

}
