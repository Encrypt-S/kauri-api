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
	getTxIdsPath := api.RouteBuilder(prefix, namespace, "v1", "getaddresstxids")
	api.OpenRouteHandler(getTxIdsPath, r, getData("txids"))

	// setup endpoint to be used for receiving raw transaction datat for all supplied addresses
	getRawTransactionsPath := api.RouteBuilder(prefix, namespace, "v1", "getrawtransactions")
	api.OpenRouteHandler(getRawTransactionsPath, r, getData("raw"))

}

type Response struct {
	Data struct {
		Results []struct {
			Currency  string `json:"currency"`
			Addresses []struct {
				Address      string `json:"address"`
				Transactions []struct {
					Txid    string `json:"txid"`
					Rawtx   string `json:"rawtx"`
					Verbose string `json:"verbose"`
				} `json:"transactions"`
			} `json:"addresses"`
		} `json:"results"`
	} `json:"data"`
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
//type OutgoingTransactionsArray struct {
//	Transactions []OutgoingTransactions `json:"transactions"`
//}

// OutgoingTransactions describes the outgoing response
//type OutgoingTransactions struct {
//	Currency          string                  `json:"currency"`
//	OutgoingAddresses []OutgoingAddressObject `json:"addressobject"`
//}

// OutgoingAddressObject contains address and array of txids
//type OutgoingAddressObject struct {
//	Address            string   `json:"address"`
//	OutgoingTxIdsArray []string `json:"txids"`
//}

// RPCGetAddressTxIDParams describes addresses array params for getaddresstxids call
type RPCGetAddressTxIDParams struct {
	Addresses []string `json:"addresses"`
}

// RPCGetDataFromAddressResponse contains RPC response :: txid array or raw tx
//type RPCGetDataFromAddressResponse struct {
//	Result []string `json:"result"`
//}

// RPCRawTxResponse contains RPC response :: raw transaction data
//type RPCRawTxResponse struct {
//	Result []string `json:"result"`
//}

// getData - ranges through transactions, returns txids or raw transactions
func getData(cmd string) http.Handler {
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

		//txarr := OutgoingTransactionsArray{}

		//txresp := TransactionReponse{}

		resp := Response{}
		results := resp.Data.Results

		for _, tx := range incomingTxs.Transactions {

			if tx.Currency == "NAV" {
				getDataFromAddresses(tx.Addresses, &results)
				//txarr.Transactions = append(txarr.Transactions, txids)
			}

		}

		apiResp.Data = resp

		apiResp.Send(w)

		return
	})
}

// getDataFromAddresses returns txids from supplied addresses
func getDataFromAddresses(addresses []string, results []string) Response {

	results := resp.Data.Results
	results......
	//txout := OutgoingTransactions{}
	//txout.Currency = "NAV"


	for _, address := range addresses {
		txIDs := getDataFromAddress(address)
		txout.OutgoingAddresses = append(txout.OutgoingAddresses, createResponseObject(address, txIDs))
	}

	return txout

}

// getDataFromAddress issuces RPC calls, returns response (txid array)
func getDataFromAddress(address string, cmd string) RPCGetDataFromAddressResponse {

	apiResp := api.Response{}

	getParams := RPCGetAddressTxIDParams{}

	getParams.Addresses = append(getParams.Addresses, address)

	n := daemonrpc.RpcRequestData{}
	n.Method = "getaddresstxids"
	n.Params = []RPCGetAddressTxIDParams{getParams}

	resp, rpcErr := daemonrpc.RequestDaemon(n, conf.NavConf)

	if rpcErr != nil {
		//daemonrpc.RpcFailed(rpcErr, w, r)
	}

	txid := RPCGetDataFromAddressResponse{}

	jsonErr := json.NewDecoder(resp.Body).Decode(&txid)

	if jsonErr != nil {
		//returnErr := api.AppRespErrors.JSONDecodeError
		//returnErr.ErrorMessage = fmt.Sprintf("TxId JSON Decode Error: %v", jsonErr)
		//apiResp.Errors = append(apiResp.Errors, returnErr)
		//apiResp.Send(w)
	}

	//rawtx := RPCGetDataFromAddressResponse{}
	//
	//resp := txid
	//
	//if cmd == "raw" {
	//	for _, address := range txid.Result {
	//		rawtx.Result = getRawTransactionsFromTxId(txid.Result, w, r)
	//	}
	//	resp = append...
	//}

	return txid

}

// getRawTransactionFromTxId return the serialized, hex-encoded data for provided 'txid'
func getRawTransactionsFromTxId(txid string) RPCRawTxResponse {

	apiResp := api.Response{}

	n := daemonrpc.RpcRequestData{}
	n.Method = "getrawtransaction"
	n.Params = txid

	resp, rpcErr := daemonrpc.RequestDaemon(n, conf.NavConf)

	if rpcErr != nil {
		//daemonrpc.RpcFailed(rpcErr, w, r)
	}

	rawTx := RPCRawTxResponse{}

	jsonErr := json.NewDecoder(resp.Body).Decode(&rawTx)

	if jsonErr != nil {
		//returnErr := api.AppRespErrors.JSONDecodeError
		//returnErr.ErrorMessage = fmt.Sprintf("Raw Tx JSON Decode Error: %v", jsonErr)
		//apiResp.Errors = append(apiResp.Errors, returnErr)
		//apiResp.Send(w)
	}

	return rawTx

}

// createResponseObject formats the address, array of txids into outgoing address object
func createResponseObject(address string, txIDs RPCGetDataFromAddressResponse) OutgoingAddressObject {

	outAddObj := OutgoingAddressObject{}
	outAddObj.Address = address
	outAddObj.OutgoingTxIdsArray = txIDs.Result

	return outAddObj

}
