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
	"github.com/Encrypt-S/kauri-api/app/manager"
)

func main() {

	// prime the app
	initMain()

	// log out server runtime OS and Architecture
	log.Println(fmt.Sprintf("Server running in %s:%s", runtime.GOOS, runtime.GOARCH))
	log.Println(fmt.Sprintf("App pid : %d.", os.Getpid()))

	// load the server config - required - contains server data
	serverConfig, err := conf.LoadServerConfig()
	if err != nil {
		log.Fatal("Failed to load the server config: " + err.Error())
	}

	// load the app config - required - contains active coin data
	err = conf.LoadAppConfig()
	if err != nil {
		log.Println("Failed to load the app config: " + err.Error())
	}

	conf.StartConfigManager()

	// load the dev config file if one is set
	//conf.LoadDevConfig()

	// start the daemon managers for active coins
	manager.StartAllDaemonManagers(conf.AppConf.ActiveCoins)

	// setup the router
	router := mux.NewRouter()

	// setup the api meta and coin meta handlers
	api.InitMetaHandlers(router, "api")

	// start the transaction handlers for active coins
	manager.StartAllTxHandlers(router, conf.AppConf.ActiveCoins)

	// set the proper server port
	port := fmt.Sprintf(":%d", serverConfig.ManagerAPIPort)

	// start http server and listen up
	http.ListenAndServe(port, router)
}

// initMain primes the pump before main init sequence
func initMain() {
	api.BuildAppErrors()
	conf.CreateRPCDetails()
}
