package conf

import (
	"github.com/spf13/viper"
)

// ServerConfig defines a structure to store server config data
type ServerConfig struct {
	ManagerAPIPort int64 `json:"managerApiPort"`
}

// LoadServerConfig sets up viper, reads and parses server config
func LoadServerConfig() (ServerConfig, error) {

	viper.SetConfigName("server-config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./app")
	viper.AddConfigPath("../")

	err := viper.ReadInConfig()

	if err != nil {
		return ServerConfig{}, err
	}

	serverConfig := parseServerConfig(ServerConfig{})

	ServerConf = serverConfig

	return serverConfig, nil

}

// parseServerConfig takes ServerConfig, parses and returns serverconf
func parseServerConfig(serverconf ServerConfig) ServerConfig {

	serverconf.ManagerAPIPort = viper.GetInt64("managerApiPort")

	return serverconf

}
