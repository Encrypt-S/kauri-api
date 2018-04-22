package wallet

import (
	"io/ioutil"
	"net/http"

	"encoding/json"
	"fmt"

	"github.com/Encrypt-S/navpi-go/app/api"
	"github.com/Encrypt-S/navpi-go/app/conf"
	"github.com/Encrypt-S/navpi-go/app/daemon/daemonrpc"
	"github.com/gorilla/mux"
)

// InitTransactionHandlers sets up handlers for transaction-related rpc commands
func InitTransactionHandlers(r *mux.Router, prefix string) {

	namespace := "transactions"

	getAddressTxIdsPath := api.RouteBuilder(prefix, namespace, "v1", "getaddresstxids")
	api.OpenRouteHandler(getAddressTxIdsPath, r, getAddressTxIds())

}

// format of transactions json payload from incoming POST
//{"transactions": [
//{"currency":  "NAV", "address": "Nkjhsdfkjh834jdu"},
//{"currency":  "NAV", "address": "Nkjhsdfkjh834jdu"},
//{"currency":  "BTC", "address": "1kjhsdfkjh834jdu"},
//{"currency":  "BTC", "address": "Nkjhsdfkjh834jdu"}
//]}

// TODO: Decode top level transactions JSON array into a slice of structs

// GetAddressTxIdsArray first decode transactions json into Transactions slice
type GetAddressTxIdsArray struct {
	Transactions []GetAddressTxIdsJSON `json:"array"`
}

// GetAddressTxIdsJSON represents the keys Transactions slice
type GetAddressTxIdsJSON struct {
	Currency string `json:"currency"`
	Address  string `json:"address"`
}

// getAddressTxIds - executes "getaddresstxids" JSON-RPC command
func getAddressTxIds() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var getAddressTxIds GetAddressTxIdsArray
		apiResp := api.Response{}

		err := json.NewDecoder(r.Body).Decode(&getAddressTxIds)

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
		n.Params = getAddressTxIds.Transactions

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
