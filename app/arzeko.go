package app

import "net/url"

type Arzeko struct {
	ServerURL *url.URL
}

func (p *Arzeko) latestVersionURL() string {
	u := p.ServerURL.JoinPath("arzeko/latestmac")
	u.RawQuery = "version=0.0.0&aarch=intel"
	return u.String()
}
