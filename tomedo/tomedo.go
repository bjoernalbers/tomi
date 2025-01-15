package tomedo

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bjoernalbers/tomi/macos"
)

// App represents an App to install.
type App interface {
	Name() string
	Path() string
	Exists() bool
	Install() error
	Configure() error
}

type Tomedo struct {
	Dir       string
	ServerURL *url.URL
}

func (p *Tomedo) Name() string {
	return "tomedo.app"
}

func (p *Tomedo) Path() string {
	return join(p.Dir, p.Name())
}

func (p *Tomedo) Exists() bool {
	return exists(p)
}

func (p *Tomedo) Install() error {
	return install(p.Dir, tomedoURL(p.ServerURL))
}

func (p *Tomedo) Configure() error {
	cmd := exec.Command("/usr/bin/defaults", "write", "com.zollsoft.tomedo-macOS", "ZS_ServerAdress", p.ServerURL.String())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%T: %v", p, err)
	}
	return nil
}

// exists returns true if the app exists, otherwise false.
func exists(a App) bool {
	_, err := os.Stat(a.Path())
	return err == nil
}

// join joins dir and name.
func join(dir, name string) string {
	return filepath.Join(dir, name)
}

// tomedoURL returns the download URL of tomedo.
func tomedoURL(serverURL *url.URL) string {
	return serverURL.JoinPath("filebyname/serverinternal/tomedo.app.tar").String()
}

// install downloads the URL and unpacks it into dir.
func install(dir, url string) error {
	filename, err := download(url)
	if err != nil {
		return err
	}
	defer os.Remove(filename)
	if err := macos.Unpack(dir, filename); err != nil {
		return err
	}
	return nil
}

// download downloads URL into temp. file and returns its filename.
// If the download fails, an error is returned.
func download(url string) (filename string, err error) {
	client := http.Client{
		Transport: &http.Transport{
			// Disabling compression will speed up the download of
			// tomedo.app.tar since it doesn't have to be
			// compressed twice.
			DisableCompression: true,
		},
	}
	response, err := client.Get(url)
	if err != nil {
		return filename, fmt.Errorf("download: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return filename, fmt.Errorf("download: %d %s", response.StatusCode, http.StatusText(response.StatusCode))
	}
	tempFile, err := os.CreateTemp("", filepath.Base(url))
	if err != nil {
		return filename, fmt.Errorf("create temp. file: %v", err)
	}
	defer tempFile.Close()
	if _, err := io.Copy(tempFile, response.Body); err != nil {
		return filename, fmt.Errorf("copy download to temp. file: %v", err)
	}
	return tempFile.Name(), nil
}
