# Kauri API

[![Build Status](https://travis-ci.org/Encrypt-S/kauri-api.svg?branch=v1.0.0-kauri)](https://travis-ci.org/Encrypt-S/kauri-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/Encrypt-S/kauri-api)](https://goreportcard.com/report/github.com/Encrypt-S/kauri-api)
[![Coverage Status](https://coveralls.io/repos/github/Encrypt-S/kauri-api/badge.svg?branch=v1.0.0-kauri)](https://coveralls.io/github/Encrypt-S/kauri-api?branch=v1.0.0-kauri)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
![GitHub release](https://img.shields.io/github/release/Encrypt-S/kauri-api.svg)
![Last commit](https://img.shields.io/github/last-commit/Encrypt-S/kauri-api.svg)

The API for the open-source, multi-currency [Kauri Wallet](https://github.com/Encrypt-S/kauri-wallet).

## API Development
Developers working on the Kauri API can leverage the following to get up and running:

    go get -u github.com/golang/dep/cmd/dep
    go get github.com/Encrypt-S/kauri-api
    dep ensure
    go run app/main.go

This should build the app and provide you with API functionality @ 127.0.0.1:9002

## UI Development
Developers working on the [Kauri Wallet](https://github.com/Encrypt-S/kauri-wallet) can leverage the following to access the Kauri API endpoints from their local dev environment:


1. Download the proper `kauri-api` binary from [releases](https://github.com/Encrypt-S/kauri-api/releases)

2. Extract `kauri-api` release

3. Locate the dir containing extracted binaries `cd [extract-dir]`

4. Run the Go app `./kauri-api`

5. Ensure you have NAV daemon (`navcoind`) running in your Activity / Process Monitor

6. Setup Postman or something similar to test Kauri API endpoints

### Swagger Spec
https://app.swaggerhub.com/apis/Encrypt-S/kauri-api/0.0.1

### /v1/getrawtransactions
The initial endpoint can be tested in as follows:

#### POST
http://127.0.0.1:9002/api/transactions/v1/getrawtransactions

#### AUTH
username: rpcuser
password: rpcpassword

#### HEADERS
`Content-Type: application/json`

#### BODY
This is the structure of the raw request body to be used in the POST:
```
{"transactions": [
    {"currency":  "NAV", "addresses": ["validNAVaddress1", "validNAVaddress2"]}
]}
```

#### MODELS
The request body above can be organised into the following models:

  **transactions** - Addresses for each currency in wallet

    Transactions {
      transactions  [...]
    }

  **WalletItem** - Object containing currency and array of addresses

    WalletItem {
      currency  string
      addresses WalletAddresses[...]
    }

  **WalletItem** - Object containing currency and array of addresses

    WalletItem {
      currency  string
      addresses WalletAddresses[...]
    }


  **WalletAddresses** - Array of addresses

    WalletAddresses [string, string, string]












