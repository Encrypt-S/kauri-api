# Developer Workflow
Notes, config, and general workflow for Kauri API development

## clone [kauri-api](https://github.com/Encrypt-S/kauri-api.git)
This is the main project repo for the Kowhai API

    git clone https://github.com/Encrypt-S/kauri-api.git

## set project WORKDIR
    cd /go/src/github.com/Encrypt-S/kauri-api/app

## install main dependencies
    go get ./

## install test dependencies
    go get -t -v ./...

## run app
    go run main.go

## run tests
    go test ./...

## clone [nav-docker](https://github.com/Encrypt-S/nav-docker)
This repo contains Docker files used to build and run containerized instances of the different Nav projects.

    git clone https://github.com/Encrypt-S/nav-docker.git

## docker-navcoind
View the info at [docker-navcoind](https://github.com/NAVCoin/nav-docker/blob/master/docker-navcoind/README.md) to run
navcoind inside a Docker container. Make sure you have the docker image building properly and that
`docker-compose up` successfully starts the service and runs the `navcoind` daemon.

## navcoin-cli access
After you run `docker-compuse up` open up a new terminal tab/window and run:

    docker exec -it dockernavcoind_testnet_1 /bin/bash

you should now be in cli mode with something like this:

    root@795b5c0525c0:/#

you will now be able to execute rpc commands accordingly

## testing endpoints
- Once the app is running it will be accessible at `127.0.0.1:9002`
- The default `managerApiPort` is `9002` and set in `server-config.json`

### getrawtransactions
returns raw transactions from supplied array of wallet addresses

    /api/transactions/v1/getrawtransactions

## daemon API
These endpoints are used in the NavPi UI.

### getstakereport
lists last single 30 day stake subtotal and last 24h, 7, 30, 365 day subtotal.

    /api/wallet/v1/getstakereport

### encryptwallet
Encrypts the wallet with _passphrase_. Once encrypted, three new commands are available:
_walletlock, walletpassphrase, walletpassphrasechange_

    /api/wallet/v1/encryptwallet





