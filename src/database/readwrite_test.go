package database

import (
	"bytes"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/logging"
	"github.com/setavenger/gobip352"
	"testing"
)

func init() {
	logging.LoadLoggersMock()
}

// todo expand test for more vectors and cases

func TestReadWriteSuccess(t *testing.T) {
	tmpTestPath := "/tmp/example1"
	passphrase := []byte("passKey")

	//empty32Arr := gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0x00}, 32))
	//empty33Arr := gobip352.ConvertToFixedLength33(bytes.Repeat([]byte{0x00}, 33))

	txid := gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0x01}, 32))
	tweak := gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0x02}, 32))
	pubKey := gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0x03}, 32))

	pubKeyLabel := gobip352.ConvertToFixedLength33(bytes.Repeat([]byte{0x04}, 33))
	tweakLabel := gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0x05}, 32))

	//

	txid2 := gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0xee}, 32))
	tweak2 := gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0xaa}, 32))
	pubKey2 := gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0xbb}, 32))

	pubKeyLabel2 := gobip352.ConvertToFixedLength33(bytes.Repeat([]byte{0x44}, 33))
	tweakLabel2 := gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0x55}, 32))

	label1 := gobip352.Label{
		PubKey:  pubKeyLabel,
		Tweak:   tweakLabel,
		Address: "this_is_my_example_address",
		M:       3,
	}
	label2 := gobip352.Label{
		PubKey:  pubKeyLabel2,
		Tweak:   tweakLabel2,
		Address: "a_different_label_address",
		M:       66,
	}

	var collection = &src.UtxoCollection{
		{
			Txid:         txid,
			Vout:         1,
			Amount:       2_200_200,
			PrivKeyTweak: tweak,
			PubKey:       pubKey,
			State:        src.StateSpent,
			Label:        &label1,
		},
		{
			Txid:         txid2,
			Vout:         33,
			Amount:       21_210_211,
			PrivKeyTweak: tweak2,
			PubKey:       pubKey2,
			State:        src.StateUnknown,
			Label:        &label2,
		},
	}

	err := WriteToDB(tmpTestPath, collection, passphrase)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	var pulledCollection src.UtxoCollection
	err = ReadFromDB(tmpTestPath, &pulledCollection, passphrase)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	if pulledCollection[0].State != src.StateSpent {
		t.Errorf("Error: States for 1 did not match %d != %d", pulledCollection[0].State, src.StateSpent)
		return
	}
	if pulledCollection[0].Amount != 2_200_200 {
		t.Errorf("Error: Amounts for 1 did not match %d != %d", pulledCollection[0].State, 2_200_200)
		return
	}
	if pulledCollection[0].Vout != 1 {
		t.Errorf("Error: Vouts for 1 did not match %d != %d", pulledCollection[0].State, 1)
		return
	}
	if !bytes.Equal(pulledCollection[0].Txid[:], txid[:]) {
		t.Errorf("Error: Txids for 1 did not match %x != %x", pulledCollection[0].Txid, txid)
		return
	}
	if !bytes.Equal(pulledCollection[0].PrivKeyTweak[:], tweak[:]) {
		t.Errorf("Error: PrivKeyTweaks for 1 did not match %x != %x", pulledCollection[0].PrivKeyTweak, tweak)
		return
	}
	if !bytes.Equal(pulledCollection[0].PubKey[:], pubKey[:]) {
		t.Errorf("Error: PubKeys for 1 did not match %x != %x", pulledCollection[0].PubKey, pubKey)
		return
	}
	if !bytes.Equal(pulledCollection[0].Label.PubKey[:], pubKeyLabel[:]) {
		t.Errorf("Error: Label PubKeys for 1 did not match %x != %x", pulledCollection[0].Label.PubKey, pubKeyLabel)
		return
	}
	if !bytes.Equal(pulledCollection[0].Label.Tweak[:], tweakLabel[:]) {
		t.Errorf("Error: Label Tweak for 1 did not match %x != %x", pulledCollection[0].Label.Tweak, tweakLabel)
		return
	}
	if pulledCollection[0].Label.Address != "this_is_my_example_address" {
		t.Errorf("Error: Label Address for 1 did not match %s != %s", pulledCollection[0].Label.Address, "this_is_my_example_address")
		return
	}
	if pulledCollection[0].Label.M != 3 {
		t.Errorf("Error: Label M for 1 did not match %d != %d", pulledCollection[0].Label.M, 3)
		return
	}

	// check second UTXO
	if pulledCollection[1].State != src.StateUnknown {
		t.Errorf("Error: States for 2 did not match %d != %d", pulledCollection[1].State, src.StateSpent)
		return
	}
	if pulledCollection[1].Amount != 21_210_211 {
		t.Errorf("Error: Amounts for 2 did not match %d != %d", pulledCollection[1].State, 21_210_211)
		return
	}
	if pulledCollection[1].Vout != 33 {
		t.Errorf("Error: Vouts for 2 did not match %d != %d", pulledCollection[1].State, 33)
		return
	}
	if !bytes.Equal(pulledCollection[1].Txid[:], txid2[:]) {
		t.Errorf("Error: Txids for 2 did not match %x != %x", pulledCollection[1].Txid, txid)
		return
	}
	if !bytes.Equal(pulledCollection[1].PrivKeyTweak[:], tweak2[:]) {
		t.Errorf("Error: PrivKeyTweaks for 2 did not match %x != %x", pulledCollection[1].PrivKeyTweak, tweak)
		return
	}
	if !bytes.Equal(pulledCollection[1].PubKey[:], pubKey2[:]) {
		t.Errorf("Error: PubKeys for 2 did not match %x != %x", pulledCollection[1].PubKey, pubKey)
		return
	}
	if !bytes.Equal(pulledCollection[1].Label.PubKey[:], pubKeyLabel2[:]) {
		t.Errorf("Error: Label PubKeys for 2 did not match %x != %x", pulledCollection[1].Label.PubKey, pubKeyLabel2)
		return
	}
	if !bytes.Equal(pulledCollection[1].Label.Tweak[:], tweakLabel2[:]) {
		t.Errorf("Error: Label Tweak for 2 did not match %x != %x", pulledCollection[1].Label.Tweak, tweakLabel2)
		return
	}
	if pulledCollection[1].Label.Address != "a_different_label_address" {
		t.Errorf("Error: Label Address for 2 did not match %s != %s", pulledCollection[1].Label.Address, "a_different_label_address")
		return
	}
	if pulledCollection[1].Label.M != 66 {
		t.Errorf("Error: Label M for 2 did not match %d != %d", pulledCollection[1].Label.M, 66)
		return
	}

}

func TestWrongPass(t *testing.T) {
	tmpTestPath := "/tmp/example2"
	passphrase := []byte("passKey")
	passphraseWrong := []byte("passKey1")

	txid := gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0x01}, 32))

	txid2 := gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0xee}, 32))

	var collection = &src.UtxoCollection{
		{
			Txid: txid,
			Vout: 1,
		},
		{
			Txid: txid2,
			Vout: 33,
		},
	}

	err := WriteToDB(tmpTestPath, collection, passphrase)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	var pulledCollection src.UtxoCollection

	err = ReadFromDB(tmpTestPath, &pulledCollection, passphraseWrong)
	if err == nil {
		t.Errorf("Error: should have throughn an error")
		return
	}

	if pulledCollection != nil {
		t.Errorf("Error: This should be nil")
		return
	}
}

func TestWrongPath(t *testing.T) {
	tmpTestPath := "/tmp/example3"
	tmpTestPathWrong := "/tmp/example3-wrong"
	passphrase := []byte("passKey")

	txid := gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0x01}, 32))

	txid2 := gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0xee}, 32))

	var collection = &src.UtxoCollection{
		{
			Txid: txid,
			Vout: 1,
		},
		{
			Txid: txid2,
			Vout: 33,
		},
	}

	err := WriteToDB(tmpTestPath, collection, passphrase)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	var pulledCollection src.UtxoCollection

	err = ReadFromDB(tmpTestPathWrong, &pulledCollection, passphrase)
	if err != nil && err.Error() != "open /tmp/example3-wrong: no such file or directory" {
		t.Errorf("Error: Should throw not exist error; threw: %s", err)
		return
	}

	if pulledCollection != nil {
		t.Errorf("Error: This should be nil")
		return
	}
}

func TestInterfaceNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	tmpTestPath := "/tmp/example3"
	passphrase := []byte("passKey")

	txid := gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0x01}, 32))

	txid2 := gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0xee}, 32))

	var collection = &src.UtxoCollection{
		{
			Txid: txid,
			Vout: 1,
		},
		{
			Txid: txid2,
			Vout: 33,
		},
	}

	err := WriteToDB(tmpTestPath, collection, passphrase)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	err = ReadFromDB(tmpTestPath, nil, passphrase)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
}
