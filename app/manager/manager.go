package manager

import (
	"github.com/Encrypt-S/kauri-api/app/conf"
	"github.com/Encrypt-S/kauri-api/app/daemon"
)

func StartAllDaemonManagers(coins []conf.Coin)  {

	for _, coin := range coins {
		daemon.CoinStartManager(coin)
	}

}
