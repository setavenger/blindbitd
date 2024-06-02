package src

import (
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

type item struct {
	mnemonic       string
	passphrase     string
	scanKeyTarget  string
	spendKeyTarget string
}

var testData = []item{
	{
		mnemonic:       "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		passphrase:     "",
		scanKeyTarget:  "78e7fd7d2b7a2c1456709d147021a122d2dccaafeada040cc1002083e2833b09",
		spendKeyTarget: "c88567742d5019d7ccc81f6e82cef8ef01997a6a3761cc9166036b580549539b",
	},
	{
		mnemonic:       "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		passphrase:     "blindbitd",
		scanKeyTarget:  "0324f67bd8ea863daf22545b6dc99829e5b30b6891b7dbddae5a43a26283c8fb",
		spendKeyTarget: "79d7b685c6d82fca04fff325b970b4e2a2117995662af59179b2b7065a7f5bca",
	},
}

func init() {
	ChainParams = &chaincfg.MainNetParams
}

func TestKeysFromMnemonic(t *testing.T) {
	for _, data := range testData {
		keys, err := KeysFromMnemonic(data.mnemonic, data.passphrase)
		if err != nil {
			t.Errorf("error deriving keys: %v", err)
			return
		}

		scanKey := hex.EncodeToString(keys.ScanSecretKey[:])
		spendKey := hex.EncodeToString(keys.SpendSecretKey[:])

		if data.scanKeyTarget != scanKey {
			t.Errorf("error deriving scan key: expected %v, got %v", data.scanKeyTarget, scanKey)
			return
		}

		if data.spendKeyTarget != spendKey {
			t.Errorf("error deriving spend key: expected %v, got %v", data.spendKeyTarget, spendKey)
			return
		}
	}
}

func TestCreateNewKeys(t *testing.T) {
	for _, data := range testData {
		keys, err := CreateNewKeys(data.passphrase)
		if err != nil {
			t.Errorf("error creating keys: %v", err)
			return
		}

		keysTarget, err := KeysFromMnemonic(keys.Mnemonic, data.passphrase)
		if err != nil {
			t.Errorf("error deriving keys: %v", err)
			return
		}

		scanKeyTarget := hex.EncodeToString(keysTarget.ScanSecretKey[:])
		spendKeyTarget := hex.EncodeToString(keysTarget.SpendSecretKey[:])

		scanKey := hex.EncodeToString(keys.ScanSecretKey[:])
		spendKey := hex.EncodeToString(keys.SpendSecretKey[:])

		if scanKeyTarget != scanKey {
			t.Errorf("error deriving scan key: expected %v, got %v", scanKeyTarget, scanKey)
			return
		}

		if spendKeyTarget != spendKey {
			t.Errorf("error deriving spend key: expected %v, got %v", spendKeyTarget, spendKey)
			return
		}
	}
}
