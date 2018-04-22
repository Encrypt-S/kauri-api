package wallet

import (
	"net/http"

	"encoding/json"
	"fmt"

	"github.com/Encrypt-S/navpi-go/app/api"
	"github.com/gorilla/mux"
)

// InitTransactionHandlers sets up handlers for transaction-related rpc commands
func InitTransactionHandlers(r *mux.Router, prefix string) {

	namespace := "transactions"

	getAddressTxIdsPath := api.RouteBuilder(prefix, namespace, "v1", "getaddresstxids")
	api.OpenRouteHandler(getAddressTxIdsPath, r, getAddressTxIds())

}

// structure of POST body to API
//{"transactions": [
//{"currency":  "NAV", "addresses": ["Nkjhsdfkjh834jdu", "Nisd8a8BAhahs"]},
//{"currency":  "BTC",  "addresses": ["Bak7ahbZAA", "B91janABsa"]}
//]}

// GetAddressTxIdsArray first decode transactions json into Transactions slice
type GetAddressTxIdsArray struct {
	Transactions []GetAddressTxIdsJSON `json:"transactions"`
}

// GetAddressTxIdsJSON represents the keys Transactions slice
type GetAddressTxIdsJSON struct {
	Currency string `json:"currency"`
	Addresses  []string `json:"addresses"`
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

		// TODO: create new function to loop over addresses array for each currency

		// then get transactions for each
		//n := daemonrpc.RpcRequestData{}
		//n.Method = "getaddresstxids"
		//n.Params = getAddressTxIds.Transactions

		//resp, err := daemonrpc.RequestDaemon(n, conf.NavConf)

		//if err != nil { // Handle errors requesting the daemon
		//	daemonrpc.RpcFailed(err, w, r)
		//	return
		//}

		// then reassemble data
		// then return formatted response

		//bodyText, err := ioutil.ReadAll(resp.Body)
		//w.WriteHeader(resp.StatusCode)
		//w.Write(bodyText)

		// write test

	})
}
