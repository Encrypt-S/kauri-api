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

	// setup endpoint to be used for receiving txids for all supplied addresses
	getTxIdsForAllAddressesPath := api.RouteBuilder(prefix, namespace, "v1", "getaddresstxids")
	api.OpenRouteHandler(getTxIdsForAllAddressesPath, r, getTxIdsForAllAddresses())

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
	Currency          string                  `json:"currency"`
	OutgoingAddresses []OutgoingAddressObject `json:"addressobject"`
}

// OutgoingAddressObject contains address and array of txids
type OutgoingAddressObject struct {
	Address            string   `json:"address"`
	OutgoingTxIdsArray []string `json:"txids"`
}

// RPCGetAddressTxIDParams describes addresses array params for getaddresstxids call
type RPCGetAddressTxIDParams struct {
	Addresses []string `json:"addresses"`
}

// RPCAddressTxIDResponse contains RPC response
type RPCAddressTxIDResponse struct {
	Result []string `json:"result"`
}

// getTxIdsForAllAddresses - ranges through transactions, returns txids for each address
func getTxIdsForAllAddresses() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiResp := api.Response{}

		var incomingTxs IncomingTransactionsArray

		err := json.NewDecoder(r.Body).Decode(&incomingTxs)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			returnErr := api.AppRespErrors.ServerError
			returnErr.ErrorMessage = fmt.Sprintf("Server error: %v", err)
			apiResp.Errors = append(apiResp.Errors, returnErr)
			apiResp.Send(w)
			return
		}

		txarr := OutgoingTransactionsArray{}

		for _, tx := range incomingTxs.Transactions {

			if tx.Currency == "NAV" {
				trans := getTxIdsFromAddresses(tx.Addresses, w, r)
				txarr.Transactions = append(txarr.Transactions, trans)
			}

		}

		apiResp.Data = txarr
		apiResp.Send(w)

		return
	})
}

// getTxIdsFromAddresses returns txids from supplied addresses
func getTxIdsFromAddresses(addresses []string, w http.ResponseWriter, r *http.Request) OutgoingTransactions {

	outTrans := OutgoingTransactions{}
	outTrans.Currency = "NAV"

	for _, address := range addresses {
		txIDs := getTxIdsFromAddress(address, w, r)
		outTrans.OutgoingAddresses = append(outTrans.OutgoingAddresses, createResponseObject(address, txIDs))
	}

	return outTrans

}

// getTxIdsFromAddress issuces RPC calls, returns response (txid array)
func getTxIdsFromAddress(address string, w http.ResponseWriter, r *http.Request) RPCAddressTxIDResponse {

	apiResp := api.Response{}

	getParams := RPCGetAddressTxIDParams{}

	getParams.Addresses = append(getParams.Addresses, address)

	n := daemonrpc.RpcRequestData{}
	n.Method = "getaddresstxids"
	n.Params = []RPCGetAddressTxIDParams{getParams}

	resp, rpcErr := daemonrpc.RequestDaemon(n, conf.NavConf)

	if rpcErr != nil {
		daemonrpc.RpcFailed(rpcErr, w, r)
	}

	txidResp := RPCAddressTxIDResponse{}

	jsonErr := json.NewDecoder(resp.Body).Decode(&txidResp)

	if jsonErr != nil {
		returnErr := api.AppRespErrors.JSONDecodeError
		returnErr.ErrorMessage = fmt.Sprintf("JSON Decode Error: %v", jsonErr)
		apiResp.Errors = append(apiResp.Errors, returnErr)
		apiResp.Send(w)
	}

	return txidResp

}

// createResponseObject formats the address, array of txids into outgoing address object
func createResponseObject(address string, txIDs RPCAddressTxIDResponse) OutgoingAddressObject {

	outAddObj := OutgoingAddressObject{}
	outAddObj.Address = address
	outAddObj.OutgoingTxIdsArray = txIDs.Result

	return outAddObj

}
