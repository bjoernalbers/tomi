package app

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

const serverURL = "http://tomedo.example.com:8080/tomedo_live/"

func ServerURL() *url.URL {
	u, err := url.Parse(serverURL)
	if err != nil {
		panic(err)
	}
	return u
}

func TestCreateUserAppsDir(t *testing.T) {
	testHome, err := os.MkdirTemp("", "home")
	if err != nil {
		t.Fatalf("create test home: %v", err)
	}
	defer os.RemoveAll(testHome)
	dir, err := CreateUserAppsDir(testHome)
	if err != nil {
		t.Fatalf("create user apps dir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(testHome, "Applications")); err != nil {
		t.Fatalf("user apps dir missing: %v", err)
	}
	if filepath.Base(dir) != "Applications" {
		t.Fatalf("wrong dirname: %q", dir)
	}
}
