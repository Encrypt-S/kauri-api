package conf

import (
	"github.com/spf13/viper"
)

// ServerConfig defines a structure to store server config data
type ServerConfig struct {
	ManagerAPIPort int64 `json:"managerApiPort"`
}

// LoadServerConfig sets up viper, reads and parses server config
func LoadServerConfig() error {

	viper.SetConfigName("server-config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./app")
	viper.AddConfigPath("../")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// parse out the config
	var serverConfig = ServerConfig{}
	err = viper.Unmarshal(&serverConfig)

	ServerConf = serverConfig

	return nil

}
