// tomi
package main

import (
	"flag"
	"fmt"
	"log"
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
	log.Print("sorry, not implemened yet.")
}

func tomedoDownloadURL(addr string) string {
	return fmt.Sprintf("http://%s:8080/tomedo_live/filebyname/serverinternal/tomedo.app.tar", addr)
}
