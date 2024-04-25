package src

import "github.com/setavenger/gobip352"

type UTXOState int8

const (
	StateUnknown UTXOState = iota - 1
	StateUnconfirmed
	StateUnspent
	StateSpent
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
	ChangeLabel    gobip352.Label   `json:"change_label"` // ChangeLabel is separate in order to make it clear that it's special and is not just shown like other labels
	NextLabelM     uint32           // NextLabelM indicates which m will be used to derive the next label
	DustLimit      uint64           `json:"dust_limit"`
	PubKeysToWatch [][32]byte       `json:"pub_keys_to_watch"`
	Addresses      `json:"addresses"`
	LabelsMapping  `json:"labels_mapping"`
}

type OwnedUTXO struct {
	Txid               [32]byte        `json:"txid,omitempty"`
	Vout               uint32          `json:"vout,omitempty"`
	Amount             uint64          `json:"amount"`
	PrivKeyTweak       [32]byte        `json:"priv_key_tweak,omitempty"`
	PubKey             [32]byte        `json:"pub_key,omitempty"`
	TimestampConfirmed uint64          `json:"timestamp_confirmed,omitempty"`
	State              UTXOState       `json:"utxo_state,omitempty"`
	Label              *gobip352.Label `json:"label"` // the pubKey associated with the label
}

type Label struct {
	Comment string
	gobip352.Label
}

// Addresses maps the address to an annotation the annotation might be empty
type Addresses map[string]string

// LabelsMapping
// the key is the label's pubKey, the value is the Label data
type LabelsMapping map[[33]byte]Label

type Recipient struct {
	Address    string
	PkScript   []byte
	Amount     int64
	Annotation map[string]any
}

// =============== IPC Transfer Types below =============== //

// todo delete those proto takes care of this

type IpcOwnedUTXO struct {
	Txid               [32]byte  `json:"txid,omitempty"`
	Vout               uint32    `json:"vout,omitempty"`
	Amount             uint64    `json:"amount"`
	PrivKeyTweak       [32]byte  `json:"priv_key_tweak,omitempty"`
	PubKey             [32]byte  `json:"pub_key,omitempty"`
	TimestampConfirmed uint64    `json:"timestamp_confirmed,omitempty"`
	State              UTXOState `json:"utxo_state,omitempty"`
	Label              *[33]byte `json:"label"` // the pubKey associated with the label
}

func ConvertUtxoForIpc(utxo OwnedUTXO) IpcOwnedUTXO {
	var labelKey *[33]byte
	if utxo.Label != nil {
		labelKey = &utxo.Label.PubKey
	}
	return IpcOwnedUTXO{
		Txid:               utxo.Txid,
		Vout:               utxo.Vout,
		Amount:             utxo.Amount,
		PrivKeyTweak:       utxo.PrivKeyTweak,
		PubKey:             utxo.PubKey,
		TimestampConfirmed: utxo.TimestampConfirmed,
		State:              utxo.State,
		Label:              labelKey,
	}
}

//
