package src

import "errors"

// todo group by logical context
// errors with `should not happen` indicate a bug in the program and not necessarily user input errors
var (
	ErrLabelAlreadyExists = errors.New("label already exists")

	ErrMinChangeAmountNotReached = errors.New("min change amount not reached")

	ErrWrongLengthRecipients = errors.New("wrong length new recipients")

	ErrRecipientIncomplete = errors.New("recipient incomplete")

	ErrNotSPAddress = errors.New("not a sp address")

	ErrNoUTXOsInWallet = errors.New("no utxos in wallet")

	ErrInsufficientFunds = errors.New("insufficient funds")

	ErrNoMatchForUTXO = errors.New("could not match UTXO to foundOutput, should not happen")

	ErrTxInputAndVinLengthMismatch = errors.New("tx inputs and vins have different length, should not happen")

	ErrNoMatchingVinFoundForTxInput = errors.New("there was no matching vin for the given tx input, should not happen")

	ErrDaemonNotSet = errors.New("daemon is not initialised")

	ErrDaemonIsLocked = errors.New("daemon is locked")

	ErrInvalidFeeRate = errors.New("invalid fee rate")

	ErrRecipientAmountIsZero = errors.New("recipient amount is zero")
)
