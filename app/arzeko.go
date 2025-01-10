package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Arzeko struct {
	ServerURL  *url.URL
	Arch       string
	InstallDir string
	Home       string
}

func (p *Arzeko) Name() string {
	return "Arzeko.app"
}

func (p *Arzeko) Path() string {
	return filepath.Join(p.InstallDir, p.Name())
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

func (p *Arzeko) Configure() error {
	dir, err := createArzekoConfigDir(p.Home)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(dir, "usersetting.json"))
	if err != nil {
		return fmt.Errorf("create Arzeko config file: %v", err)
	}
	defer file.Close()
	if _, err := io.Copy(file, strings.NewReader(p.userSettings())); err != nil {
		return fmt.Errorf("create Arzeko config file: %v", err)
	}
	return nil
}

func (p *Arzeko) userSettings() string {
	protocol := p.ServerURL.Scheme
	ipAddress := p.ServerURL.Hostname()
	httpPort := p.ServerURL.Port()
	basePath := p.ServerURL.Path
	basePath = strings.TrimPrefix(basePath, "/")
	basePath = strings.TrimSuffix(basePath, "/")
	content := fmt.Sprintf(`{
  "protocol": "%s",
  "authPassword": "tomedo",
  "xApiKey": "tomedo",
  "ipAddress": "%s",
  "oneWaySSLPort": 9444,
  "authPortServerTools": 9442,
  "authPort": 9443,
  "httpsPort": 8443,
  "httpPort": %s,
  "basePath": "%s",
  "autoUpdatePath": "arzeko/latestmac?version=",
  "calendarOptions": {
    "weekNumbers": true,
    "dayMaxEvents": 5,
    "minTime": "06:00:00",
    "maxTime": "23:00:00",
    "printSorting": "start,title",
    "snapDuration": "00:15:00",
    "slotDuration": "00:30:00",
    "allDaySlot": true,
    "showNamesMonth": false,
    "eventLimitWeek": 5,
    "showNamesInsteadOfKuerzel": true,
    "selectable": true,
    "restorePosition": false,
    "showAbwesenMonth": false,
    "hiddenDays": [],
    "showNutzerSchichtzeiten": true,
    "showPopovers": true,
    "pauseDuration": 0
  },
  "exportOptions": {
    "abwesenExport": false,
    "eigeneSchichtenOnly": false
  },
  "timezone": 1,
  "disableUpdates": false,
  "logLevel": "info"
}`, protocol, ipAddress, httpPort, basePath)
	return content
}

func createArzekoConfigDir(home string) (string, error) {
	dir := filepath.Join(home, "Library", "Application Support", "Arzeko")
	if _, err := os.Stat(dir); err == nil {
		return dir, nil
	}
	if err := os.Mkdir(dir, 0700); err != nil {
		return "", fmt.Errorf("create Arzeko config dir: %v", err)
	}
	return dir, nil
}
