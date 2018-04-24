package wallet

import (
	"net/http"

	"encoding/json"
	"fmt"

	"log"

	"github.com/Encrypt-S/kauri-api/app/api"
	"github.com/gorilla/mux"
	"github.com/Encrypt-S/kauri-api/app/daemon/daemonrpc"
	"github.com/Encrypt-S/kauri-api/app/conf"
)

// InitTransactionHandlers sets up handlers for transaction-related rpc commands
func InitTransactionHandlers(r *mux.Router, prefix string) {

	namespace := "transactions"

	// setup endpoint to be used for receiving txids for supplied addresses
	getAddressTxIdsPath := api.RouteBuilder(prefix, namespace, "v1", "getaddresstxids")
	api.OpenRouteHandler(getAddressTxIdsPath, r, getAddressTxIds())

}

// TransactionsArrayStruct describes the Transactions array
type TransactionsArrayStruct struct {
	Transactions []TransactionsStruct `json:"transactions"`
}

// TransactionsStruct describes the keys in Transactions array
type TransactionsStruct struct {
	Currency       string   `json:"currency"`
	Addresses      []string `json:"addresses"`
}

type TransactionResponseStruct struct {
	Transactions []ParsedTransactionDataStruct `json:"transactions"`
}

type ParsedTransactionDataStruct struct {
	Currency     string        `json:"currency"`
	Address      string        `json:"address"`
	Transactions []interface{} `json:"transactions"`
}

// getAddressTxIds - ranges through transactions, returns raw transactions
func getAddressTxIds() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var getAddressTxIds TransactionsArrayStruct
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

		// range over the Transactions slice
		for _, tx := range getAddressTxIds.Transactions {
			// isolate NAV addresses
			if tx.Currency == "NAV" {
				getNavTxIds(tx.Addresses)
			}
		}

	})
}

type ParsedAddresses struct {
	Address string
	TxIds []string
}

type TxLookups struct {
	TxIdsToLookup []ParsedAddresses
}


type rpcGetaddressTxIds struct {
	Addresses []string `json:"addresses"`
}

func getNavTxIds(addresses []string) {

	// loop through all the NAV addresses
	for _, add := range addresses {

		// add each address to
		txId := ParsedAddresses{}
		txId.Address = add

		// bring in transactions struct
		tx := TransactionsStruct{}

		// add current address to addresses array
		tx.Addresses = append(tx.Addresses, add)

		// then get transactions for each
		n := daemonrpc.RpcRequestData{}
		n.Method = "getaddresstxids"
		n.Params = []TransactionsStruct{tx}

		// override credentials temporarily
		conf.NavConf.RPCUser = "user"
		conf.NavConf.RPCPassword = "hi"

		// issue rpc call
		resp, err := daemonrpc.RequestDaemon(n, conf.NavConf)

		// handle errors
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(resp)

		// append to array

	}

	// return responseTx

}
