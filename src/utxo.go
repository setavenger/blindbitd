package src

import (
	"encoding/json"
	"github.com/setavenger/gobip352"
)

type OwnedUTXO struct {
	Txid         [32]byte        `json:"txid,omitempty"`
	Vout         uint32          `json:"vout,omitempty"`
	Amount       uint64          `json:"amount"`
	PrivKeyTweak [32]byte        `json:"priv_key_tweak,omitempty"`
	PubKey       [32]byte        `json:"pub_key,omitempty"`
	Timestamp    uint64          `json:"timestamp,omitempty"`
	State        UTXOState       `json:"utxo_state,omitempty"`
	Label        *gobip352.Label `json:"label"` // the pubKey associated with the label
}

func (u *OwnedUTXO) LabelPubKey() []byte {
	if u.Label != nil {
		return u.Label.PubKey[:]
	} else {
		return nil
	}
}

type UtxoCollection []OwnedUTXO

func (c *UtxoCollection) Serialise() ([]byte, error) {
	return json.Marshal(c)
}

func (c *UtxoCollection) DeSerialise(data []byte) error {
	// either write directly or do some extra manipulation
	err := json.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return nil
}
