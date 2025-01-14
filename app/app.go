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
}

// Install performs the actual app installation.
func Install(p App) error {
	u, err := p.DownloadURL()
	if err != nil {
		return err
	}
	filename, err := macos.Download(u)
	if err != nil {
		return err
	}
	defer os.Remove(filename)
	if err := macos.Unpack(p.Dir(), filename); err != nil {
		return err
	}
	return nil
}
