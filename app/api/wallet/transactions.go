package wallet

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Encrypt-S/navpi-go/app/api"
	"github.com/Encrypt-S/navpi-go/app/conf"
	"github.com/Encrypt-S/navpi-go/app/daemon/daemonrpc"
	"github.com/Encrypt-S/navpi-go/app/middleware"
	"github.com/gorilla/mux"
)

// InitTransactionHandlers sets up handlers for transaction-related rpc commands
func InitTransactionHandlers(r *mux.Router, prefix string) {

	namespace := "transactions"

	getAddressTxIdsPath := api.RouteBuilder(prefix, namespace, "v1", "getaddresstxids")
	api.ProtectedRouteHandler(getAddressTxIdsPath, r, getAddressTxIds(), http.MethodPost)

	r.Handle(getAddressTxIdsPath, middleware.Adapt(getAddressTxIds(),
		middleware.JwtHandler())).Methods("GET")

}

// getAddressTxIds - returns the txids for an address(es) (requires addressindex to be enabled).
func getAddressTxIds() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("getAddressTxIds")

		n := daemonrpc.RpcRequestData{}
		n.Method = "getaddresstxids"

		resp, err := daemonrpc.RequestDaemon(n, conf.NavConf)

		if err != nil { // Handle errors requesting the daemon
			daemonrpc.RpcFailed(err, w, r)
			return
		}

		bodyText, err := ioutil.ReadAll(resp.Body)
		w.WriteHeader(resp.StatusCode)
		w.Write(bodyText)
		io.WriteString(w, "hello world\n")
	})
}
