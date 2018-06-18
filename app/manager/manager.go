package manager

import (
	"github.com/Encrypt-S/kauri-api/app/conf"
	"github.com/Encrypt-S/kauri-api/app/daemon"
	"github.com/Encrypt-S/kauri-api/app/daemon/daemonapi"
	"github.com/gorilla/mux"
	"log"
)

// StartAllDaemonManagers ranges through coins, starts daemons
func StartAllDaemonManagers(activeCoins []conf.CoinData)  {

	log.Println("ranging through active coins, starting daemons")

	for _, coinData := range activeCoins {
		daemon.StartManager(coinData)
	}

}

// StartWalletHandlers ranges through activeCoins, inits handlers
func StartWalletHandlers(r *mux.Router, activeCoins []conf.CoinData)  {

	log.Println("ranging through active coins, initialising wallet handlers")

	for _, coinData := range activeCoins {
		daemonapi.InitWalletHandlers(r, coinData, "api")
	}

}
