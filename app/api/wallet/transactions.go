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
)

// InitTransactionHandlers sets up handlers for transaction-related rpc commands
func InitTransactionHandlers(r *mux.Router, prefix string) {

	namespace := "transactions"

	// get raw transactions endpoint :: provides raw transaction data for supplied wallet addresses
	getRawTransactionsPath := api.RouteBuilder(prefix, namespace, "v1", "getrawtransactions")
	api.OpenRouteHandler(getRawTransactionsPath, r, getRawTransactionsHandler())

}

// Response describes top-level response object
type Response struct {
	Results []Result `json:"results"`
}

// Result describes each item returned in results array
type Result struct {
	Currency  string    `json:"currency"`
	Addresses []Address `json:"addresses"`
}

// Address describes each item returned in IncomingTxItems array
type Address struct {
	Address      string        `json:"address"`
	Transactions []Transaction `json:"transactions"`
}

// Transaction describes each item returned in transactions array
type Transaction struct {
	TxID    string      `json:"txid"`
	RawTx   string      `json:"rawtx"`
	Verbose interface{} `json:"verbose"`
}

// IncomingTransactions describes the incoming transactions in POST body
type IncomingTransactions struct {
	IncomingTxItems []IncomingTxItem `json:"transactions"`
}

// IncomingTxItem describes incoming transaction items
type IncomingTxItem struct {
	Currency  string   `json:"currency"`
	Addresses []string `json:"addresses"`
}

// GetTxIDParams describes addresses array params for 'getaddresstxids' RPC call
type GetTxIDParams struct {
	Addresses []string `json:"addresses"`
}

// GetTxIDsResp describes Result of RPC response > txid (array)
type GetTxIDsResp struct {
	Result []string `json:"result"`
}

// GetRawTxResp descibes Result of RPC response > raw (string)
type GetRawTxResp struct {
	Result string `json:"result"`
}

// getRawTransactionsHandler ranges through transactions, returns RPC response data
func getRawTransactionsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiResp := api.Response{}

		var incomingTxs IncomingTransactions

		err := json.NewDecoder(r.Body).Decode(&incomingTxs)

		if err != nil {
			returnErr := api.AppRespErrors.JSONDecodeError
			returnErr.ErrorMessage = fmt.Sprintf("JSON decode error: %v", err)
			apiResp.Errors = append(apiResp.Errors, returnErr)
			apiResp.Send(w)
			return
		}

		resp, err := buildResponse(incomingTxs)

		if err != nil {
			returnErr := api.AppRespErrors.RPCResponseError
			returnErr.ErrorMessage = fmt.Sprintf("Response error: %v", err)
			apiResp.Errors = append(apiResp.Errors, returnErr)
			apiResp.Send(w)
			return
		}

		apiResp.Data = resp.Results

		apiResp.Send(w)

		return
	})
}

// buildResponse takes address and returns response data
func buildResponse(incomingAddreses IncomingTransactions) (Response, error) {

	resp := Response{}

	// loop through all the lines that we received
	for _, item := range incomingAddreses.IncomingTxItems {

		if item.Currency == "NAV" {

			// setup the result struct we will use later
			result := Result{}
			result.Currency = "NAV"

			// get transaction related to the address and store them in the result
			result.Addresses, _ = getTransactionsForAddresses(item.Addresses)

			// append the result to the array of results
			resp.Results = append(resp.Results, result)

		}

	}

	return resp, nil

}

// getTransactionsForAddresses takes addresses array and returns data for each
func getTransactionsForAddresses(addresses []string) ([]Address, error) {

	adds := []Address{}

	for _, addressStr := range addresses {

		addStruct := Address{}
		addStruct.Address = addressStr
		rpcTxIDsResp, err := getTxIdsRPC(addressStr)

		if err != nil {
			return nil, err
		}

		// for all the txIDs from the rpc we need to create a transaction
		for _, txID := range rpcTxIDsResp.Result {

			rawTx, _ := getRawTx(txID)

			if err != nil {
				return nil, err
			}

			verboseTx, _ := getRawTxVerbose(txID)

			if err != nil {
				return nil, err
			}

			trans := Transaction{TxID: txID, RawTx: rawTx.Result, Verbose: verboseTx}
			addStruct.Transactions = append(addStruct.Transactions, trans)

		}

		adds = append(adds, addStruct)

	}

	return adds, nil
}

// getTxIdsRPC takes address and returns array of txids
func getTxIdsRPC(address string) (GetTxIDsResp, error) {

	getParams := GetTxIDParams{}

	getParams.Addresses = append(getParams.Addresses, address)

	n := daemonrpc.RpcRequestData{}
	n.Method = "getaddresstxids"
	n.Params = []GetTxIDParams{getParams}

	resp, err := daemonrpc.RequestDaemon(n, conf.NavConf)

	if err != nil {
		return GetTxIDsResp{}, err
	}

	rpcTxIDResults := GetTxIDsResp{}

	err = json.NewDecoder(resp.Body).Decode(&rpcTxIDResults)

	if err != nil {
		return GetTxIDsResp{}, err
	}

	return rpcTxIDResults, nil

}

// getRawTx takes txid and returns raw transaction data
func getRawTx(txid string) (GetRawTxResp, error) {

	n := daemonrpc.RpcRequestData{}
	n.Method = "getrawtransaction"
	n.Params = []string{txid}

	resp, err := daemonrpc.RequestDaemon(n, conf.NavConf)
	if err != nil {
		return GetRawTxResp{}, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var rawData map[string]string

	if err := json.Unmarshal(bodyBytes, &rawData); err != nil {
		return GetRawTxResp{}, err
	}

	rawResp := GetRawTxResp{}
	rawResp.Result = rawData["result"]

	return rawResp, nil

}

// getRawTxVerbose takes txid and returns verbose transaction data
func getRawTxVerbose(txid string) (interface{}, error) {

	n := daemonrpc.RpcRequestData{}
	n.Method = "getrawtransaction"
	n.Params = []interface{}{txid, 1}

	resp, err := daemonrpc.RequestDaemon(n, conf.NavConf)

	if err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, err
	}

	return data["result"], nil

}
