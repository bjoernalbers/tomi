package macos

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// App represents a macOS app by its path, i.e. "/Applications/Firefox.app".
type App struct {
	Path string
}

// Exists returns true if the app exists, otherwise false.
func (a *App) Exists() bool {
	_, err := os.Stat(a.Path)
	return err == nil
}

// Name returns the basename of the App, i.e. "Firefox.app".
func (a *App) Name() string {
	return filepath.Base(a.Path)
}

// Dir returns the parent directory of the App.
func (a *App) Dir() string {
	return filepath.Dir(a.Path)
}

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
