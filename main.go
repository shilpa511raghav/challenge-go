package main

import (
	"fmt"

	"challenge-go/decrypt"
	"challenge-go/payment"
)

func main() {
	encyptedFileName := "./data/fng.1000.csv.rot128"
	decryptedFileName := "./data/decrypted_fng.1000.csv.rot128"
	err := decrypt.DecryptCSVFile(encyptedFileName, decryptedFileName)
	if err != nil {
		fmt.Println("Error in decryption process:", err)
		return
	}
	donors, totalReceived, successfullyDonated, faultyDonation, averagePerPerson, err := payment.ProcessPayments(decryptedFileName)
	if err != nil {
		fmt.Println("Error processing payments:", err)
		return
	}
	payment.PrintSummary(donors, totalReceived, successfullyDonated, faultyDonation, averagePerPerson)

}
