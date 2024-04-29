package database

import (
	"os"
)

func WriteToDB(path string, dataStruct Serialiser, pass []byte) error {

	data, err := dataStruct.Serialise()
	if err != nil {
		return err
	}
	encryptedData, err := EncryptWithPass(data, pass)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, encryptedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// ReadFromDB
// Reads data from a file and decrypts its content parsing it into the given Serialiser Interface.
func ReadFromDB(path string, dataStruct Serialiser, pass []byte) error {
	encryptedData, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	decryptedData, err := DecryptWithPass(encryptedData, pass)
	if err != nil {
		return err
	}

	err = dataStruct.DeSerialise(decryptedData)
	if err != nil {
		return err
	}

	return nil
}
