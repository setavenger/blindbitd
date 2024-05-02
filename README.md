# BlindBitd (not functional yet)

Receive and simple send functions but not integrated and not functional yet.
This is the daemon for the BlindBit Wallet.

## Todo

### Priority 1

- [ ] Create a coin selector that incorporates the fees
- [ ] Binary encoding to reduce bandwidth save time on decoding (protoBuffs)
- [ ] Add Transaction History
- [x] Mark UTXOs as spent (or similar) if used for a transaction

### Priority 2

- [ ] UTXO export - similar to a backup to avoid rescanning from birthHeight
- [ ] Separate spending password
- [ ] Out-of-band communications
    - share tweak and tx data directly with the receiver to reduce scanning efforts (follow blindbit standard set for
      the mobile app)
- [ ] Balance checks for UTXOs: account for more than one UTXO per script

## IPC

### Todos

- [x] Create New label
    - Give comment
    - Returns label address
- [ ] Create Tx and broadcast
- [ ] Broadcast raw Tx
- [ ] Pause/Resume scanning
- [ ] List UTXOs by label
- [ ] Broadcast transaction
