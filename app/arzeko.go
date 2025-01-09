package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Arzeko struct {
	ServerURL *url.URL
	Arch      string
}

func (p *Arzeko) Name() string {
	return "Arzeko.app"
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

func (p *Arzeko) DownloadURL() (string, error) {
	errPrefix := "Arzeko download URL:"
	resp, err := http.Get(p.latestVersionURL())
	if err != nil {
		return "", fmt.Errorf("%s: %v", errPrefix, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: HTTP %d %s (%s)", errPrefix, resp.StatusCode, http.StatusText(resp.StatusCode), resp.Request.URL)
	}
	latestVersion := struct{ URL string }{}
	if err := json.NewDecoder(resp.Body).Decode(&latestVersion); err != nil {
		return "", fmt.Errorf("%s: parse JSON response: %v", errPrefix, err)
	}
	return latestVersion.URL, nil
}

func (p *Arzeko) Configure(home string) error {
	return nil
}
