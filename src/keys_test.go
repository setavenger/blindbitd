package src

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/tyler-smith/go-bip39"
	"testing"
)

func TestDeriveKeysFromMaster(t *testing.T) {

	scanKeyTarget, _ := hex.DecodeString("78e7fd7d2b7a2c1456709d147021a122d2dccaafeada040cc1002083e2833b09")
	spendKeyTarget, _ := hex.DecodeString("c88567742d5019d7ccc81f6e82cef8ef01997a6a3761cc9166036b580549539b")

	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	seed := bip39.NewSeed(mnemonic, "")

	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		t.Errorf("error creating master key: %v", err)
		return
	}

	keys, err := DeriveKeysFromMaster(master)
	if err != nil {
		t.Errorf("error deriving keys: %v", err)
		return
	}

	if !bytes.Equal(keys.ScanSecretKey[:], scanKeyTarget) {
		t.Errorf("error deriving keys: expected %v, got %v", scanKeyTarget, keys.ScanSecretKey)
		return
	}
	if !bytes.Equal(keys.SpendSecretKey[:], spendKeyTarget) {
		t.Errorf("error deriving keys: expected %v, got %v", scanKeyTarget, keys.ScanSecretKey)
		return
	}
}
