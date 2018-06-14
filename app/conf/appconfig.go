package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/spf13/viper"
)

// AppConfig defines a structure to store app config data
// at present all config is mapped to Coins array
type AppConfig struct {
	Coins []Coin `json:"coin"`
}

// Coins defines structure of supported coin data
type Coin struct {
	Name             string `json:"name"`
	CurrencyCode     string `json:"currencyCode"`
	DaemonHeartbeat  int64  `json:"daemonHeartbeat"`
	DaemonVersion    string `json:"daemonVersion"`
	DataDir          string `json:"dataDir"`
	LatestReleaseAPI string `json:"latestReleaseApi"`
	ReleaseAPI       string `json:"ReleaseApi"`
	LivePort         int64  `json:"livePort"`
	TestNetPort      int64  `json:"testnetPort"`
	CmdAddressIndex  string `json:"cmdAddressIndex"`
	CmdNetwork       string `json:"cmdNetwork"`
}

// StartConfigManager sets up the ticker loop to load app config
func StartConfigManager() {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for range ticker.C {
			LoadAppConfig()
		}
	}()
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
	var appConfig AppConfig
	err = viper.UnmarshalKey("config", &appConfig)

	if err != nil {
		return err
	}

	AppConf = appConfig

	return nil
}

// SaveAppConfig formats/indents json and saves to app-config.json
func SaveAppConfig() error {

	jsonData, err := json.MarshalIndent(AppConfig{
		Coins: AppConf.Coins,
	}, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))

	path := "app/app-config.json"

	log.Println("attempting to write json data to " + path)

	err = ioutil.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil

}
