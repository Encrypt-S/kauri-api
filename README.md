# Kauri API

[![Build Status](https://travis-ci.org/Encrypt-S/kauri-api.svg?branch=v1.0.0-kauri)](https://travis-ci.org/Encrypt-S/kauri-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/Encrypt-S/kauri-api)](https://goreportcard.com/report/github.com/Encrypt-S/kauri-api)
[![Coverage Status](https://coveralls.io/repos/github/Encrypt-S/kauri-api/badge.svg?branch=v1.0.0-kauri)](https://coveralls.io/github/Encrypt-S/kauri-api?branch=v1.0.0-kauri)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
![GitHub release](https://img.shields.io/github/release/Encrypt-S/kauri-api.svg)
![Last commit](https://img.shields.io/github/last-commit/Encrypt-S/kauri-api.svg)

The API for the open-source, multi-currency [Kauri Wallet](https://navhub.org/projects/kauri/)

## API Development

Installation for Developers working on the Kauri API

    go get -u github.com/golang/dep/cmd/dep
    go get github.com/Encrypt-S/kauri-api
    dep ensure
    go run app/main.go

This should build the app and provide you with API functionality @ 127.0.0.1:9002

## UI Development

Installation for Developers working on the [Kauri Wallet](https://github.com/Encrypt-S/kauri-wallet)

1. Download the latest `kauri-api` binary from [releases](https://github.com/Encrypt-S/kauri-api/releases)

2. Extract `kauri-api` release

3. CD into the dir containing extracted binaries `cd [extract-dir]`

4. Run the app `./kauri-api`

5. Ensure you have NAV daemon (`navcoind`) running in your Activity / Process Monitor

6. Setup Postman or something similar to test Kauri API endpoints

### Swagger API Spec

https://app.swaggerhub.com/apis/Encrypt-S/kauri-api/0.1.0

### API Calls

The initial endpoint can be tested in Postman, Paw, Shell, Angular app, etc...

#### POST to /v1/getrawtransactions

    http://127.0.0.1:9002/api/transactions/v1/getrawtransactions

#### Auth

    username: rpcuser
    password: rpcpassword

#### Headers
`Content-Type: application/json`

#### Body
This is the structure of the raw request body to be used in the POST:

    {"transactions": [
        {"currency":  "NAV", "addresses": ["NW7uXr4ZAeJKigMGnKbSLfCBQY59cH1T8G", "NUDke42E3fwLqaBbBFRyVSTETuhWAi7ugk"]}
    ]}

#### Models

The request body has the following model structure:

  **transactions** - Addresses for each currency in wallet

    Transactions {
      transactions  [...]
    }

  **WalletItem** - Object containing currency and array of addresses

    WalletItem {
      currency  string
      addresses WalletAddresses[...]
    }

  **WalletAddresses** - Array of addresses

    WalletAddresses [string, string, string]

#### Response `200 OK`

A successful response will contain raw transaction data for supplied wallet addresses.

    {
      "data": [
        {
          "currency": "NAV",
          "addresses": [
            {
              "address": "NW7uXr4ZAeJKigMGnKbSLfCBQY59cH1T8G",
              "transactions": [
                  {
                    "txid": "11a7071a43a8da2b9ac116865a6cd92c985c3f7cbde63933d253f88dffaa311a",
                    "rawtx": "",
                    "verbose": null
                  },
                  {
                    "txid": "c8dad515d5e5c7a45bc5b3814fcf5e1f63474c9b67f84ee2ab9803f809e94929",
                    "rawtx": "01000000f0f33457011a31aaff8df853d23339e6bd7c3f5c982cd96c5a8616c19a2bdaa8431a07a711010000006a47304402202fbb2c5955013fc4806420a66e5c9116902c0263fe7920ae104ff1818ef62efd022040857e3108ae8f30e8a0800f8f892c8a97aa88b67b8e40032e2ba33d3445230e012103f6c3b8154a19327783dd46e0dda13f812f57b00f9246387f62d5ece8bed767b4ffffffff0300000000000000000000debdfcc1c60100232103f6c3b8154a19327783dd46e0dda13f812f57b00f9246387f62d5ece8bed767b4ac3688d6fcc1c60100232103f6c3b8154a19327783dd46e0dda13f812f57b00f9246387f62d5ece8bed767b4ac00000000",
                    "verbose": {
                        "anon-destination": "",
                        "blockhash": "52260690630225abb5b9bd1f9b72774ced5f9b74e18ac2ab7dd5b76d229fbfdd",
                        "blocktime": 1463088112,
                        "confirmations": 55769,
                        "hash": "c8dad515d5e5c7a45bc5b3814fcf5e1f63474c9b67f84ee2ab9803f809e94929",
                        "height": 523,
                        "hex": "01000000f0f33457011a31aaff8df853d23339e6bd7c3f5c982cd96c5a8616c19a2bdaa8431a07a711010000006a47304402202fbb2c5955013fc4806420a66e5c9116902c0263fe7920ae104ff1818ef62efd022040857e3108ae8f30e8a0800f8f892c8a97aa88b67b8e40032e2ba33d3445230e012103f6c3b8154a19327783dd46e0dda13f812f57b00f9246387f62d5ece8bed767b4ffffffff0300000000000000000000debdfcc1c60100232103f6c3b8154a19327783dd46e0dda13f812f57b00f9246387f62d5ece8bed767b4ac3688d6fcc1c60100232103f6c3b8154a19327783dd46e0dda13f812f57b00f9246387f62d5ece8bed767b4ac00000000",
                        "locktime": 0,
                        "size": 258,
                        "time": 1463088112,
                        "txid": "c8dad515d5e5c7a45bc5b3814fcf5e1f63474c9b67f84ee2ab9803f809e94929",
                        "version": 1,
                        "vin": [
                            {
                              "scriptSig": {
                                  "asm": "304402202fbb2c5955013fc4806420a66e5c9116902c0263fe7920ae104ff1818ef62efd022040857e3108ae8f30e8a0800f8f892c8a97aa88b67b8e40032e2ba33d3445230e[ALL] 03f6c3b8154a19327783dd46e0dda13f812f57b00f9246387f62d5ece8bed767b4",
                                  "hex": "47304402202fbb2c5955013fc4806420a66e5c9116902c0263fe7920ae104ff1818ef62efd022040857e3108ae8f30e8a0800f8f892c8a97aa88b67b8e40032e2ba33d3445230e012103f6c3b8154a19327783dd46e0dda13f812f57b00f9246387f62d5ece8bed767b4"
                              },
                              "sequence": 4294967295,
                              "txid": "11a7071a43a8da2b9ac116865a6cd92c985c3f7cbde63933d253f88dffaa311a",
                              "vout": 1
                            }
                        ],
                        "vout": [
                            {
                              "n": 0,
                              "scriptPubKey": {
                                  "asm": "",
                                  "hex": "",
                                  "type": "nonstandard"
                              },
                              "value": 0,
                              "valueSat": 0
                            },
                            {
                              "n": 1,
                              "scriptPubKey": {
                                  "addresses": [
                                      "NW7uXr4ZAeJKigMGnKbSLfCBQY59cH1T8G"
                                  ],
                                  "asm": "03f6c3b8154a19327783dd46e0dda13f812f57b00f9246387f62d5ece8bed767b4 OP_CHECKSIG",
                                  "hex": "2103f6c3b8154a19327783dd46e0dda13f812f57b00f9246387f62d5ece8bed767b4ac",
                                  "reqSigs": 1,
                                  "type": "pubkey"
                              },
                              "value": 5000114.48,
                              "valueSat": 500011448000000
                            },
                            {
                              "n": 2,
                              "scriptPubKey": {
                                  "addresses": [
                                      "NW7uXr4ZAeJKigMGnKbSLfCBQY59cH1T8G"
                                  ],
                                  "asm": "03f6c3b8154a19327783dd46e0dda13f812f57b00f9246387f62d5ece8bed767b4 OP_CHECKSIG",
                                  "hex": "2103f6c3b8154a19327783dd46e0dda13f812f57b00f9246387f62d5ece8bed767b4ac",
                                  "reqSigs": 1,
                                  "type": "pubkey"
                              },
                              "value": 5000114.49616438,
                              "valueSat": 500011449616438
                            }
                        ],
                        "vsize": 258
                      }
                  }
              ]
            }
          ]
        }
      ]
    }













