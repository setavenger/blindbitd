# BlindBitd (not functional yet)

Receive and simple send functions but not integrated and not functional yet.
This is the daemon for the BlindBit Wallet.

## Todo

### Priority 1

- [ ] Bring test coverage to a meaningful level
- [x] Create a coin selector that incorporates the fees
- [x] Binary encoding to reduce bandwidth save time on decoding (protoBuffs)
- [ ] Add Transaction History
- [x] Mark UTXOs as spent (or similar) if used for a transaction
- [ ] Sometimes unlock does not work on first try, needs a restart of the daemon
- [ ] Add gRPC credentials

### Priority 2

- [ ] More tests for coin selector
  - Selector seems very accurate, but should rather do +1 sat to exceed fee and don't go below
- [ ] Coin selector allow float fees
- [ ] UTXO export - similar to a backup to avoid rescanning from birthHeight
- [ ] Separate spending password
- [ ] Out-of-band notifications
    - share tweak and tx data directly with the receiver to reduce scanning efforts (follow blindbit standard set for
      the mobile app)
- [ ] Balance checks for UTXOs: account for more than one UTXO per script
- [ ] Expand logging especially on errors
- [ ] Check which panics to keep
- [ ] Automatically make annotation in tx-history if sent to sp-address, not possible to reconstruct in hindsight
- [ ] Don't always add change in coin selector (see todo)

## IPC

### Todos

- [x] Create New label
    - Give comment
    - Returns label address
- [x] Create Tx and broadcast
- [x] Broadcast raw Tx
- [ ] Pause/Resume scanning
- [ ] List UTXOs by label
