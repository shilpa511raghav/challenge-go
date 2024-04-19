package decrypt

import (
	"challenge-go/cipher"
	"io/ioutil"
	"os"
)

func DecryptCSVFile(inputFileName string, outputFileName string) (err error) {

	//open encrypted file for reading
	encryptedFile, err := os.Open(inputFileName)
	if err != nil {
		return err
	}
	defer encryptedFile.Close()

	//read an encryptedFile
	encryptedData, err := ioutil.ReadAll(encryptedFile)

	if err != nil {
		return err
	}
	//decrypt data
	decryptedData := cipher.Rot128Decrypt(encryptedData)

	//write the decrypted data into a new file
	err2 := os.WriteFile(outputFileName, decryptedData, 0644)
	if err2 != nil {
		return err2
	}

	return nil
}
