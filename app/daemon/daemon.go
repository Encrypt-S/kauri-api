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
	"github.com/Encrypt-S/kauri-api/app/daemon/daemonrpc"
	"github.com/Encrypt-S/kauri-api/app/fs"
)

const (
	WindowsDaemonName string = "navcoind.exe"
	DarwinDaemonName  string = "navcoind"
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
var minHeartbeat int64 = 1000 // the lowest value the hb checker can be set to

var isGettingDaemon = false

// CoinStartManager is a simple system that checks if the coin's daemon
// is alive. If not it tries to startCoinDaemon it with proper config
func CoinStartManager(coin conf.Coin) {

	// set the heartbeat interval but make sure it is not
	// less than the min heartbeat setting
	hbInterval := minHeartbeat
	if coin.DaemonHeartbeat > hbInterval {
		hbInterval = coin.DaemonHeartbeat
	}

	ticker := time.NewTicker(time.Duration(hbInterval) * time.Millisecond)
	go func() {
		for range ticker.C {

			// check to see if the daemon is alive
			if isAlive() {
				log.Println("NAVCoin daemon is alive!")
			} else {

				// only do thing if we are already not getting the daemon
				if !isGettingDaemon {
					log.Println("NAVCoin daemon is unresponsive...")

					if runningDaemon != nil {
						Stop(runningDaemon)
					}

					// startCoinDaemon the daemon and download it if necessary
					cmd, err := DownloadAndStartCoin(coin)

					if err != nil {
						log.Println(err)
					} else {
						runningDaemon = cmd
					}
				}

			}

		}
	}()
}

// isAlive performs a simple rpc command to the Daemon
// returns false on error
func isAlive() bool {

	isLiving := true

	n := daemonrpc.RPCRequestData{}
	n.Method = "getblockcount"

	_, err := daemonrpc.RequestDaemon(n, conf.NavConf)

	if err != nil {
		isLiving = false
	}

	return isLiving

}

// DownloadAndStartCoin checks for supplied coin daemon
// and either downloads it or starts it up if detected
func DownloadAndStartCoin(coin conf.Coin) (*exec.Cmd, error) {

	path, err := CheckForCoinDaemon(coin)

	// download coin daemon if not found
	if err != nil {
		downloadDaemon(coin)
	} else {
		return startCoinDaemon(coin, path), nil
	}

	return startCoinDaemon(coin, path), nil

}

func Stop(cmd *exec.Cmd) {

	if err := cmd.Process.Kill(); err != nil {
		log.Fatal("failed to kill: ", err)
	}
}

// CheckForCoinDaemon checks for supplied coin's current daemon
// in appropriate path and reports back to DownLoadAndStartCoin
func CheckForCoinDaemon(coin conf.Coin) (string, error) {

	// get the latest release info
	releaseVersion := coin.DaemonVersion

	log.Println("Checking NAVCoin daemon for v" + releaseVersion)

	// get the apps current path
	path, err := fs.GetCurrentPath()
	if err != nil {
		return "", err
	}

	// build the path
	path += "/lib/navcoin-" + releaseVersion + "/bin/" + getOSInfo().DaemonName
	log.Println("Searching for NAVCoin daemon at " + path)

	// check the daemon exists
	if !fs.Exists(path) {
		log.Println("NAVCoin daemon not found for v" + releaseVersion)
		return "", errors.New("NAVCoin daemon found for v" + releaseVersion)
	} else {
		log.Println("NAVCoin daemon located for v" + releaseVersion)
	}

	return path, nil

}

// startCoinDaemon pulls in proper config for supplied coin(s),
// builds the command path, and executes start command
func startCoinDaemon(coin conf.Coin, daemonPath string) *exec.Cmd {

	log.Println("Booting NAVCoin daemon")

	// build up the command flags from config
	rpcUser := fmt.Sprintf("-rpcuser=%s", conf.NavConf.RPCUser)
	rpcPassword := fmt.Sprintf("-rpcpassword=%s", conf.NavConf.RPCPassword)
	addressIndex := fmt.Sprintf("-addressindex=%s", coin.CmdAddressIndex)
	network := coin.CmdNetwork

	dataDir :=

	// setup to index transactions (required for API functionality)
	cmd := exec.Command(daemonPath, rpcUser, rpcPassword, addressIndex, dataDir, network)

	err := cmd.Start()

	if err != nil {
		log.Fatal("Failed to startCoinDaemon the daemon: " + err.Error())
	}

	return cmd

}

// getOSInfo supplies current OS info and the Daemon name for said OS
func getOSInfo() OSInfo {

	osInfo := OSInfo{}

	// TODO: put goosList and goarchList in a config to be loaded in via Viper
	//const goosList = "android darwin dragonfly freebsd linux nacl netbsd openbsd plan9 solaris windows zos "
	//const goarchList = "386 amd64 amd64p32 arm armbe arm64 arm64be ppc64 ppc64le mips mipsle mips64 mips64le mips64p32 mips64p32le ppc s390 s390x sparc sparc64"

	switch runtime.GOARCH {

	case "amd64":

		switch runtime.GOOS {

		case "windows":

			osInfo.DaemonName = WindowsDaemonName
			osInfo.OS = "win64"
			break

		case "darwin":

			osInfo.DaemonName = DarwinDaemonName
			osInfo.OS = "osx64"
			break
		}

		break
	}

	return osInfo

}

// downloadDaemon pieces together release info, path, name
// and passes that info to DownloadExtract function
func downloadDaemon(coin conf.Coin) {

	releaseInfo, _ := getReleaseDataForVersion(coin)

	dlPath, dlName, _ := getDownloadPathAndName(releaseInfo)

	isGettingDaemon = true // flag we are getting the daemon

	fs.DownloadExtract(dlPath, dlName)

	isGettingDaemon = false // flag we have finished

}

// getReleaseDataForVersion ranges through the releases and matches
// release TagName to the coin's DaemonVersion via gitHubReleaseInfo function
func getReleaseDataForVersion(coin conf.Coin) (GitHubReleaseData, error) {

	log.Println("Attempting to get release data for NAVCoin v" + coin.DaemonVersion)

	releases, err := gitHubReleaseInfo(coin.ReleaseAPI)

	var e GitHubReleaseData = GitHubReleaseData{}

	for _, elem := range releases {
		if elem.TagName == coin.DaemonVersion {
			log.Println("Release data found for NAVCoin v" + coin.DaemonVersion)
			e = elem.GitHubReleaseData
		}
	}

	return e, err

}

// gitHubReleaseInfo takes the coin's ReleaseAPI and queries for data
func gitHubReleaseInfo(releaseAPI string) (GitHubReleases, error) {
	log.Println("Retrieving NAVCoin Github release data from: " + releaseAPI)
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
		log.Fatal(jsonErr)
	}

	return c, nil
}

// getDownloadPathAndName ranges through release assets
// and builds/returns downloadPath and downloadName
func getDownloadPathAndName(gitHubReleaseData GitHubReleaseData) (string, string, error) {

	log.Println("Getting download path/name for OS from release assest data")

	releaseInfo := gitHubReleaseData

	downloadPath := ""
	downloadName := ""

	for e := range releaseInfo.Assets {

		asset := releaseInfo.Assets[e]

		if strings.Contains(asset.Name, getOSInfo().OS) {
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
