package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/setavenger/blindbitd/src/logging"
	"github.com/setavenger/go-bip352"
)

// StandardAddressComment Identifier for the base non-labelled address
const StandardAddressComment = "standard"

type Wallet struct {
	secretKeyScan  [32]byte
	secretKeySpend [32]byte        // todo might not populate it and only load it on spend
	PubKeyScan     [33]byte        `json:"pub_key_scan"`
	PubKeySpend    [33]byte        `json:"pub_key_spend"`
	BirthHeight    uint64          `json:"birth_height,omitempty"`
	LastScanHeight uint64          `json:"last_scan,omitempty"`
	UTXOs          UtxoCollection  `json:"utxos,omitempty"`
	Labels         []*bip352.Label `json:"labels"`            // Labels contains all labels except for the change label
	ChangeLabel    *bip352.Label   `json:"change_label"`      // ChangeLabel is separate in order to make it clear that it's special and is not just shown like other labels
	NextLabelM     uint32          `json:"next_label_m"`      // NextLabelM indicates which m will be used to derive the next label
	PubKeysToWatch [][32]byte      `json:"pub_keys_to_watch"` // todo this can be taken out, we can check for new payments to old pub-keys via electrum
	Addresses      Addresses       `json:"addresses"`
	LabelsMapping  LabelsMapping   `json:"labels_mapping"` // never show LabelsMapping addresses to the user - it includes the change label which should NEVER be shown to normal users
	UTXOMapping    UTXOMapping     `json:"utxo_mapping"`   // used to keep track of utxos and not add the same twice
}

func NewWallet(birthHeight uint64) *Wallet {
	return &Wallet{
		Addresses:      Addresses{},
		LabelsMapping:  LabelsMapping{},
		LastScanHeight: 1, // always at least 1 to avoid genesis block
		NextLabelM:     1,
		BirthHeight:    birthHeight,
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

func (w *Wallet) LoadKeys(secretKeyScan, secretKeySpend [32]byte) {

	w.secretKeyScan = secretKeyScan
	w.secretKeySpend = secretKeySpend

	_, pubKeyScan := btcec.PrivKeyFromBytes(secretKeyScan[:])
	_, pubKeySpend := btcec.PrivKeyFromBytes(secretKeySpend[:])

	w.PubKeyScan = bip352.ConvertToFixedLength33(pubKeyScan.SerializeCompressed())
	w.PubKeySpend = bip352.ConvertToFixedLength33(pubKeySpend.SerializeCompressed())

	return
}

func (w *Wallet) GenerateAddress() (string, error) {
	var mainnet bool
	if ChainParams.Name == chaincfg.MainNetParams.Name {
		mainnet = true
	}
	address, err := bip352.CreateAddress(w.PubKeyScan, w.PubKeySpend, mainnet, 0)
	if err != nil {
		return "", err
	}
	w.Addresses[address] = StandardAddressComment
	return address, err
}

func (w *Wallet) GenerateNewLabel(comment string) (*Label, error) {
	var mainnet bool
	if ChainParams.Name == chaincfg.MainNetParams.Name {
		mainnet = true
	}
	// we don't allow m = 0 as it's reserved for the change label and should also never be exposed
	if w.NextLabelM == 0 {
		w.NextLabelM = 1
	}

	m := w.NextLabelM
	label, err := bip352.CreateLabel(w.secretKeyScan, m)
	if err != nil {
		return nil, err
	}

	BmKey, err := bip352.AddPublicKeys(w.PubKeySpend, label.PubKey)
	address, err := bip352.CreateAddress(w.PubKeyScan, BmKey, mainnet, 0)
	if err != nil {
		return nil, err
	}

	label.Address = address

	_, exists := w.LabelsMapping[label.PubKey]

	if exists {
		// users should not create the same label twice
		return nil, ErrLabelAlreadyExists
	}

	w.Addresses[address] = fmt.Sprintf("label-%d: %s", m, comment)
	w.NextLabelM++

	wideLabel := Label{Label: &label, Comment: comment}
	w.LabelsMapping[label.PubKey] = wideLabel
	w.Labels = append(w.Labels, &label)
	return &wideLabel, err
}

func (w *Wallet) GenerateChangeLabel() (string, error) {
	var mainnet bool
	if ChainParams.Name == chaincfg.MainNetParams.Name {
		mainnet = true
	}
	// the change label is always m = 0 as defined in the BIP
	var m uint32 = 0
	label, err := bip352.CreateLabel(w.secretKeyScan, m)
	if err != nil {
		return "", err
	}

	BmKey, err := bip352.AddPublicKeys(w.PubKeySpend, label.PubKey)
	address, err := bip352.CreateAddress(w.PubKeyScan, BmKey, mainnet, 0)
	if err != nil {
		return "", err
	}

	label.Address = address
	w.ChangeLabel = &label
	w.LabelsMapping[label.PubKey] = Label{Label: &label, Comment: "change"}
	return address, err
}

func (w *Wallet) AddUTXOs(utxos []*OwnedUTXO) error {
	for _, utxo := range utxos {
		key, err := utxo.GetKey()
		if err != nil {
			logging.ErrorLogger.Println(err)
			return err
		}
		_, exists := w.UTXOMapping[key]
		if exists {
			continue
		}

		w.UTXOs = append(w.UTXOs, utxo)
	}

	for _, utxo := range utxos {
		w.PubKeysToWatch = append(w.PubKeysToWatch, utxo.PubKey)
	}
	return nil
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

func (w *Wallet) FreeBalance() uint64 {
	var balance uint64 = 0
	for _, utxo := range w.UTXOs {
		if utxo.State == StateUnspent {
			balance += utxo.Amount
		}
	}
	return balance
}

func (w *Wallet) GetFreeUTXOs() UtxoCollection {
	var utxos UtxoCollection
	for _, utxo := range w.UTXOs {
		if utxo.State == StateUnspent {
			utxos = append(utxos, utxo)
		}
	}
	return utxos
}

func (w *Wallet) GetUTXOsByStates(states ...UTXOState) UtxoCollection {
	var utxos UtxoCollection
	for _, utxo := range w.UTXOs {
		for _, state := range states {
			if utxo.State == state {
				utxos = append(utxos, utxo)
			}
		}
	}
	return utxos
}

func (w *Wallet) CheckAndInitialiseFields() error {
	secretKeyScan := w.SecretKeyScan()
	if bytes.Equal(secretKeyScan[:], Empty32Arr[:]) {
		return errors.New("empty scan secret key")
	}

	secretKeySpend := w.SecretKeySpend()
	if bytes.Equal(secretKeySpend[:], Empty32Arr[:]) {
		return errors.New("empty spend secret key")
	}

	if bytes.Equal(w.PubKeyScan[:], Empty33Arr[:]) {
		// if the secret keys are not zero then the pubKeys should be generated without problems
		w.LoadKeys(w.SecretKeyScan(), w.SecretKeySpend())
	}

	if bytes.Equal(w.PubKeySpend[:], Empty33Arr[:]) {
		// if the secret keys are not zero then the pubKeys should be generated without problems
		w.LoadKeys(w.SecretKeyScan(), w.SecretKeySpend())
	}

	if w.LabelsMapping == nil {
		w.LabelsMapping = LabelsMapping{}
	}

	if w.Addresses == nil {
		w.Addresses = Addresses{}
	}

	if w.ChangeLabel == nil {
		_, err := w.GenerateChangeLabel()
		if err != nil {
			logging.ErrorLogger.Println(err)
			return err
		}
	}

	if len(w.UTXOMapping) == 0 {
		w.UTXOMapping = make(map[[36]byte]struct{})
		var newCollection UtxoCollection
		for _, utxo := range w.UTXOs {
			key, err := utxo.GetKey()
			if err != nil {
				logging.ErrorLogger.Println(err)
				return err
			}
			_, ok := w.UTXOMapping[key]
			if !ok {
				// if an utxo is already present remove it from the set of UTXOs
				// or actually create a new slice only with valid utxos
				// might have happened due to prior bug where duplicate utxos were added on rescan
				newCollection = append(newCollection, utxo)
			}
			w.UTXOMapping[key] = struct{}{}
		}
		w.UTXOs = newCollection
	}

	_, err := w.GenerateAddress()
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}

	return nil
}

type Address struct {
	Address string
	Comment string
}

func (w *Wallet) SortedAddresses() ([]Address, error) {
	var addresses []Address

	addresses = append(addresses)
	var nextM = 1

	for address, comment := range w.Addresses {
		if comment == StandardAddressComment {
			addresses = append(addresses, Address{
				Address: address,
				Comment: comment,
			})
			break
		}
	}

	//check:
	// todo make a goto GO-label based approach
	for nextM < len(w.LabelsMapping) {
		//var found bool

		// todo can this end up in an infinite loop
		for _, label := range w.LabelsMapping {
			if label.M == uint32(nextM) {
				addresses = append(addresses, Address{
					Address: label.Address,
					Comment: label.Comment,
				})
				nextM++
				//found = true
				//goto check
				break
			}
		}
	}

	if len(addresses) != len(w.LabelsMapping) {
		return nil, errors.New("addresses not of equal length")
	}

	return addresses, nil
}

func ConvertOwnedUTXOIntoVin(utxo *OwnedUTXO) bip352.Vin {
	vin := bip352.Vin{
		Txid:         utxo.Txid,
		Vout:         utxo.Vout,
		Amount:       utxo.Amount,
		ScriptPubKey: append([]byte{0x51, 0x20}, utxo.PubKey[:]...),
		SecretKey:    &utxo.PrivKeyTweak,
		Taproot:      true,
	}
	return vin
}
