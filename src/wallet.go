package src

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/setavenger/gobip352"
)

type UTXOState uint8

const (
	StateSpent UTXOState = iota
	StatePending
	StateUnspent
)

type Wallet struct {
	secretKeyScan  [32]byte
	secretKeySpend [32]byte         // todo might not populate it and only load it on spend
	PubKeyScan     [33]byte         `json:"pub_key_scan"`
	PubKeySpend    [33]byte         `json:"pub_key_spend"`
	BirthHeight    uint64           `json:"birth_height,omitempty"`
	LastScan       uint64           `json:"last_scan,omitempty"`
	UTXOs          []OwnedUTXO      `json:"utxos,omitempty"`
	Labels         []gobip352.Label `json:"labels"`
	// ChangeLabel is separate in order to make it clear that it's special and is not just shown like other labels
	ChangeLabel gobip352.Label `json:"change_label"`
	// nextLabelM indicates which m will be used to derive the next label
	nextLabelM     uint32
	DustLimit      uint64     `json:"dust_limit"`
	PubKeysToWatch [][32]byte `json:"pub_keys_to_watch"`
	Addresses      `json:"addresses"`
	LabelsMapping  `json:"labels_mapping"`
}

type OwnedUTXO struct {
	Txid               [32]byte  `json:"txid,omitempty"`
	Vout               uint32    `json:"vout,omitempty"`
	Amount             uint64    `json:"amount"`
	PrivKeyTweak       [32]byte  `json:"priv_key_tweak,omitempty"`
	PubKey             [32]byte  `json:"pub_key,omitempty"`
	TimestampConfirmed uint64    `json:"timestamp_confirmed,omitempty"`
	State              UTXOState `json:"utxo_state,omitempty"`
	Label              uint32    `json:"label"` // the pubKey associated with the label
}

// Addresses maps the address to an annotation the annotation might be empty
type Addresses map[string]string

// LabelsMapping
// the key is the label's pubKey, the value is the Label data
type LabelsMapping map[[33]byte]gobip352.Label

func NewWallet() *Wallet {
	return &Wallet{Addresses: Addresses{}, nextLabelM: 1}
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

func (w *Wallet) GenerateNewLabel() (string, error) {
	// we don't allow m = 0 as it's reserved for the change label and should also never be exposed
	if w.nextLabelM == 0 {
		w.nextLabelM = 1
	}

	m := w.nextLabelM
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
		return "", LabelAlreadyExistsErr{}
	}

	w.Addresses[address] = fmt.Sprintf("label: %d", m)
	w.nextLabelM++

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

	return address, err
}

func (w *Wallet) ConvertOwnedUTXOIntoVin(utxo OwnedUTXO) gobip352.Vin {
	fullSecretKey := gobip352.AddPrivateKeys(w.secretKeySpend, utxo.PrivKeyTweak)
	vin := gobip352.Vin{
		Txid:         utxo.Txid,
		Vout:         utxo.Vout,
		ScriptPubKey: append([]byte{0x51, 0x20}, utxo.PubKey[:]...),
		SecretKey:    &fullSecretKey,
		Taproot:      true,
	}
	return vin
}
