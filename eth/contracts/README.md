# State smart contracts

```
abigen --sol State.sol --pkg=contracts --out state.go
```

Generated from git@github.com:iden3/contracts.git commit cff032a968254133195a6cf28ae27ff928ccf28b with `pragma solidity ^0.6.0;`.

# Poseidon raw smart contract

```
cd js
npm ci
./poseidon_gencontract.js
cd ..
abigen --bin=js/poseidon.bin --abi=js/poseidon.abi --pkg=contracts --type PoseidonCircomlib --out=poseidon.go
```
