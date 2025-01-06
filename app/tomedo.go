package app

import (
	"net/url"
	"os"
	"path/filepath"
)

type Tomedo struct {
	ServerURL *url.URL
}

func (p *Tomedo) DownloadURL() string {
	return p.ServerURL.JoinPath("filebyname/serverinternal/tomedo.app.tar").String()
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
	filename, err := Download(p.DownloadURL())
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
