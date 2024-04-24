package src

import (
	"fmt"
	"github.com/setavenger/gobip352"
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

// ConvertSPRecipient converts a gobip352.Recipient to a Recipient native to this program
func ConvertSPRecipient(recipient *gobip352.Recipient) *Recipient {
	return &Recipient{
		Address:    recipient.SilentPaymentAddress,
		PkScript:   append([]byte{0x51, 0x20}, recipient.Output[:]...),
		Amount:     int64(recipient.Amount),
		Annotation: recipient.Data,
	}
}
