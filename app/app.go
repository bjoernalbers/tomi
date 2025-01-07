package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// App represents an App to install.
type App interface {
	Name() string
	DownloadURL() (string, error)
}

// Install performs the actual app installation.
func Install(home string, p App) error {
	userAppsDir, err := CreateUserAppsDir(home)
	if err != nil {
		return err
	}
	appDir := filepath.Join(userAppsDir, p.Name())
	if _, err := os.Stat(appDir); err == nil {
		return nil
	}
	u, err := p.DownloadURL()
	if err != nil {
		return err
	}
	filename, err := Download(u)
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