package database

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// todo rename functions from text to more accurate naming

func ConvertPassphraseToKey(pass []byte) []byte {
	hash := sha256.Sum256(pass)
	return hash[:]
}

func EncryptWithPass(data, pass []byte) ([]byte, error) {
	key := ConvertPassphraseToKey(pass)
	return Encrypt(data, key)
}

func DecryptWithPass(encryptedData, pass []byte) ([]byte, error) {
	key := ConvertPassphraseToKey(pass)
	return Decrypt(encryptedData, key)
}

// Encrypt
// gotten from https://github.com/nbd-wtf/go-nostr/blob/master/nip04/nip04.go
func Encrypt(message []byte, key []byte) ([]byte, error) {
	// block size is 16 bytes
	iv := make([]byte, 16)

	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("error creating initization vector: %w", err)
	}

	// automatically picks aes-256 based on key length (32 bytes)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating block cipher: %w", err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)

	plaintext := message

	// add padding
	base := len(plaintext)

	// this will be a number between 1 and 16 (including), never 0
	padding := block.BlockSize() - base%block.BlockSize()

	// encode the padding in all the padding bytes themselves
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	paddedMsgBytes := append(plaintext, padtext...)

	ciphertext := make([]byte, len(paddedMsgBytes))
	mode.CryptBlocks(ciphertext, paddedMsgBytes)

	var result []byte
	result = append(result, iv...)
	result = append(result, ciphertext...)

	return result, nil
}

// Decrypt decrypts a content string using the shared secret key.
// The inverse operation to message -> Encrypt(message, key).
func Decrypt(content []byte, key []byte) ([]byte, error) {
	ciphertext := content[16:]

	iv := content[:16]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating block cipher: %w", err)
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// remove padding
	var message []byte
	var plaintextLen = len(plaintext)

	if plaintextLen > 0 {
		padding := int(plaintext[plaintextLen-1]) // the padding amount is encoded in the padding bytes themselves
		if padding > plaintextLen {
			return nil, fmt.Errorf("invalid padding amount: %d", padding)
		}
		message = plaintext[0 : plaintextLen-padding]
	}

	return message, nil
}
