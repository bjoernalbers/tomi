package app

import (
	"fmt"
	"net/url"
	"os/exec"
)

type Tomedo struct {
	ServerURL *url.URL
}

func (p *Tomedo) Name() string {
	return "tomedo.app"
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
