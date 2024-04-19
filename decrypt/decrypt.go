package decrypt

import (
	"challenge-go/cipher"
	"fmt"
	"io/ioutil"
	"os"
)

func DecryptCSVFile(fileName string) {

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
	decryptedData := cipher.Rot128Decrypt(encryptedData)

	//write the decrypted data into a new file
	decryptedFileName := "./data/decrypted_fng.1000.csv.rot128"
	os.WriteFile(decryptedFileName, decryptedData, 0644)
	if err != nil {
		fmt.Printf("error in decrpting file %v", err)
	}

}
