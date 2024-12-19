// tomi
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("tomi: ")
}

func main() {
	if os.Geteuid() == 0 {
		log.Fatal("please run as regular user, not as root or with sudo!")
	}
	home := os.Getenv("HOME")
	if home == "" {
		log.Fatal("$HOME is not set")
	}
	userAppsDir := filepath.Join(home, "Applications")
	tomedoAppDir := filepath.Join(userAppsDir, "tomedo.app")
	if _, err := os.Stat(tomedoAppDir); err == nil {
		os.Exit(0)
	}
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("server address missing")
	}
	serverAddress := flag.Args()[0]
	// Create localized user applications directory if missing.
	if _, err := os.Stat(userAppsDir); err != nil {
		if err := os.Mkdir(userAppsDir, 0700); err != nil {
			log.Fatalf("could not create user applications directory: %v", err)
		}
		if f, err := os.Create(filepath.Join(userAppsDir, ".localized")); err != nil {
			log.Fatalf("could not localize user applications directory: %v", err)
		} else {
			f.Close()
		}
	}
	response, err := http.Get(tomedoDownloadURL(serverAddress))
	if err != nil {
		log.Fatalf("failed to download tomedo: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Fatalf("failed to download tomedo: %d %s", response.StatusCode, http.StatusText(response.StatusCode))
	}
	tempFile, err := os.CreateTemp("", "tomedo.app.tar.*")
	if err != nil {
		log.Fatal("failed to create temporary file: %v", err)
	}
	defer tempFile.Close()
	if _, err := io.Copy(tempFile, response.Body); err != nil {
		log.Fatal("failed to download tomedo: %v", err)
	}
	log.Print(tempFile.Name())
}

func tomedoDownloadURL(addr string) string {
	return fmt.Sprintf("http://%s:8080/tomedo_live/filebyname/serverinternal/tomedo.app.tar", addr)
}
