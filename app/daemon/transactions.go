package daemon


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


