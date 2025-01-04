package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestServerURL(t *testing.T) {
	server := Server{Host: "tomedo.example.com", Port: 8080, Path: "tomedo_live/"}
	want := "http://tomedo.example.com:8080/tomedo_live/"
	if got := server.URL(); got != want {
		t.Fatalf("%#v.URL():\ngot:\t%q\nwant:\t%q", server, got, want)
	}
}

func TestServerTomedoDownloadURL(t *testing.T) {
	server := Server{Host: "tomedo.example.com", Port: 8080, Path: "tomedo_live/"}
	want := "http://tomedo.example.com:8080/tomedo_live/filebyname/serverinternal/tomedo.app.tar"
	if got := server.TomedoDownloadURL(); got != want {
		t.Fatalf("%#v.TomedoDownloadURL():\ngot:\t%q\nwant:\t%q", server, got, want)
	}
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
