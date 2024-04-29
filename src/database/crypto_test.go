package database

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {

	data := []byte("Hello, World!")

	passphrase := []byte("passKey")

	key := ConvertPassphraseToKey(passphrase)

	encryptedBytes, err := Encrypt(data, key)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	decryptedBytes, err := Decrypt(encryptedBytes, key)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	if !bytes.Equal(decryptedBytes, data) {
		t.Errorf("Error: did not match %s != %s", decryptedBytes, data)
		return
	}

	encryptedBytes, err = EncryptWithPass(data, passphrase)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	decryptedBytes, err = DecryptWithPass(encryptedBytes, passphrase)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	if len(decryptedBytes) == 0 {
		t.Errorf("Error: should not be zero length")
		return
	}

	if !bytes.Equal(decryptedBytes, data) {
		t.Errorf("Error: did not match %s != %s", decryptedBytes, data)
		return
	}
}

func TestEncryptDecryptWrongPassphrase(t *testing.T) {

	data := []byte("Hello, world!")

	passphrase := []byte("passKey")
	passphraseWrong := []byte("passKey1")

	key := ConvertPassphraseToKey(passphrase)
	keyWrong := ConvertPassphraseToKey(passphraseWrong)

	encryptedBytes, err := Encrypt(data, key)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	decryptedBytes, err := Decrypt(encryptedBytes, keyWrong)
	if err.Error()[:22] != "invalid padding amount" {
		t.Errorf("Error: %s", err)
		return
	}

	if bytes.Equal(decryptedBytes, data) {
		t.Errorf("Error: results should not match")
		return
	}

	encryptedBytes, err = EncryptWithPass(data, passphrase)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	decryptedBytes, err = DecryptWithPass(encryptedBytes, passphraseWrong)
	if err.Error()[:22] != "invalid padding amount" {
		t.Errorf("Error: %s", err)
		return
	}

	if bytes.Equal(decryptedBytes, data) {
		t.Errorf("Error: results should not match")
		return
	}
}
