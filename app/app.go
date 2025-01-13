package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bjoernalbers/tomi/macos"
)

// App represents an App to install.
type App interface {
	Name() string
	Path() string
	DownloadURL() (string, error)
	Configure() error
}

// Exists returns true if the app exists, otherwise false.
func Exists(a App) bool {
	_, err := os.Stat(a.Path())
	return err == nil
}

// Install performs the actual app installation.
func Install(p App, dock *macos.Dock) error {
	u, err := p.DownloadURL()
	if err != nil {
		return err
	}
	filename, err := Download(u)
	if err != nil {
		return err
	}
	defer os.Remove(filename)
	installDir := filepath.Dir(p.Path())
	if err := Unpack(installDir, filename); err != nil {
		return err
	}
	if err := p.Configure(); err != nil {
		return err
	}
	if err := dock.Add(p.Path()); err != nil {
		return err
	}
	return nil
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
