package daemonapi

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

// InitTxHandlers sets up handlers for transaction-related rpc commands
func InitTxHandlers(r *mux.Router, coinData conf.CoinData, prefix string) {

	namespace := "transactions"

	// get raw transactions endpoint :: provides raw transaction data for supplied wallet addresses
	getRawTransactionsPath := api.RouteBuilder(prefix, namespace, "v1", "getrawtransactions")
	api.OpenRouteHandler(getRawTransactionsPath, r, getRawTxHandler(coinData))

}

// TxResponse is the top-level response object
type TxResponse struct {
	Results []WalletResult `json:"results"`
}

// WalletResult is all transactions for array of address by currency
type WalletResult struct {
	Currency  string                `json:"currency"`
	Addresses []AddressTransactions `json:"addresses"`
}

// AddressTransactions contains array of transactions for address
type AddressTransactions struct {
	Address      string        `json:"address"`
	Transactions []Transaction `json:"transactions"`
}

// Transaction contains txid, raw transactions txid, verbose boolean
type Transaction struct {
	TxID    string      `json:"txid"`
	RawTx   string      `json:"rawtx"`
	Verbose interface{} `json:"verbose"`
}

// IncomingTransactions are the incoming transactions in POST body
type IncomingTransactions struct {
	IncomingTxItems []WalletItem `json:"transactions"`
}

// WalletItem is the incoming currency and corresponding addresses
type WalletItem struct {
	Currency  string   `json:"currency"`
	Addresses []string `json:"addresses"`
}

// GetTxIDParams are the addresses array params for 'getaddresstxids' RPC call
type GetTxIDParams struct {
	Addresses []string `json:"addresses"`
}

// GetTxIDsResp is the Result of RPC response > txid (array)
type GetTxIDsResp struct {
	Result []string `json:"result"`
}

// GetRawTxResp is the Result of RPC response > raw (string)
type GetRawTxResp struct {
	Result string `json:"result"`
}

// getRawTxHandler ranges through transactions, returns RPC response data
func getRawTxHandler(coinData conf.CoinData) http.Handler {
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

		resp, err := buildResponse(coinData, incomingTxs)

		if err != nil {
			returnErr := api.AppRespErrors.RPCResponseError
			returnErr.ErrorMessage = fmt.Sprintf("TxResponse error: %v", err)
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
func buildResponse(coinData conf.CoinData, incomingAddreses IncomingTransactions) (TxResponse, error) {

	resp := TxResponse{}

	// loop through all the lines that we received
	for _, item := range incomingAddreses.IncomingTxItems {

		if item.Currency == "NAV" {

			// setup the result struct we will use later
			result := WalletResult{}
			result.Currency = "NAV"

			// get transaction related to the address and store them in the result
			result.Addresses, _ = getTxForAddresses(coinData, item.Addresses)

			// append the result to the array of results
			resp.Results = append(resp.Results, result)

		}

	}

	return resp, nil

}

// getTxForAddresses takes addresses array and returns data for each
func getTxForAddresses(coinData conf.CoinData, addresses []string) ([]AddressTransactions, error) {

	adds := []AddressTransactions{}

	for _, addressStr := range addresses {

		addStruct := AddressTransactions{}
		addStruct.Address = addressStr
		rpcTxIDsResp, err := getTxIdsRPC(coinData, addressStr)

		if err != nil {
			return nil, err
		}

		// for all the txIDs from the rpc we need to create a transaction
		for _, txID := range rpcTxIDsResp.Result {

			rawTx, _ := getRawTx(coinData, txID)

			if err != nil {
				return nil, err
			}

			verboseTx, _ := getRawTxVerbose(coinData, txID)

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
func getTxIdsRPC(coinData conf.CoinData, address string) (GetTxIDsResp, error) {

	getParams := GetTxIDParams{}

	getParams.Addresses = append(getParams.Addresses, address)

	reqData := daemonrpc.RPCRequestData{}
	reqData.Method = "getaddresstxids"
	reqData.Params = []GetTxIDParams{getParams}

	resp, err := daemonrpc.RequestDaemon(coinData, reqData, conf.DaemonConf)

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
func getRawTx(coinData conf.CoinData, txid string) (GetRawTxResp, error) {

	rpcData := daemonrpc.RPCRequestData{}
	rpcData.Method = "getrawtransaction"
	rpcData.Params = []string{txid}

	resp, err := daemonrpc.RequestDaemon(coinData, rpcData, conf.DaemonConf)
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
func getRawTxVerbose(coinData conf.CoinData, txid string) (interface{}, error) {

	rpcData := daemonrpc.RPCRequestData{}
	rpcData.Method = "getrawtransaction"
	rpcData.Params = []interface{}{txid, 1}

	resp, err := daemonrpc.RequestDaemon(coinData, rpcData, conf.DaemonConf)

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
