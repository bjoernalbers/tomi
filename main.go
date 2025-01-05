// tomi
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

// version gets set via ldflags
var version = "unset"

func init() {
	log.SetFlags(0)
	log.SetPrefix("tomi: ")
	flag.Usage = Usage
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Println("server address missing")
		Usage()
		os.Exit(1)
	}
	serverURL := flag.Args()[0]
	if os.Geteuid() == 0 {
		log.Fatal("please run as regular user, not as root or with sudo!")
	}
	u, err := url.Parse(serverURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := installTomedo(u); err != nil {
		log.Fatalf("install tomedo: %v", err)
	}
}

type Tomedo struct {
}

func (p *Tomedo) DownloadURL(serverURL *url.URL) string {
	return serverURL.JoinPath("filebyname/serverinternal/tomedo.app.tar").String()
}

func installTomedo(serverURL *url.URL) error {
	home, err := GetHome()
	if err != nil {
		return err
	}
	userAppsDir, err := CreateUserAppsDir(home)
	if err != nil {
		return err
	}
	appDir := filepath.Join(userAppsDir, "tomedo.app")
	if _, err := os.Stat(appDir); err == nil {
		return nil
	}
	tomedo := Tomedo{}
	filename, err := Download(tomedo.DownloadURL(serverURL))
	if err != nil {
		return err
	}
	defer os.Remove(filename)
	if err := Unpack(userAppsDir, filename); err != nil {
		return err
	}
	if err := AddFileToDock(appDir); err != nil {
		return err
	}
	return nil
}

// GetHome returns the path to the user's home directory.
// If $HOME is not set, an error is returned.
func GetHome() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return home, fmt.Errorf("$HOME is not set")
	}
	return home, nil
}

// CreateUserAppsDir creates a localized user application directory if missing
// and returns its path.
func CreateUserAppsDir(home string) (string, error) {
	dir := filepath.Join(home, "Applications")
	if _, err := os.Stat(dir); err == nil {
		return dir, nil
	}
	if err := os.Mkdir(dir, 0700); err != nil {
		return "", fmt.Errorf("create user apps dir: %v", err)
	}
	if f, err := os.Create(filepath.Join(dir, ".localized")); err != nil {
		return "", fmt.Errorf("localize user apps dir: %v", err)
	} else {
		f.Close()
	}
	return dir, nil
}

func Usage() {
	header := fmt.Sprintf(`tomi - the tomedo-installer (version %s)

Usage: tomi <tomedo_server_url>`, version)
	fmt.Fprintln(flag.CommandLine.Output(), header)
}

// Download downloads URL into temp. file and returns its filename.
// If the download fails, an error is returned.
func Download(url string) (filename string, err error) {
	response, err := http.Get(url)
	if err != nil {
		return filename, fmt.Errorf("download: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return filename, fmt.Errorf("download: %d %s", response.StatusCode, http.StatusText(response.StatusCode))
	}
	tempFile, err := os.CreateTemp("", filepath.Base(url))
	if err != nil {
		return filename, fmt.Errorf("create temp. file: %v", err)
	}
	defer tempFile.Close()
	if _, err := io.Copy(tempFile, response.Body); err != nil {
		return filename, fmt.Errorf("copy download to temp. file: %v", err)
	}
	return tempFile.Name(), nil
}

// Unpack extracts content from archive file into dir.
// The archive is extracted using the macOS tar command which can handle many
// different archive formats including tar and zip (see "man tar").
func Unpack(dir, file string) error {
	cmd := exec.Command("/usr/bin/tar", "-x", "-f", file, "-C", dir)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("unpack: %s", string(output))
	}
	return nil
}

// AddFileToDock adds the file to the macOS Dock.
func AddFileToDock(file string) error {
	// https://stackoverflow.com/questions/66106001/add-files-to-dock-on-macos-using-the-command-line
	plist := fmt.Sprintf("<dict><key>tile-data</key><dict><key>file-data</key><dict><key>_CFURLString</key><string>file://%s/</string><key>_CFURLStringType</key><integer>15</integer></dict></dict></dict>", file)
	if err := exec.Command("/usr/bin/defaults", "write", "com.apple.dock", "persistent-apps", "-array-add", plist).Run(); err != nil {
		return fmt.Errorf("add file to Dock: %q", file)
	}
	if output, err := exec.Command("/usr/bin/killall", "Dock").CombinedOutput(); err != nil {
		return fmt.Errorf("restart Dock: %s", string(output))
	}
	return nil
}
