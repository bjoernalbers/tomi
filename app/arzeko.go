package app

import (
	"fmt"
	"net/url"
	"strings"
)

type Arzeko struct {
	ServerURL *url.URL
	Arch      string
}

// latestVersionURL returns the URL to query the download URL of the latest
// Arzeko version.
func (p *Arzeko) latestVersionURL() string {
	u := p.ServerURL.JoinPath("arzeko/latestmac")
	// Arzeko itself queries latest version with "aarch=intel" on
	// Intel-based Macs and "aarch=arm" on Apple Silicon.
	arch := "intel"
	if strings.HasPrefix(p.Arch, "arm") {
		arch = "arm"
	}
	u.RawQuery = fmt.Sprintf("version=0.0.0&aarch=%s", arch)
	return u.String()
}
