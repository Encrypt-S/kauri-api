package wallet

import (
	"io/ioutil"
	"net/http"

	"github.com/Encrypt-S/navpi-go/app/api"
	"github.com/Encrypt-S/navpi-go/app/conf"
	"github.com/Encrypt-S/navpi-go/app/daemon/daemonrpc"
	"github.com/Encrypt-S/navpi-go/app/middleware"
	"github.com/gorilla/mux"
	"encoding/json"
	"fmt"
)

// InitTransactionHandlers sets up handlers for transaction-related rpc commands
func InitTransactionHandlers(r *mux.Router, prefix string) {

	namespace := "transactions"

	getAddressTxIdsPath := api.RouteBuilder(prefix, namespace, "v1", "getaddresstxids")
	api.ProtectedRouteHandler(getAddressTxIdsPath, r, getAddressTxIds(), http.MethodPost)

	r.Handle(getAddressTxIdsPath, middleware.Adapt(getAddressTxIds(),
		middleware.JwtHandler())).Methods("GET")

}

// GetAddressTxIdsCmd defines the "getaddresstxids" JSON-RPC command.
type GetAddressTxIdsCmd struct {
	addresses string `json:"addresses"`
}

// getAddressTxIds - executes "getaddresstxids" JSON-RPC command
// arguments - addresses array, start block height, end block height
// returns the txids for an address(es) (requires addressindex to be enabled).
func getAddressTxIds() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var getAddressTxIdsCmd GetAddressTxIdsCmd
		apiResp := api.Response{}

		err := json.NewDecoder(r.Body).Decode(&getAddressTxIdsCmd)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			returnErr := api.AppRespErrors.ServerError
			returnErr.ErrorMessage = fmt.Sprintf("Server error: %v", err)
			apiResp.Errors = append(apiResp.Errors, returnErr)
			apiResp.Send(w)
			return
		}

		n := daemonrpc.RpcRequestData{}
		n.Method = "getaddresstxids"
		n.Params = []string{getAddressTxIdsCmd.addresses}

		resp, err := daemonrpc.RequestDaemon(n, conf.NavConf)

		if err != nil { // Handle errors requesting the daemon
			daemonrpc.RpcFailed(err, w, r)
			return
		}

		bodyText, err := ioutil.ReadAll(resp.Body)
		w.WriteHeader(resp.StatusCode)
		w.Write(bodyText)

	})
}
