// tomi
package main

import (
	"flag"
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
	log.Print("sorry, not implemened yet.")
}
