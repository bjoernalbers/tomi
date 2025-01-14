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

	"github.com/bjoernalbers/tomi/macos"
)

type arzekoURL struct {
	ServerURL *url.URL
	Arch      string
}

// String returns the download URL of Arzeko or an error, if the URL could not
// be determined.
func (d *arzekoURL) String() (string, error) {
	resp, err := http.Get(d.autoUpdate())
	if err != nil {
		return "", fmt.Errorf("%T: %v", d, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%T: HTTP %d %s (%s)", d, resp.StatusCode, http.StatusText(resp.StatusCode), resp.Request.URL)
	}
	latestVersion := struct{ URL string }{}
	if err := json.NewDecoder(resp.Body).Decode(&latestVersion); err != nil {
		return "", fmt.Errorf("%T: parse JSON response: %v", d, err)
	}
	return d.replaceHost(latestVersion.URL)
}

// autoUpdate returns the URL to query the download URL of Arzeko.
func (d *arzekoURL) autoUpdate() string {
	u := d.ServerURL.JoinPath("arzeko/latestmac")
	// Arzeko itself queries latest version with "aarch=intel" on
	// Intel-based Macs and "aarch=arm" on Apple Silicon.
	arch := "intel"
	if strings.HasPrefix(d.Arch, "arm") {
		arch = "arm"
	}
	u.RawQuery = fmt.Sprintf("version=0.0.0&aarch=%s", arch)
	return u.String()
}

// replaceHost replaces the host in given download URL by the server host,
// because the tomedo demo server returns unreachable download URLs containing
// localhost:
//
// $ curl 'http://allgemeinmedizin.demo.tomedo.org:8080/tomedo_live/arzeko/latestmac?version=0.0.0'
// {"url":"http://127.0.0.1:9901/tomedo_live/filebyname/serverinternalzip/arzeko/Arzeko-1.146.6-mac.zip","name":"1.146.6","notes":"Aarch: null","pub_date":null}%
func (d *arzekoURL) replaceHost(downloadURL string) (string, error) {
	u, err := url.Parse(downloadURL)
	if err != nil {
		return "", fmt.Errorf("%T: replace host: %v", d, err)
	}
	u.Host = d.ServerURL.Host
	return u.String(), nil
}

type Arzeko struct {
	macos.App
	ServerURL *url.URL
	Arch      string
	Home      string
}

func (p *Arzeko) Path() string {
	return p.App.Path
}

func (p *Arzeko) DownloadURL() (string, error) {
	d := &arzekoURL{ServerURL: p.ServerURL, Arch: p.Arch}
	return d.String()
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
