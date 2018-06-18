package daemon

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"fmt"

	"github.com/Encrypt-S/kauri-api/app/conf"
	"github.com/Encrypt-S/kauri-api/app/fs"
	"github.com/Encrypt-S/kauri-api/app/daemon/daemonrpc"
)

type OSInfo struct {
	DaemonName string
	OS         string
}

type GitHubReleases []struct {
	GitHubReleaseData
}

type GitHubReleaseData struct {
	URL             string `json:"url"`
	AssetsURL       string `json:"assets_url"`
	UploadURL       string `json:"upload_url"`
	HTMLURL         string `json:"html_url"`
	ID              int    `json:"id"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Draft           bool   `json:"draft"`
	Author          struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		URL      string      `json:"url"`
		ID       int         `json:"id"`
		Name     string      `json:"name"`
		Label    interface{} `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
	Body       string `json:"body"`
}

var runningDaemon *exec.Cmd
var minHeartbeat = 1000 // the lowest value the hb checker can be set to

var isGettingDaemon = false

// StartManager is a simple system that checks if the coin's daemon
// is alive. If not it tries to startCoinDaemons it with proper config
// It is called from the StartAllDaemonManagers function in managers pkg
func StartManager(coinData conf.CoinData) {

	log.Println("starting manager for " + coinData.CurrencyCode + " daemon...")

	// set the heartbeat interval but make sure it is not
	// less than the min heartbeat setting
	hbInterval := minHeartbeat
	if coinData.DaemonHeartbeat > hbInterval {
		hbInterval = coinData.DaemonHeartbeat
	}

	// check to see if the daemon is alive
	if isAlive(coinData) {
		log.Println(coinData.CurrencyCode + " daemon is alive!")
	} else {
		log.Println(coinData.CurrencyCode + " daemon is not yet alive...")

		// only do thing if we are already not getting the daemon
		if !isGettingDaemon {
			log.Println(coinData.CurrencyCode + " daemon is unresponsive...")

			if runningDaemon != nil {
				log.Println(coinData.CurrencyCode + " daemon is already running, just stop")
				Stop(coinData, runningDaemon)
			}

			// kick off goroutine for DownloadAndStart
			go func() {

				cmd, err := DownloadAndStart(coinData)

				if err != nil {
					log.Println(err)
				} else {
					runningDaemon = cmd
				}

			}()

		}
	}
}

// returns false on error
func isAlive(coinData conf.CoinData) bool {

	isLiving := true

	n := daemonrpc.RPCRequestData{}
	n.Method = "getblockcount"

	_, err := daemonrpc.RequestDaemon(coinData, n, conf.DaemonConf)

	if err != nil {
		isLiving = false
	}

	return isLiving

}

// DownloadAndStart checks for current coin's daemon
// and either downloads it or starts it up if already detected
func DownloadAndStart(coinData conf.CoinData) (*exec.Cmd, error) {

	path, err := CheckForDaemon(coinData)

	// download daemon if not found
	if err != nil {
		downloadDaemons(coinData)
	} else {
		return startCoinDaemons(coinData, path), nil
	}

	// if found, just start it up
	return startCoinDaemons(coinData, path), nil

}

func Stop(coinData conf.CoinData, cmd *exec.Cmd) {

	if err := cmd.Process.Kill(); err != nil {
		log.Fatal("failed to kill" + coinData.CurrencyCode + "process" + err.Error())
	}
}

// CheckForDaemon checks for current coin's daemon
// in appropriate path and reports back to DownLoadAndStartDaemons
func CheckForDaemon(coinData conf.CoinData) (string, error) {

	// get the latest release version, equal to daemon version
	releaseVersion := coinData.DaemonVersion

	log.Println("Checking" + coinData.CurrencyCode + "daemon for v" + releaseVersion)

	// get the apps current path
	path, err := fs.GetCurrentPath()
	if err != nil {
		return "", err
	}

	// build the path for current daemon
	path += "/lib/" + coinData.LibPath + "-" + releaseVersion + "/bin/" + getOSInfo(coinData).DaemonName
	log.Println("Searching for" + coinData.CurrencyCode + "daemon at " + path)

	// check the current daemon exists
	if !fs.Exists(path) {
		log.Println(coinData.CurrencyCode + "daemon not found for v" + releaseVersion)
		return "", errors.New(coinData.CurrencyCode + "daemon found for v" + releaseVersion)
	} else {
		log.Println(coinData.CurrencyCode + "daemon located for v" + releaseVersion)
	}

	return path, nil

}

// startCoinDaemons pulls in config for coin data, daemonPath,
// builds the command arguments, and executes start command
func startCoinDaemons(coinData conf.CoinData, daemonPath string) *exec.Cmd {

	log.Println("Booting" + coinData.CurrencyCode + "daemon")

	// build up the command flags from daemon config
	rpcUser := fmt.Sprintf("-rpcuser=%s", conf.DaemonConf.RPCUser)
	rpcPassword := fmt.Sprintf("-rpcpassword=%s", conf.DaemonConf.RPCPassword)

	cmdStr := []string{rpcUser, rpcPassword}



	if coinData.UseTestNet {
		cmdStr = append(cmdStr, "-testnet")
	}

	if coinData.IndexTransactions {
		cmdStr = append(cmdStr, "-addressindex=1")
	}

	fs.CreateDataDir(coinData.DataDir)
	p, _ := fs.GetCurrentPath()
	p += coinData.DataDir

	s := fmt.Sprintf("-datadir=%s", p)
	cmdStr = append(cmdStr, s)

	// setup to index transactions (required for API functionality)
	cmd := exec.Command(daemonPath, cmdStr...)

	err := cmd.Start()

	if err != nil {
		log.Fatal("Failed to startCoinDaemons the " + coinData.CurrencyCode + "daemon: " + err.Error())
	}

	return cmd

}

// getOSInfo supplies current OS info and the Daemon name for said OS
func getOSInfo(coinData conf.CoinData) OSInfo {

	osInfo := OSInfo{}

	// TODO: put goosList and goarchList in a config to be loaded in via Viper
	//const goosList = "android darwin dragonfly freebsd linux nacl netbsd openbsd plan9 solaris windows zos "
	//const goarchList = "386 amd64 amd64p32 arm armbe arm64 arm64be ppc64 ppc64le mips mipsle mips64 mips64le mips64p32 mips64p32le ppc s390 s390x sparc sparc64"

	// switch on arch then switch on OS
	switch runtime.GOARCH {

	case "amd64":

		switch runtime.GOOS {

		case "windows":
			osInfo.DaemonName = coinData.WindowsDaemonName
			osInfo.OS = "win64"
			break

		case "darwin":
			osInfo.DaemonName = coinData.DarwinDaemonName
			osInfo.OS = "osx64"
			break
		}

		break
	}

	return osInfo

}

// downloadDaemons pieces together release info, path, name
// and passes that info to DownloadExtract function
func downloadDaemons(coinData conf.CoinData) {

	releaseInfo, _ := getReleaseDataForVersion(coinData)

	dlPath, dlName, _ := getDownloadPathAndName(coinData, releaseInfo)

	isGettingDaemon = true // flag we are getting the daemon

	fs.DownloadExtract(dlPath, dlName)

	isGettingDaemon = false // flag we have finished

}

// getReleaseDataForVersion ranges through the releases and matches
// release TagName to the coin's DaemonVersion via gitHubReleaseInfo function
func getReleaseDataForVersion(coinData conf.CoinData) (GitHubReleaseData, error) {

	log.Println("Attempting to get release data for NAVCoin v" + coinData.DaemonVersion)

	releases, err := gitHubReleaseInfo(coinData.CurrencyCode, coinData.ReleaseAPI)

	var e GitHubReleaseData = GitHubReleaseData{}

	for _, elem := range releases {
		if elem.TagName == coinData.DaemonVersion {
			log.Println("Release data found for NAVCoin v" + coinData.DaemonVersion)
			e = elem.GitHubReleaseData
		}
	}

	return e, err

}

// gitHubReleaseInfo takes the coin's ReleaseAPI and queries for data
func gitHubReleaseInfo(currencyCode string, releaseAPI string) (GitHubReleases, error) {
	log.Println("Retrieving " + currencyCode + "Github release data from: " + releaseAPI)
	response, err := http.Get(releaseAPI)

	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		return GitHubReleases{}, err
	}

	// read the data out to json
	data, _ := ioutil.ReadAll(response.Body)
	c := GitHubReleases{}
	jsonErr := json.Unmarshal(data, &c)

	if jsonErr != nil {
		return GitHubReleases{}, jsonErr
	}

	return c, nil
}

// getDownloadPathAndName ranges through release assets
// and builds/returns downloadPath and downloadName
func getDownloadPathAndName(coinData conf.CoinData, gitHubReleaseData GitHubReleaseData) (string, string, error) {

	log.Println("Getting download path/name for OS from " + coinData.CurrencyCode + " release assest data")

	releaseInfo := gitHubReleaseData

	downloadPath := ""
	downloadName := ""

	for e := range releaseInfo.Assets {

		asset := releaseInfo.Assets[e]

		if strings.Contains(asset.Name, getOSInfo(coinData).OS) {
			// windows os check to provide .zip
			if strings.Contains(asset.Name, "win") {
				if filepath.Ext(asset.Name) == ".zip" {
					log.Println("win64 detected - preparing NAVCoin .zip download")
					downloadPath = releaseInfo.Assets[e].BrowserDownloadURL
					downloadName = releaseInfo.Assets[e].Name
				}
			}
			// osx64 check to provide gzip package :: tar.gz
			if strings.Contains(asset.Name, "osx64") {
				log.Println("osx64 detected - preparing NAVCoin tar.gz download")
				downloadPath = releaseInfo.Assets[e].BrowserDownloadURL
				downloadName = releaseInfo.Assets[e].Name
			} else {
				// TODO: more checks to be added for other systems
				// fall through to defaults :: fire-in-the-hole mode
				downloadPath = releaseInfo.Assets[e].BrowserDownloadURL
				downloadName = releaseInfo.Assets[e].Name
			}
		}
	}

	return downloadPath, downloadName, nil

}
