package wallet

import (
	"net/http"

	"encoding/json"
	"fmt"

	"github.com/Encrypt-S/kauri-api/app/api"
	"github.com/Encrypt-S/kauri-api/app/conf"
	"github.com/Encrypt-S/kauri-api/app/daemon/daemonrpc"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
)

// InitTransactionHandlers sets up handlers for transaction-related rpc commands
func InitTransactionHandlers(r *mux.Router, prefix string) {

	namespace := "transactions"

	// setup "getrawtransactions" endpoint - provides raw transaction data for supplied wallet addresses
	getRawTransactionsPath := api.RouteBuilder(prefix, namespace, "v1", "getrawtransactions")
	api.OpenRouteHandler(getRawTransactionsPath, r, getRawTransactions())

}

type Response struct {
	Results []Result `json:"results"`
}
type Result struct {
	Currency  string    `json:"currency"`
	Addresses []Address `json:"addresses"`
}
type Address struct {
	Address      string        `json:"address"`
	Transactions []Transaction `json:"transactions"`
}
type Transaction struct {
	Txid    string      `json:"txid"`
	Rawtx   string      `json:"rawtx"`
	Verbose interface{} `json:"verbose"`
}

// IncomingTransactionsArray describes the Transactions array
type AddressesReq struct {
	Addresses []AddressReqItem `json:"transactions"`
}

// IncomingTransactions describes the incoming Transactions array
type AddressReqItem struct {
	Currency  string   `json:"currency"`
	Addresses []string `json:"addresses"`
}

// RPCGetAddressTxIDParams describes addresses array params for getaddresstxids call
type RPCGetAddressTxIDParams struct {
	Addresses []string `json:"addresses"`
}

// GetTxIdsRPCResp contains RPC response :: txid array
type GetTxIdsRPCResp struct {
	Result []string `json:"result"`
}

// RPCGetAddressRawTx contains RPC response :: raw tx
type RPCGetAddressRawTx struct {
	Result string `json:"result"`
}

// getRawTransactions - ranges through transactions, returns rpc response data
func getRawTransactions() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiResp := api.Response{}

		var incomingTxs AddressesReq

		err := json.NewDecoder(r.Body).Decode(&incomingTxs)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			returnErr := api.AppRespErrors.ServerError
			returnErr.ErrorMessage = fmt.Sprintf("Server error: %v", err)
			apiResp.Errors = append(apiResp.Errors, returnErr)
			apiResp.Send(w)
			return
		}

		resp := buildResponse(incomingTxs)

		apiResp.Data = resp

		apiResp.Send(w)

		return
	})
}

// buildResponse takes address and returns response data
func buildResponse(incomingAddreses AddressesReq) Response {

	resp := Response{}

	// loop through all the lines that we received
	for _, item := range incomingAddreses.Addresses {

		if item.Currency == "NAV" {

			// setup the result struct we will use later
			result := Result{}
			result.Currency = "NAV"

			// get transaction related to the address and store them in the result
			result.Addresses = getTransactionsForAddresses(item.Addresses)

			// append the result to the array of results
			resp.Results = append(resp.Results, result)

		}

	}

	return resp

}

// getTransactionsForAddresses takes addresses array and returns data for each
func getTransactionsForAddresses(addresses []string) []Address {

	adds := []Address{}

	for _, addressStr := range addresses {

		addStruct := Address{}
		addStruct.Address = addressStr
		rpcTxIDsResp := getTxIdsRPC(addressStr)

		// for all the txIds from the rpc we need to create a transaction
		for _, txId := range rpcTxIDsResp.Result {

			rawTx := getRawTx(txId)
			verboseTx := getVerboseTx(txId)

			trans := Transaction{Txid: txId, Rawtx: rawTx.Result, Verbose: verboseTx}
			addStruct.Transactions = append(addStruct.Transactions, trans)

		}

		adds = append(adds, addStruct)

	}

	return adds
}

// getTxIdsRPC takes address and returns array of txids
func getTxIdsRPC(address string) GetTxIdsRPCResp {

	getParams := RPCGetAddressTxIDParams{}

	getParams.Addresses = append(getParams.Addresses, address)

	n := daemonrpc.RpcRequestData{}
	n.Method = "getaddresstxids"
	n.Params = []RPCGetAddressTxIDParams{getParams}

	resp, rpcErr := daemonrpc.RequestDaemon(n, conf.NavConf)

	if rpcErr != nil {
		//daemonrpc.RpcFailed(rpcErr, w, r)
	}

	rpcTxIdResults := GetTxIdsRPCResp{}

	jsonErr := json.NewDecoder(resp.Body).Decode(&rpcTxIdResults)

	if jsonErr != nil {
		//returnErr := api.AppRespErrors.JSONDecodeError
		//returnErr.ErrorMessage = fmt.Sprintf("TxId JSON Decode Error: %v", jsonErr)
		//apiResp.Errors = append(apiResp.Errors, returnErr)
		//apiResp.Send(w)
	}

	return rpcTxIdResults

}

func getRawTx(txid string) RPCGetAddressRawTx {

	n := daemonrpc.RpcRequestData{}
	n.Method = "getrawtransaction"
	n.Params = []string{txid}

	resp, rpcErr := daemonrpc.RequestDaemon(n, conf.NavConf)
	if rpcErr != nil {
		//log.Printf()
		log.Println("getRawTx rpcErr", rpcErr)
	}

	// TODO: PAUL look at why this is not auto unmarshalling
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var dat map[string]string

	if err := json.Unmarshal(bodyBytes, &dat); err != nil {
		panic(err)
	}

	rawResp := RPCGetAddressRawTx{}
	rawResp.Result = dat["result"]

	return rawResp

}

// TODO: Write test for this...
func getVerboseTx(txid string) interface{} {

	n := daemonrpc.RpcRequestData{}
	n.Method = "getrawtransaction"
	n.Params = []interface{}{txid, 1}

	resp, rpcErr := daemonrpc.RequestDaemon(n, conf.NavConf)
	if rpcErr != nil {
		log.Println("getRawTx rpcErr", rpcErr)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var dat map[string]interface{}

	if err := json.Unmarshal(bodyBytes, &dat); err != nil {
		panic(err)
	}

	return dat["result"]

}
