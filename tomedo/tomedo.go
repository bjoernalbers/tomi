package tomedo

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"

	"github.com/bjoernalbers/tomi/macos"
)

// App represents an App to install.
type App interface {
	Name() string
	Path() string
	Dir() string
	Exists() bool
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

type Tomedo struct {
	macos.App
	ServerURL *url.URL
}

func (p *Tomedo) Install() error {
	return install(p.Dir(), tomedoURL(p.ServerURL))
}

func (p *Tomedo) Path() string {
	return p.App.Path
}

func (p *Tomedo) Configure() error {
	cmd := exec.Command("/usr/bin/defaults", "write", "com.zollsoft.tomedo-macOS", "ZS_ServerAdress", p.ServerURL.String())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%T: %v", p, err)
	}
	return nil
}

// tomedoURL returns the download URL of tomedo.
func tomedoURL(serverURL *url.URL) string {
	return serverURL.JoinPath("filebyname/serverinternal/tomedo.app.tar").String()
}
