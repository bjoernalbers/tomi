package app

import (
	"os"

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
	Install() error
}

// install downloads the URL and unpacks it into dir.
func install(dir, url string) error {
	filename, err := macos.Download(url)
	if err != nil {
		return err
	}
	defer os.Remove(filename)
	if err := macos.Unpack(dir, filename); err != nil {
		return err
	}
	return nil
}
