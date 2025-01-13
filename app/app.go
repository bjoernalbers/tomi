package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bjoernalbers/tomi/macos"
)

// App represents an App to install.
type App interface {
	Name() string
	Path() string
	Dir() string
	Exists() bool
	DownloadURL() (string, error)
	Configure() error
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
	if err := macos.Unpack(p.Dir(), filename); err != nil {
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
