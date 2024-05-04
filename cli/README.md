# BlindBit-cli

This cli application controls the [blindbit daemon (blindbitd)](https://github.com/setavenger/BlindBitd). Send and
receive are the first features to come.
The daemon already has some more features which are not exposed through this cli. Those can be accessed using gRPC tools
like grpcui (unix socket is exposed via `socat`). The cli will be expanded as this project moves forward.
The daemon will also be expanded, stay tuned...

## Commands Todo

### Priority 1

- [x] Status
- [x] SyncHeight
- [x] Unlock
- [x] Shutdown
- [x] ListUTXOs
    - [x] Show balance filtered by state
    - [x] List utxos by state
- [x] ListAddresses
- [ ] CreateTransaction
  - [ ] single recipient
  - [ ] multi recipients
- [ ] BroadcastRawTx
- [ ] GetMnemonic
- [x] CreateNewWallet
- [ ] ForceRescanFromHeight

### Priority 2

- [ ] Overview (holistic wrapper to show key information in one view)

## General Todos

- [ ] take care of modules deprecation warnings
- [ ] make outputs pretty and easy to understand
    - Format outputs in tables to make data more clear and appear more structured
    - add coloring?