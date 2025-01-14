package app

import (
	"fmt"
	"net/http"
	"net/url"
	"os/exec"

	"github.com/bjoernalbers/tomi/macos"
)

type tomedoDownloader struct {
	ServerURL *url.URL
}

func (d *tomedoDownloader) Get() (*http.Response, error) {
	resp, err := http.Get(d.URL())
	if err != nil {
		return nil, fmt.Errorf("%T: %v", d, err)
	}
	return resp, nil
}

func (p *tomedoDownloader) URL() string {
	return p.ServerURL.JoinPath("filebyname/serverinternal/tomedo.app.tar").String()
}

type Tomedo struct {
	macos.App
	ServerURL *url.URL
}

func (p *Tomedo) Path() string {
	return p.App.Path
}

func (p *Tomedo) DownloadURL() (string, error) {
	return p.ServerURL.JoinPath("filebyname/serverinternal/tomedo.app.tar").String(), nil
}

func (p *Tomedo) Configure() error {
	cmd := exec.Command("/usr/bin/defaults", "write", "com.zollsoft.tomedo-macOS", "ZS_ServerAdress", p.ServerURL.String())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%T: %v", p, err)
	}
	return nil
}
