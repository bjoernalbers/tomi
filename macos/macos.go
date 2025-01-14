package macos

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// Dock represents the Dock on macOS.
type Dock struct {
	changed bool
}

// Add adds the app to the Dock.
func (d *Dock) Add(path string) error {
	// https://stackoverflow.com/questions/66106001/add-files-to-dock-on-macos-using-the-command-line
	plist := fmt.Sprintf("<dict><key>tile-data</key><dict><key>file-data</key><dict><key>_CFURLString</key><string>file://%s/</string><key>_CFURLStringType</key><integer>15</integer></dict></dict></dict>", path)
	cmd := exec.Command("/usr/bin/defaults", "write", "com.apple.dock", "persistent-apps", "-array-add", plist)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Dock: add app: %q", path)
	}
	d.changed = true
	return nil
}

// Restart restarts the Dock in order to apply changes.
func (d *Dock) Restart() error {
	if output, err := exec.Command("/usr/bin/killall", "Dock").CombinedOutput(); err != nil {
		return fmt.Errorf("Dock: restart: %s", string(output))
	}
	d.changed = false
	return nil
}

// Changed returns true if the Dock has been changed, otherwise false.
func (d *Dock) Changed() bool {
	return d.changed
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
