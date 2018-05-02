package conf

// DaemonConfig defines a structure to store rpc data
type DaemonConfig struct {
	RPCUser     string `json:"rpcUser"`
	RPCPassword string `json:"rpcPassword"`
}

// CreateRPCDetails generate rpc details for this run
func CreateRPCDetails() {

	// only changed to hardcoded values for development
	NavConf.RPCUser = "user"
	//NavConf.RPCUser, _ = utils.GenerateRandomString(32)

	// only changed to hardcoded values for development
	NavConf.RPCPassword = "hi"
	//NavConf.RPCPassword, _ = utils.GenerateRandomString(32)

}
