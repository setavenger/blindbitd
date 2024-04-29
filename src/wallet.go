package src

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/setavenger/gobip352"
)

type Wallet struct {
	secretKeyScan  [32]byte
	secretKeySpend [32]byte         // todo might not populate it and only load it on spend
	PubKeyScan     [33]byte         `json:"pub_key_scan"`
	PubKeySpend    [33]byte         `json:"pub_key_spend"`
	BirthHeight    uint64           `json:"birth_height,omitempty"`
	LastScanHeight uint64           `json:"last_scan,omitempty"`
	UTXOs          UtxoCollection   `json:"utxos,omitempty"`
	Labels         []gobip352.Label `json:"labels"`
	ChangeLabel    gobip352.Label   `json:"change_label"` // ChangeLabel is separate in order to make it clear that it's special and is not just shown like other labels
	NextLabelM     uint32           `json:"next_label_m"` // NextLabelM indicates which m will be used to derive the next label
	DustLimit      uint64           `json:"dust_limit"`
	PubKeysToWatch [][32]byte       `json:"pub_keys_to_watch"`
	Addresses      Addresses        `json:"addresses"`
	LabelsMapping  LabelsMapping    `json:"labels_mapping"` // never show LabelsMapping addresses to the user - it includes the change label which should NEVER be shown to normal users
}

func NewWallet(birthHeight uint64) *Wallet {
	return &Wallet{
		Addresses:     Addresses{},
		LabelsMapping: LabelsMapping{},
		NextLabelM:    1,
		BirthHeight:   200, // todo set to var birthHeight
	}
}

func (w *Wallet) Serialise() ([]byte, error) {
	return json.Marshal(w)
}

func (w *Wallet) DeSerialise(data []byte) error {
	// either write directly or do some extra manipulation
	err := json.Unmarshal(data, w)
	if err != nil {
		return err
	}

	return nil
}

func (w *Wallet) LoadWalletFromKeys(secretKeyScan, secretKeySpend [32]byte) {

	w.secretKeyScan = secretKeyScan
	w.secretKeySpend = secretKeySpend

	_, pubKeyScan := btcec.PrivKeyFromBytes(secretKeyScan[:])
	_, pubKeySpend := btcec.PrivKeyFromBytes(secretKeySpend[:])

	w.PubKeyScan = gobip352.ConvertToFixedLength33(pubKeyScan.SerializeCompressed())
	w.PubKeySpend = gobip352.ConvertToFixedLength33(pubKeySpend.SerializeCompressed())

	return
}

func (w *Wallet) GenerateAddress() (string, error) {
	address, err := gobip352.CreateAddress(w.PubKeyScan, w.PubKeySpend, false, 0)
	if err != nil {
		return "", err
	}
	w.Addresses[address] = "standard"
	return address, err
}

func (w *Wallet) GenerateNewLabel(comment string) (string, error) {
	// we don't allow m = 0 as it's reserved for the change label and should also never be exposed
	if w.NextLabelM == 0 {
		w.NextLabelM = 1
	}

	m := w.NextLabelM
	label, err := gobip352.CreateLabel(w.secretKeyScan, m)
	if err != nil {
		return "", err
	}

	BmKey, err := gobip352.AddPublicKeys(w.PubKeySpend, label.PubKey)
	address, err := gobip352.CreateAddress(w.PubKeyScan, BmKey, false, 0)
	if err != nil {
		return "", err
	}

	label.Address = address

	_, exists := w.LabelsMapping[label.PubKey]

	if exists {
		// users should not create the same label twice
		return "", ErrLabelAlreadyExists
	}

	w.Addresses[address] = fmt.Sprintf("label-%d: %s", m, comment)
	w.NextLabelM++

	w.LabelsMapping[label.PubKey] = Label{Label: label, Comment: comment}
	w.Labels = append(w.Labels, label)
	return address, err
}

func (w *Wallet) GenerateChangeLabel() (string, error) {
	// the change label is always m = 0 as defined in the BIP
	var m uint32 = 0
	label, err := gobip352.CreateLabel(w.secretKeyScan, m)
	if err != nil {
		return "", err
	}

	BmKey, err := gobip352.AddPublicKeys(w.PubKeySpend, label.PubKey)
	address, err := gobip352.CreateAddress(w.PubKeyScan, BmKey, false, 0)
	if err != nil {
		return "", err
	}

	label.Address = address
	w.ChangeLabel = label
	w.LabelsMapping[label.PubKey] = Label{Label: label, Comment: "change"}
	return address, err
}

func ConvertOwnedUTXOIntoVin(utxo OwnedUTXO) gobip352.Vin {
	vin := gobip352.Vin{
		Txid:         utxo.Txid,
		Vout:         utxo.Vout,
		Amount:       utxo.Amount,
		ScriptPubKey: append([]byte{0x51, 0x20}, utxo.PubKey[:]...),
		SecretKey:    &utxo.PrivKeyTweak,
		Taproot:      true,
	}
	return vin
}

// FindLabelByPubKey
// returns the pointer to a Label stored in the wallet, will be nil if none was found.
// This is basically a wrapper function around LabelsMapping but adds the change label.
// Chose this approach to avoid accidentally exposing the change address.
func (w *Wallet) FindLabelByPubKey(pubKey [33]byte) *Label {
	panic("implement me")
	return nil
}

func (w *Wallet) SecretKeyScan() [32]byte {
	return w.secretKeyScan
}

func (w *Wallet) SecretKeySpend() [32]byte {
	return w.secretKeySpend
}

func (w *Wallet) CheckAndInitialiseFields() {
	if w.LabelsMapping == nil {
		w.LabelsMapping = LabelsMapping{}
	}
	if w.Addresses == nil {
		w.Addresses = Addresses{}
	}
}

//func (w *Wallet) MakeWalletReady() error {
//	_, err := w.GenerateAddress()
//	if err != nil {
//		return err
//	}
//	_, err = w.GenerateChangeLabel()
//	if err != nil {
//		return err
//	}
//	var i uint32
//	for i = 1; i < w.NextLabelM; i++ {
//		w.GenerateNewLabel()
//	}
//}
