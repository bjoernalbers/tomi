package app

import (
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
)

type Tomedo struct {
	ServerURL  *url.URL
	InstallDir string
}

func (p *Tomedo) Name() string {
	return "tomedo.app"
}

func (p *Tomedo) Path() string {
	return filepath.Join(p.InstallDir, p.Name())
}

func (p *Tomedo) DownloadURL() (string, error) {
	return p.ServerURL.JoinPath("filebyname/serverinternal/tomedo.app.tar").String(), nil
}

func (p *Tomedo) Configure(home string) error {
	cmd := exec.Command("/usr/bin/defaults", "write", "com.zollsoft.tomedo-macOS", "ZS_ServerAdress", p.ServerURL.String())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%T: %v", p, err)
	}
	return nil
}
