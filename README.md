# go-iden3-core

Go implementation of the iden3 system.

[![Go Report Card](https://goreportcard.com/badge/github.com/iden3/go-iden3-core)](https://goreportcard.com/report/github.com/iden3/go-iden3-core)
[![Build Status](https://travis-ci.org/iden3/go-iden3-core.svg?branch=master)](https://travis-ci.org/iden3/go-iden3-core)

## Install
```
$ go get github.com/iden3/go-iden3-core
```

## Documentation

Go Modules documentation:
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/testgen?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/testgen) testgen
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/merkletree?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/merkletree) merkletree
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/keystore?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/keystore) keystore
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/identity?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/identity) identity
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/identity/issuer?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/identity/issuer) identity/issuer
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/identity/holder?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/identity/holder) identity/holder
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/eth?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/eth) eth
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/eth/contracts?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/eth/contracts) eth/contracts
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/core?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/core) core
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/core/claims?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/core/claims) core/claims
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/core/genesis?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/core/genesis) core/genesis
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/core/proof?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/core/proof) core/proof
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/common?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/common) common
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/utils?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/utils) utils
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/utils/noncedb?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/utils/noncedb) utils/noncedb
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/db?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/db) db
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/components?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/components) components
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/components/httpclient?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/components/httpclient) components/httpclient
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/components/idenpubonchain?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/components/idenpubonchain) components/idenpubonchain
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/components/idenpubonchain/mock?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/components/idenpubonchain/mock) components/idenpubonchain/mock
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/components/verifier?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/components/verifier) components/verifier
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/components/idensigner?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/components/idensigner) components/idensigner
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/components/idenpuboffchain?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/components/idenpuboffchain) components/idenpuboffchain
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/components/idenpuboffchain/readermock?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/components/idenpuboffchain/readermock) components/idenpuboffchain/readermock
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/components/idenpuboffchain/readerhttp?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/components/idenpuboffchain/readerhttp) components/idenpuboffchain/readerhttp
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/components/idenpuboffchain/writerhttp?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/components/idenpuboffchain/writerhttp) components/idenpuboffchain/writerhttp
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/components/idenpuboffchain/writermock?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/components/idenpuboffchain/writermock) components/idenpuboffchain/writermock
- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/crypto?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/crypto) crypto

## Testing
`go test ./...`



### WARNING
All code here is experimental and WIP

## License
go-iden3-core is part of the iden3 project copyright 2018 0kims association and published with GPL-3 license, please check the LICENSE file for more details.
