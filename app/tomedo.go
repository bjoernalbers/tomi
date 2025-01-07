package app

import (
	"net/url"
	"os"
	"path/filepath"
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

func (p *Tomedo) Install() error {
	home, err := GetHome()
	if err != nil {
		return err
	}
	userAppsDir, err := CreateUserAppsDir(home)
	if err != nil {
		return err
	}
	appDir := filepath.Join(userAppsDir, "tomedo.app")
	if _, err := os.Stat(appDir); err == nil {
		return nil
	}
	u, err := p.DownloadURL()
	if err != nil {
		return err
	}
	filename, err := Download(u)
	if err != nil {
		return err
	}
	defer os.Remove(filename)
	if err := Unpack(userAppsDir, filename); err != nil {
		return err
	}
	if err := AddFileToDock(appDir); err != nil {
		return err
	}
	return nil
}
