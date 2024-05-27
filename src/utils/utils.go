package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/setavenger/blindbitd/src/logging"
	"github.com/setavenger/go-bip352"
)

func ConvertToFixedLength34(input []byte) [34]byte {
	if len(input) != 34 {
		panic(fmt.Sprintf("wrong length expected 32 got %d", len(input)))
	}
	var output [34]byte
	copy(output[:], input)
	return output
}

// IsSilentPaymentAddress determines whether an address is a silent payment address.
// Works only for silent payment v0
func IsSilentPaymentAddress(address string) bool {
	// only works for v1
	if len(address) == 116 && address[:2] == "sp" {
		return true
	}
	if len(address) == 117 && address[:3] == "tsp" {
		return true
	}
	return false
}

func CopyBytes(bytes []byte) []byte {
	result := make([]byte, len(bytes))
	copy(result, bytes)
	return result
}

// ConvertPubKeyToScriptHash
// Converts the given taproot pubKey to a scriptHash which can be checked with electrumX
func ConvertPubKeyToScriptHash(pubKey [32]byte) string {
	data := append([]byte{0x51, 0x20}, pubKey[:]...)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(bip352.ReverseBytesCopy(hash[:]))
}

// SerialiseVinToOutpoint serialises a vin to an outpoint in LE encoding
func SerialiseVinToOutpoint(vin bip352.Vin) ([36]byte, error) {
	var buf bytes.Buffer
	buf.Write(bip352.ReverseBytesCopy(vin.Txid[:]))
	err := binary.Write(&buf, binary.LittleEndian, vin.Vout)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return [36]byte{}, err
	}
	var outpoint [36]byte
	copy(outpoint[:], buf.Bytes())
	return outpoint, nil
}
