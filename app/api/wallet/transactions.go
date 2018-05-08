package wallet

import (
	"net/http"

	"encoding/json"
	"fmt"

	"io/ioutil"

	"github.com/Encrypt-S/kauri-api/app/api"
	"github.com/Encrypt-S/kauri-api/app/conf"
	"github.com/Encrypt-S/kauri-api/app/daemon/daemonrpc"
	"github.com/gorilla/mux"
	"github.com/Encrypt-S/kauri-api/app/daemon"
	"github.com/Encrypt-S/kauri-api/app/navdaemon"
)

// InitTransactionHandlers sets up handlers for transaction-related rpc commands
func InitTransactionHandlers(r *mux.Router, prefix string) {

	namespace := "transactions"

	// get raw transactions endpoint :: provides raw transaction data for supplied wallet addresses
	getRawTransactionsPath := api.RouteBuilder(prefix, namespace, "v1", "getrawtransactions")
	api.OpenRouteHandler(getRawTransactionsPath, r, getRawTransactionsHandler())

}




// getRawTransactionsHandler ranges through transactions, returns RPC response data
func getRawTransactionsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiResp := api.Response{}

		var incomingTxs daemon.IncomingTransactions

		err := json.NewDecoder(r.Body).Decode(&incomingTxs)

		if err != nil {
			returnErr := api.AppRespErrors.JSONDecodeError
			returnErr.ErrorMessage = fmt.Sprintf("JSON decode error: %v", err)
			apiResp.Errors = append(apiResp.Errors, returnErr)
			apiResp.Send(w)
			return
		}

		res := []daemon.Result{}

		//resp, err :=

			/*
			buildResponse(incomingTxs)

		if err != nil {
			returnErr := api.AppRespErrors.RPCResponseError
			returnErr.ErrorMessage = fmt.Sprintf("Response error: %v", err)
			apiResp.Errors = append(apiResp.Errors, returnErr)
			apiResp.Send(w)
			return
		}

			*/

		apiResp.Data = resp.Results

		apiResp.Send(w)

		return
	})
}


