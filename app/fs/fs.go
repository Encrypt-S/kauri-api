package fs

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
)

// WriteCounter counts the number of bytes written to it
type WriteCounter struct {
	Total uint64
}

// wc writes bytes and prints progress
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// PrintProgress will print current status of download
func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

// DownloadExtract sets up and runs the functions
// needed for downloading and extracting of assets
func DownloadExtract(url string, assetName string) error {

	path, err := GetCurrentPath()

	downloadLocation := path + "/" + assetName

	extractPath := path + "/lib"

	Download(url, downloadLocation)

	Extract(assetName, downloadLocation, extractPath)

	if err != nil {
		return err
	}

	return nil

}

// CreateDataDir makes a directory for the data dir
// of active coin's daemon at app's path context
func CreateDataDir(dirPath string) {

	path, err := GetCurrentPath()
	if err != nil {
		log.Println(err)
	}

	// build the path for current daemon
	path += dirPath

	err = os.MkdirAll(path, 0777)

	if err != nil {
		log.Println(err)
	}
}

// Download performs file download of the given url
// this method provides no feedback to the system
func Download(url string, downloadTofileName string) {

	log.Println("Downloading", url)
	log.Println("Destination", downloadTofileName)
	log.Println("This could take a few mins :)")

	output, err := os.Create(downloadTofileName)

	response, err := http.Get(url)
	if err != nil {
		log.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		log.Println("Error while downloading", url, "-", err)
		return
	}

	log.Println(n, "bytes downloaded")

}

// Extract will call Unzip or Untar depending on the detected file extension
func Extract(assetName string, downloadLocation string, extractPath string) {

	log.Println("Extracting " + assetName + " from " + extractPath + " to " + downloadLocation)

	switch filepath.Ext(assetName) {
	case ".zip":
		Unzip(downloadLocation, extractPath)
	case ".gz":
		Untar(downloadLocation, extractPath)
	}

	log.Println("File extracted to " + extractPath)

}

// Unzip takes a src and destination path and unzips accordingly
func Unzip(src, dest string) error {

	log.Println("Unzip " + src + " from " + dest)

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Untar takes a gzStream and destination path;
// opens gzStream, creates new reader, checks headers
// passes on target and header to
func Untar(gzStream string, dst string) error {

	log.Println("Untar the " + gzStream + " to " + dst)

	r, err := os.Open(gzStream)
	gzr, err := gzip.NewReader(r)

	defer gzr.Close()
	if err != nil {
		return err
	}

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(dst, header.Name)

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

			// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer f.Close()

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
		}

	}
}

// GetCurrentPath gets the path of the go app
func GetCurrentPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)
	log.Println("GetCurrentPath = " + exPath)
	return exPath, nil
}

// Exists reports if the named file or directory exists
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
