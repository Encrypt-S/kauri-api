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
	"github.com/Encrypt-S/kauri-api/app/daemon/daemonapi"
	"github.com/Encrypt-S/kauri-api/app/manager"
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
	//conf.LoadDevConfig()

	//daemon.StartManager()
	manager.StartAllDaemonManagers(conf.AppConf.Coins)


	// setup the router
	router := mux.NewRouter()

	// setup the api
	api.InitMetaHandlers(router, "api")

	// init the transaction handlers
	daemonapi.InitTxHandlers(router, "api")

	// Start http server
	port := fmt.Sprintf(":%d", serverConfig.ManagerAPIPort)
	http.ListenAndServe(port, router)
}

// Start everything before we get going
func initMain() {

	api.BuildAppErrors()
	conf.CreateRPCDetails()

}
