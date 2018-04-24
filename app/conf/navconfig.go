package conf

// DaemonConfig defines a structure to store rpc data
type DaemonConfig struct {
	RPCUser     string `json:"rpcUser"`
	RPCPassword string `json:"rpcPassword"`
}

// CreateRPCDetails generate rpc details for this run
func CreateRPCDetails() {

	//NavConf.RPCUser, _ = utils.GenerateRandomString(32)
	NavConf.RPCUser = "user"
	//NavConf.RPCPassword, _ = utils.GenerateRandomString(32)
	NavConf.RPCPassword = "hi"

}
