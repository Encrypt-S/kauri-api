package main

import (
	"fmt"
	"log"
	"net/http"

	"os"
	"runtime"

	"github.com/Encrypt-S/kauri-api/app/api"
	"github.com/Encrypt-S/kauri-api/app/conf"
	"github.com/gorilla/mux"
	"github.com/Encrypt-S/kauri-api/app/api/wallet"
)

func main() {

	initMain()

	// log out server runtime OS and Architecture
	log.Println(fmt.Sprintf("Server running in %s:%s", runtime.GOOS, runtime.GOARCH))
	log.Println(fmt.Sprintf("App pid : %d.", os.Getpid()))

	// load the server config - this is required otherwise we die right here
	serverConfig, err := conf.LoadServerConfig()
	if err != nil {
		log.Fatal("Failed to load the server config: " + err.Error())
	}

	// Load the App config
	err = conf.LoadAppConfig()
	if err != nil {
		log.Println("Failed to load the app config: " + err.Error())
	}

	conf.StartConfigManager()

	//load the dev config file if one is set
	conf.LoadDevConfig()

	// setup the router and the api
	router := mux.NewRouter()
	api.InitMetaHandlers(router, "api")

	//daemon.StartManager()

	/*
		// check to see if we have a defined running config
		// If not we are only going to boot the setup apis, otherwise we will start the app
		if conf.AppConf.RunningNavVersion == "" {

			log.Println("No App Config starting the setup api")
			setupapi.InitSetupHandlers(router, "api")

		} else {

			log.Println("App config found :: booting all apis!")
			// we have a user config so start the app in running mode
			// TODO: make dependent on the dev config
			daemon.StartManager()

			// stat all app API's
			managerapi.InitManagerhandlers(router, "api")
			daemonapi.InitChainHandlers(router, "api")
			daemonapi.InitWalletHandlers(router, "api")
			user.InitSetupHandlers(router, "api")

		}
	*/

	// init the transaction handlers
	wallet.InitTransactionHandlers(router, "api")

	// Start http server
	port := fmt.Sprintf(":%d", serverConfig.ManagerAPIPort)
	http.ListenAndServe(port, router)
}

// Start everything before we get going
func initMain() {

	api.BuildAppErrors()
	conf.CreateRPCDetails()


}
