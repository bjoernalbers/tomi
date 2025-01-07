package app

import (
	"net/url"
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
