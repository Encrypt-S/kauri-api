package manager

import (
	"github.com/Encrypt-S/kauri-api/app/conf"
	"github.com/Encrypt-S/kauri-api/app/daemon"
	"github.com/Encrypt-S/kauri-api/app/daemon/daemonapi"
	"github.com/gorilla/mux"
)

// StartAllDaemonManagers ranges through coins, starts daemons
func StartAllDaemonManagers(activeCoins []conf.CoinData)  {

	for _, coin := range activeCoins {
		daemon.StartDaemonManager(coin)
	}

}

// StartAllTxHandlers ranges through activeCoins, inits handlers
func StartAllTxHandlers(r *mux.Router, activeCoins []conf.CoinData)  {

	for _, coin := range activeCoins {
		daemonapi.InitTxHandlers(r, coin, "api")
	}

}
