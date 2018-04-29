package wallet

import (
	"testing"
	"encoding/json"
	"bytes"
	"github.com/stretchr/testify/assert"
	"log"
	"gopkg.in/jarcoal/httpmock.v1"
)

func getTxIdsData() string {
	return `{"result":[
  "11a7071a43a8da2b9ac116865a6cd92c985c3f7cbde63933d253f88dffaa311a",
  "c6b6063a0512ed40958bff62a48168b4b30f89cb6bce22b722f8a6d00fcb9d98",
  "08f87e9de0fd9be71bc91f42d45c48bb9494df5d5df47df7354eec0adbf35731",
  "c8dad515d5e5c7a45bc5b3814fcf5e1f63474c9b67f84ee2ab9803f809e94929",
  "52489abff43212445d432f6042e5b9faf99b3c843a79210629b5383f52694ec5",
  "01f7b0831f174beb8a9b0990ca8bae197f6f1e4fe3d306c755d9f52da5687a9d"
]}`

}

// This functions mocks out the data struct
func setupIncomingTestData(t *testing.T) AddressesReq  {
	data := `
	{"transactions": [
    {"currency":  "NAV", "addresses": ["NW7uXr4ZAeJKigMGnKbSLfCBQY59cH1T8G", "NUDke42E3fwLqaBbBFRyVSTETuhWAi7ugk"]},
    {"currency":  "BTC",  "addresses": ["Bak7ahbZAA", "B91janABsa"]},
    {"currency":  "NAV", "addresses": ["NTkWY7kqiwoETFz8FUiaQoATLnwSYWTgvJ", "NeDXdkRkqDxav1KX5JxDLevGiSLDuhEBVY"]},
    {"currency":  "BTC",  "addresses": ["Bak7ahbZAA", "B91janABsa"]}
	]}`

	r := bytes.NewReader([]byte(data))

	var incomingAddressesReq AddressesReq
	json.NewDecoder(r).Decode(&incomingAddressesReq)


	//Preflight checks
	assert.Equal(t, "NAV", incomingAddressesReq.Addresses[0].Currency)
	assert.Equal(t, "NW7uXr4ZAeJKigMGnKbSLfCBQY59cH1T8G", incomingAddressesReq.Addresses[0].Addresses[0])
	assert.Equal(t, "NUDke42E3fwLqaBbBFRyVSTETuhWAi7ugk", incomingAddressesReq.Addresses[0].Addresses[1])


	assert.Equal(t, "BTC", incomingAddressesReq.Addresses[1].Currency)
	assert.Equal(t, "Bak7ahbZAA", incomingAddressesReq.Addresses[1].Addresses[0])
	assert.Equal(t, "B91janABsa", incomingAddressesReq.Addresses[1].Addresses[1])



	return incomingAddressesReq
}


func Test_buildResponse(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://127.0.0.1:0",
		httpmock.NewStringResponder(200, getTxIdsData()))

	incomingAddreses := setupIncomingTestData(t)
	resp := buildResponse(incomingAddreses)


	// Assert we have only nav currencies
	for i := range resp.Results {
		assert.Equal(t, "NAV", resp.Results[i].Currency)
	}

	// we have the right amount of results
	assert.Equal(t, 2, len(resp.Results))

	// top level scan of the struct to make sure things are in order - comprehensive tests are performed at each function
	assert.Equal(t, 2, len(resp.Results[0].Addresses))
	assert.Equal(t, "NW7uXr4ZAeJKigMGnKbSLfCBQY59cH1T8G", resp.Results[0].Addresses[0].Address)
	assert.Equal(t, "NUDke42E3fwLqaBbBFRyVSTETuhWAi7ugk", resp.Results[0].Addresses[1].Address)

}

func Test_getTransactionsForAddress(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://127.0.0.1:0",
		httpmock.NewStringResponder(200, getTxIdsData()))


	incomingAddresses := setupIncomingTestData(t)

	adds := getTransactionsForAddresses(incomingAddresses.Addresses[0].Addresses)

	// check we have the right amount of addresses
	assert.Equal(t, 2, len(adds))

	//check we have the addresses
	assert.Equal(t, "NW7uXr4ZAeJKigMGnKbSLfCBQY59cH1T8G", adds[0].Address)
	assert.Equal(t, "NUDke42E3fwLqaBbBFRyVSTETuhWAi7ugk", adds[1].Address)

	//assert.Equal(t, "11a7071a43a8da2b9ac116865a6cd92c985c3f7cbde63933d253f88dffaa311a", adds[1].Transactions[0].Txid)
	//assert.Equal(t, "c6b6063a0512ed40958bff62a48168b4b30f89cb6bce22b722f8a6d00fcb9d98", adds[1].Transactions[1].Txid)

}


func Test_getTxIDForAddressFromDaemon(t *testing.T) {


	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://127.0.0.1:0",
		httpmock.NewStringResponder(200, getTxIdsData()))

	rpcResp := getTxIDForAddressFromDaemon("NW7uXr4ZAeJKigMGnKbSLfCBQY59cH1T8G")
	log.Println(rpcResp)

	assert.Equal(t, "11a7071a43a8da2b9ac116865a6cd92c985c3f7cbde63933d253f88dffaa311a", rpcResp.Result[0])
	assert.Equal(t, "52489abff43212445d432f6042e5b9faf99b3c843a79210629b5383f52694ec5", rpcResp.Result[4])

}