package app

import (
	"fmt"
	"net/url"
	"os/exec"

	"github.com/bjoernalbers/tomi/macos"
)

type Tomedo struct {
	macos.App
	ServerURL *url.URL
}

func (p *Tomedo) Path() string {
	return p.App.Path
}

func (p *Tomedo) DownloadURL() (string, error) {
	return tomedoURL(p.ServerURL), nil
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
