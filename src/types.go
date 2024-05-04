package src

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/setavenger/gobip352"
)

var Empty32Arr = gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0x00}, 32))
var Empty33Arr = gobip352.ConvertToFixedLength33(bytes.Repeat([]byte{0x00}, 33))

type UTXOState int8

const (
	StateUnknown UTXOState = iota - 1
	StateUnconfirmed
	StateUnspent
	StateSpent
	StateUnconfirmedSpent
)

type Recipient struct {
	Address    string
	PkScript   []byte
	Amount     int64
	Annotation string
	Data       map[string]any
}

type Label struct {
	// todo add created_at field
	Comment         string `json:"comment"`
	*gobip352.Label `json:"label"`
}

// Addresses maps the address to an annotation the annotation might be empty
type Addresses map[string]string

// LabelsMapping
// the key is the label's pubKey, the value is the Label data
type LabelsMapping map[[33]byte]Label

func (lm *LabelsMapping) MarshalJSON() ([]byte, error) {
	// Convert map to a type that can be marshaled by the standard JSON package
	aux := make(map[string]Label)
	for k, v := range *lm {
		key := fmt.Sprintf("%x", k) // Convert byte array to hex string
		aux[key] = v
	}
	return json.Marshal(aux)
}

func (lm *LabelsMapping) UnmarshalJSON(data []byte) error {
	aux := make(map[string]Label)
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	*lm = make(LabelsMapping)
	for k, v := range aux {
		var key [33]byte
		_, err := hex.Decode(key[:], []byte(k))
		if err != nil {
			return err
		}
		(*lm)[key] = v
	}
	return nil
}

// UTXOMapping
// the key is the utxos (txid||vout)
type UTXOMapping map[[36]byte]struct{}

func (um *UTXOMapping) MarshalJSON() ([]byte, error) {
	// Convert map to a type that can be marshaled by the standard JSON package
	aux := make(map[string]struct{})
	for k, v := range *um {
		key := fmt.Sprintf("%x", k) // Convert byte array to hex string
		aux[key] = v
	}
	return json.Marshal(aux)
}

func (um *UTXOMapping) UnmarshalJSON(data []byte) error {
	aux := make(map[string]struct{})
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	*um = make(UTXOMapping)
	for k, v := range aux {
		var key [36]byte
		_, err := hex.Decode(key[:], []byte(k))
		if err != nil {
			return err
		}
		(*um)[key] = v
	}
	return nil
}
