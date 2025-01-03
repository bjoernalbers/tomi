// tomi
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
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
	serverAddress := flag.Args()[0]
	if os.Geteuid() == 0 {
		log.Fatal("please run as regular user, not as root or with sudo!")
	}
	if err := installTomedo(serverAddress); err != nil {
		log.Fatalf("install tomedo: %v", err)
	}

}

func installTomedo(serverAddress string) error {
	home := os.Getenv("HOME")
	if home == "" {
		return fmt.Errorf("$HOME is not set")
	}
	userAppsDir := filepath.Join(home, "Applications")
	tomedoAppDir := filepath.Join(userAppsDir, "tomedo.app")
	if _, err := os.Stat(tomedoAppDir); err == nil {
		return nil
	}
	// Create localized user applications directory if missing.
	if _, err := os.Stat(userAppsDir); err != nil {
		if err := os.Mkdir(userAppsDir, 0700); err != nil {
			return fmt.Errorf("failed to create user applications directory: %v", err)
		}
		if f, err := os.Create(filepath.Join(userAppsDir, ".localized")); err != nil {
			return fmt.Errorf("failed to localize user applications directory: %v", err)
		} else {
			f.Close()
		}
	}
	// Download and unpack tomedo.
	response, err := http.Get(tomedoDownloadURL(serverAddress))
	if err != nil {
		return fmt.Errorf("failed to download tomedo: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download tomedo: %d %s", response.StatusCode, http.StatusText(response.StatusCode))
	}
	tempFile, err := os.CreateTemp("", "tomedo.app.tar.*")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	if _, err := io.Copy(tempFile, response.Body); err != nil {
		return fmt.Errorf("failed to download tomedo: %v", err)
	}
	if err := Unpack(userAppsDir, tempFile.Name()); err != nil {
		return err
	}
	// Add tomedo to Dock.
	if err := addFileToDock(tomedoAppDir); err != nil {
		return err
	}
	return nil
}

func Usage() {
	header := fmt.Sprintf(`tomi - the tomedo-installer (version %s)

Usage: tomi <tomedo_server_address>`, version)
	fmt.Fprintln(flag.CommandLine.Output(), header)
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

func tomedoDownloadURL(addr string) string {
	return fmt.Sprintf("http://%s:8080/tomedo_live/filebyname/serverinternal/tomedo.app.tar", addr)
}

func addFileToDock(path string) error {
	// https://stackoverflow.com/questions/66106001/add-files-to-dock-on-macos-using-the-command-line
	plist := fmt.Sprintf("<dict><key>tile-data</key><dict><key>file-data</key><dict><key>_CFURLString</key><string>file://%s/</string><key>_CFURLStringType</key><integer>15</integer></dict></dict></dict>", path)
	if err := exec.Command("/usr/bin/defaults", "write", "com.apple.dock", "persistent-apps", "-array-add", plist).Run(); err != nil {
		return fmt.Errorf("failed to add file to Dock: %q", path)
	}
	if output, err := exec.Command("/usr/bin/killall", "Dock").CombinedOutput(); err != nil {
		return fmt.Errorf("failed to restart Dock: %s", string(output))
	}
	return nil
}
