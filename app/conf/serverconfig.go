package conf

import (
	"github.com/spf13/viper"
)

type ServerConfig struct {
	ManagerApiPort   int64
	DaemonApiPort    int64
	SetupApiPort     int64
	LatestReleaseAPI string
	ReleaseAPI       string
	DaemonHeartbeat  int64
}

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

func parseServerConfig(config ServerConfig) ServerConfig {

	config.ManagerApiPort = viper.GetInt64("managerApiPort")
	config.LatestReleaseAPI = viper.GetString("latestReleaseAPI")
	config.ReleaseAPI = viper.GetString("releaseAPI")
	config.DaemonHeartbeat = viper.GetInt64("daemonHeartbeat")

	return config

}