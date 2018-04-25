package wallet

import (
	"net/http"

	"encoding/json"
	"fmt"

	"github.com/Encrypt-S/kauri-api/app/api"
	"github.com/Encrypt-S/kauri-api/app/conf"
	"github.com/Encrypt-S/kauri-api/app/daemon/daemonrpc"
	"github.com/gorilla/mux"
)

// InitTransactionHandlers sets up handlers for transaction-related rpc commands
func InitTransactionHandlers(r *mux.Router, prefix string) {

	namespace := "transactions"

	// setup endpoint to be used for receiving txids for supplied addresses
	getAddressTxIdsPath := api.RouteBuilder(prefix, namespace, "v1", "getaddresstxids")
	api.OpenRouteHandler(getAddressTxIdsPath, r, getAddressTxIds())

}

// IncomingTransactionsArray describes the Transactions array
type IncomingTransactionsArray struct {
	Transactions []IncomingTransactions `json:"transactions"`
}

// IncomingTransactions describes the incoming Transactions array
type IncomingTransactions struct {
	Currency  string   `json:"currency"`
	Addresses []string `json:"addresses"`
}

// OutgoingTransactionsArray describes the parsed transactions data array
type OutgoingTransactionsArray struct {
	Transactions []OutgoingTransactions `json:"transactions"`
}

// OutgoingTransactions describes the outgoing response
type OutgoingTransactions struct {
	Currency          string   `json:"currency"`
	OutgoingAddressObject []interface{} `json:"addressobject"`
}

// OutgoingAddressObject contains address and array of txids
type OutgoingAddressObject struct {
	Address string `json:"address"`
	OutgoingTxIdsArray []string `json:"txids"`
}

// OutgoingTxIdsArray describes the outgoing addresses array
type OutgoingTxIdsArray struct {
	TxIds []string `json:"txids"`
}

// GetAddressTxIdRPCParams describes addresses array params for getaddresstxids call
type GetAddressTxIdRPCParams struct {
	Addresses []string `json:"addresses"`
}

// GetAddressTxIdRPCResponse contains RPC response
type GetAddressTxIdRPCResponse struct {
	Result []string `json:"result"`
}

// getAddressTxIds - ranges through transactions, returns raw transactions
func getAddressTxIds() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var incomingTxs IncomingTransactionsArray
		apiResp := api.Response{}

		err := json.NewDecoder(r.Body).Decode(&incomingTxs)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			returnErr := api.AppRespErrors.ServerError
			returnErr.ErrorMessage = fmt.Sprintf("Server error: %v", err)
			apiResp.Errors = append(apiResp.Errors, returnErr)
			apiResp.Send(w)
			return
		}

		// start building outgoing transactions
		outgoingTxs := OutgoingTransactions{}

		// start building outgoing addresses array
		outgoingAddressObject := OutgoingAddressObject{}

		// start building outgoing txid array
		//outgoingTxIds := OutgoingTxIdsArray{}

		// range over the Transactions
		for _, tx := range incomingTxs.Transactions {

			// isolate NAV addresses
			if tx.Currency == "NAV" {

				// set currency for this iteration
				outgoingTxs.Currency = "NAV"

				// range over the NAV addresses
				for _, address := range tx.Addresses {

					// add NAV address to address object
					outgoingAddressObject.Address = address

					// payload should contain array of address objects
					// with array of transactions in each object

					// bring in rpc get address txids struct to setup rpc params
					// this is done each iteration to create fresh address params
					// we only want to send one address in rpc call (this iteration)
					getParams := GetAddressTxIdRPCParams{}

					// add current address to addresses array for rpc params
					getParams.Addresses = append(getParams.Addresses, address)

					// prepare rpc call
					n := daemonrpc.RpcRequestData{}
					// set method
					n.Method = "getaddresstxids"
					// set params
					n.Params = []GetAddressTxIdRPCParams{getParams}

					// issue rpc call
					resp, rpcErr := daemonrpc.RequestDaemon(n, conf.NavConf)

					// handle error with rpc call
					if rpcErr != nil {
						daemonrpc.RpcFailed(rpcErr, w, r)
						return
					}

					// handle the nav daemon response
					txidResp := GetAddressTxIdRPCResponse{}

					// get the json from the response Body
					jsonErr := json.NewDecoder(resp.Body).Decode(&txidResp)

					// handle error decoding json
					if jsonErr != nil {
						w.WriteHeader(http.StatusInternalServerError)
						returnErr := api.AppRespErrors.ServerError
						returnErr.ErrorMessage = fmt.Sprintf("Server error: %v", jsonErr)
						apiResp.Errors = append(apiResp.Errors, returnErr)
						apiResp.Send(w)
						return
					}

					//println(txidResp)

					// add the decoded txid response to the outgoing TxIds array
					//outgoingTxIds.TxIds = append(outgoingTxIds.TxIds, txidResp)

				}

			}
		}

		return

	})
}
