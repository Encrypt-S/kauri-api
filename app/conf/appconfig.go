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
type AppConfig struct {
	NavConf           string   `json:"navconf"`
	RunningNavVersion string   `json:"runningNavVersion"`
	AllowedIps        []string `json:"allowedIps"`
	UIPassword        string   `json:"uiPassword"`
	Coins             []Coins  `json:"coins"`
}

// Coins defines structure of supported coin data
type Coins struct {
	DaemonVersion    string `json:"daemonVersion"`
	CurrencyCode     string `json:"currencyCode"`
	LatestReleaseAPI string `json:"latestReleaseApi"`
	ReleaseAPI       string `json:"ReleaseApi"`
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

	//parse out the config
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
		NavConf:           AppConf.NavConf,
		RunningNavVersion: AppConf.RunningNavVersion,
		AllowedIps:        AppConf.AllowedIps,
		UIPassword:        AppConf.UIPassword,
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
