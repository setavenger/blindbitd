package src

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/setavenger/blindbitd/src/logging"
	"github.com/setavenger/go-bip352"
)

type OwnedUTXO struct {
	Txid         [32]byte      `json:"txid,omitempty"`
	Vout         uint32        `json:"vout,omitempty"`
	Amount       uint64        `json:"amount"`
	PrivKeyTweak [32]byte      `json:"priv_key_tweak,omitempty"`
	PubKey       [32]byte      `json:"pub_key,omitempty"`
	Timestamp    uint64        `json:"timestamp,omitempty"`
	State        UTXOState     `json:"utxo_state,omitempty"`
	Label        *bip352.Label `json:"label"` // the pubKey associated with the label
}

func (u *OwnedUTXO) LabelPubKey() []byte {
	if u.Label != nil {
		return u.Label.PubKey[:]
	} else {
		return nil
	}
}

func (u *OwnedUTXO) LabelComment(mapping LabelsMapping) *string {
	if mapping == nil {
		logging.ErrorLogger.Println("labels mapping is nil")
		panic(errors.New("labels mapping is nil")) // todo change to return ""/nil after initial test phase
	}
	if u.Label == nil {
		return nil
	}
	label, ok := mapping[u.Label.PubKey]
	if !ok {
		logging.ErrorLogger.Println("label not found")
		return nil
	}
	return &label.Comment
}

type UtxoCollection []*OwnedUTXO

func (c *UtxoCollection) Serialise() ([]byte, error) {
	return json.Marshal(c)
}

func (c *UtxoCollection) DeSerialise(data []byte) error {
	// either write directly or do some extra manipulation
	err := json.Unmarshal(data, c)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}

	return nil
}

func (u *OwnedUTXO) GetKey() ([36]byte, error) {
	var buf bytes.Buffer
	buf.Write(u.Txid[:])
	err := binary.Write(&buf, binary.BigEndian, u.Vout)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return [36]byte{}, err
	}

	var result [36]byte
	copy(result[:], buf.Bytes())

	return result, nil
}
