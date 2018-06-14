package daemonrpc

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"fmt"

	"github.com/Encrypt-S/kauri-api/app/conf"
)

// RPCRequestData defines method and params
type RPCRequestData struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// RPCResponse defines code, data, and message
type RPCResponse struct {
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

// RequestDaemon requests the data via the daemon's RPC api
// allows auto switches between the testnet and live depending on the config
func RequestDaemon(coinData conf.CoinData, rpcReqData RPCRequestData, daemonConf conf.DaemonConfig) (*http.Response, error) {

	username := daemonConf.RPCUser
	password := daemonConf.RPCPassword

	client := &http.Client{}

	jsonValue, _ := json.Marshal(rpcReqData)

	// set the live port
	port := coinData.LivePort

	// build the url
	url := fmt.Sprintf("http://127.0.0.1:%d", port)

	// build the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	req.SetBasicAuth(username, password)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	return resp, err

}

// RPCFailed handles errors encountered when requesting daemon
func RPCFailed(err error, w http.ResponseWriter) {

	resp := RPCResponse{}

	w.WriteHeader(http.StatusFailedDependency)
	resp.Code = http.StatusFailedDependency
	resp.Message = "Failed to run command: " + err.Error()
	log.Fatal("Failed to run command: " + err.Error())

	respJSON, err := json.Marshal(resp)

	if err != nil {

	}

	w.Write(respJSON)

}
