package wallet

import (
	"log"
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

// OutgoingTransactions describes individual address contents
type OutgoingTransactions struct {
	Currency     string        `json:"currency"`
	Address      string        `json:"address"`
	Transactions []interface{} `json:"transactions"`
}

// RPCGetAddressTxIdsParams describes params needed for getaddresstxids daemonrpc call
type RPCGetAddressTxIdsParams struct {
	Addresses []string `json:"addresses"`
}

// RPCGetAddressTxIdsResponse describes returned address and txids
type RPCGetAddressTxIdsResponse struct {
	Address string
	TxIds   []string
}

// RPCTxIdsArray describes array of txids returned from rpc call
type RPCTxIdsArray struct {
	TxIds []string `json:"addresses"`
}

// RPCTxLookupArray describes the array of txids to lookup
type RPCTxLookupArray struct {
	TxIdsToLookup []RPCGetAddressTxIdsResponse
}

// getAddressTxIds - ranges through transactions, returns raw transactions
func getAddressTxIds() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var getAddressTxIds IncomingTransactionsArray
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

		// range over the Transactions
		for _, tx := range getAddressTxIds.Transactions {

			// isolate NAV addresses
			if tx.Currency == "NAV" {

				// range over the NAV addresses
				for _, address := range tx.Addresses {

					// add each address to response data struct for future use
					navData := RPCGetAddressTxIdsResponse{}
					navData.Address = address

					// bring in rpc get address txids struct to setup rpc params
					rpcGetParams := RPCGetAddressTxIdsParams{}

					// add current address to addresses array for rpc params
					rpcGetParams.Addresses = append(rpcGetParams.Addresses, address)

					// prepare rpc call, method and params (addresses array)
					n := daemonrpc.RpcRequestData{}
					n.Method = "getaddresstxids"
					n.Params = []RPCGetAddressTxIdsParams{rpcGetParams}

					// override credentials temporarily
					//conf.NavConf.RPCUser = "user"
					//conf.NavConf.RPCPassword = "hi"

					// issue rpc call
					resp, err := daemonrpc.RequestDaemon(n, conf.NavConf)

					// handle errors
					if err != nil {
						daemonrpc.RpcFailed(err, w, r)
						return
					}

					log.Println(resp)

					// returnedAddresses := RPCGetAddressTxIdsResponse{}

					// append response (array of txids) to TxIds array in returned address struct
					// returnedAddresses.TxIds = append(returnedAddresses.TxIds, resp)
				}

			}
		}

		return

	})
}
