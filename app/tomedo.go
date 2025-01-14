package app

import (
	"fmt"
	"net/url"
	"os/exec"

	"github.com/bjoernalbers/tomi/macos"
)

type tomedoURL struct {
	ServerURL *url.URL
}

// String returns the download URL of tomedo.
func (u *tomedoURL) String() string {
	return u.ServerURL.JoinPath("filebyname/serverinternal/tomedo.app.tar").String()
}

type Tomedo struct {
	macos.App
	ServerURL *url.URL
}

func (p *Tomedo) Path() string {
	return p.App.Path
}

func (p *Tomedo) DownloadURL() (string, error) {
	u := &tomedoURL{p.ServerURL}
	return u.String(), nil
}

func (p *Tomedo) Configure() error {
	cmd := exec.Command("/usr/bin/defaults", "write", "com.zollsoft.tomedo-macOS", "ZS_ServerAdress", p.ServerURL.String())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%T: %v", p, err)
	}
	return nil
}
