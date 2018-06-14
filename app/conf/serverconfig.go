package conf

import (
	"github.com/spf13/viper"
)

// ServerConfig defines a structure to store server config data
type ServerConfig struct {
	ManagerAPIPort int64 `json:"managerApiPort"`
}

// LoadServerConfig sets up viper, reads and parses server config
//func LoadServerConfig() (ServerConfig, error) {
//
//	viper.SetConfigName("server-config")
//	viper.AddConfigPath(".")
//	viper.AddConfigPath("./")
//	viper.AddConfigPath("./app")
//	viper.AddConfigPath("../")
//
//	err := viper.ReadInConfig()
//
//	if err != nil {
//		return ServerConfig{}, err
//	}
//
//	serverConfig := parseServerConfig(ServerConfig{})
//
//	ServerConf = serverConfig
//
//	return serverConfig, nil
//
//}

func LoadServerConfig() error {

	viper.SetConfigName("server-config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./app")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// parse out the config
	var serverConfig ServerConfig
	err = viper.UnmarshalKey("config", &serverConfig)

	if err != nil {
		return err
	}

	ServerConf = serverConfig

	return nil
}

// parseServerConfig takes ServerConfig, parses and returns serverconf
//func parseServerConfig(serverconf ServerConfig) ServerConfig {
//
//	// parse out the config
//	var serverConfig ServerConfig
//	err = viper.UnmarshalKey("config", &appConfig)
//
//	serverconf.ManagerAPIPort = viper.GetInt64("managerApiPort")
//
//	serverconf.LivePort = viper.GetInt64("navCoinPorts.livePort")
//	serverconf.TestPort = viper.GetInt64("navCoinPorts.testnetPort")
//	serverconf.UseTestnet = viper.GetBool("useTestnet")
//
//	return serverconf
//
//}
