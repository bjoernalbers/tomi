package tomedo

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bjoernalbers/tomi/macos"
)

// App represents an App to install.
type App interface {
	Name() string
	Path() string
	Exists() bool
	Install() error
	Configure() error
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

type Tomedo struct {
	Dir       string
	ServerURL *url.URL
}

func (p *Tomedo) Name() string {
	return "tomedo.app"
}

func (p *Tomedo) Path() string {
	return join(p.Dir, p.Name())
}

func (p *Tomedo) Exists() bool {
	return exists(p)
}

func (p *Tomedo) Install() error {
	return install(p.Dir, tomedoURL(p.ServerURL))
}

func (p *Tomedo) Configure() error {
	cmd := exec.Command("/usr/bin/defaults", "write", "com.zollsoft.tomedo-macOS", "ZS_ServerAdress", p.ServerURL.String())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%T: %v", p, err)
	}
	return nil
}

// exists returns true if the app exists, otherwise false.
func exists(a App) bool {
	_, err := os.Stat(a.Path())
	return err == nil
}

// join joins dir and name.
func join(dir, name string) string {
	return filepath.Join(dir, name)
}

// tomedoURL returns the download URL of tomedo.
func tomedoURL(serverURL *url.URL) string {
	return serverURL.JoinPath("filebyname/serverinternal/tomedo.app.tar").String()
}
