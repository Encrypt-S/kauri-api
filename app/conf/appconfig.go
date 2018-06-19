package conf

import (
	"github.com/spf13/viper"
)

// AppConfig defines a structure to store app config data
// at present all config is mapped to Coins array
type AppConfig struct {
	Coins []CoinData `json:"coins"`
}

// Coins defines properties of active coin
type CoinData struct {
	Name              string `json:"name"`
	CurrencyCode      string `json:"currencyCode"`
	LibPath           string `json:"libPath"`
	DataDir           string `json:"dataDir"`
	DaemonHeartbeat   int    `json:"daemonHeartbeat"`
	DaemonVersion     string `json:"daemonVersion"`
	WindowsDaemonName string `json:"windowsDaemonName"`
	DarwinDaemonName  string `json:"darwinDaemonName"`
	LatestReleaseAPI  string `json:"latestReleaseApi"`
	ReleaseAPI        string `json:"ReleaseApi"`
	LivePort          int    `json:"livePort"`
	TestNetPort       int    `json:"testnetPort"`
	UseTestNet        bool   `json:"useTestNet"`
	IndexTransactions bool   `json:"indexTransactions"`
}

// LoadAppConfig sets up viper, reads and parses app config
func LoadAppConfig() error {

	viper.SetConfigName("app-config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./app")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// parse out the config
	var appConfig = AppConfig{}
	err = viper.Unmarshal(&appConfig)

	if err != nil {
		return err
	}

	AppConf = appConfig

	return nil
}
