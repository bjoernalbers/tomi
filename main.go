// tomi
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bjoernalbers/tomi/app"
	"github.com/bjoernalbers/tomi/macos"
	"github.com/bjoernalbers/tomi/server"
)

// version gets set via ldflags
var version = "unset"

func init() {
	log.SetFlags(0)
	log.SetPrefix("tomi: ")
	flag.Usage = Usage
}

func main() {
	dock := &macos.Dock{}
	s := server.Default()
	flag.StringVar(&s.Addr, "a", s.Addr, "address of tomedo server")
	flag.StringVar(&s.Port, "p", s.Port, "port of tomedo server")
	flag.StringVar(&s.Path, "P", s.Path, "path of tomedo server")
	installArzeko := flag.Bool("A", false, "install Arzeko as well")
	flag.Parse()
	if os.Geteuid() == 0 {
		log.Fatal("please run as regular user, not as root or with sudo!")
	}
	home := os.Getenv("HOME")
	if home == "" {
		log.Fatalf("$HOME is not set")
	}
	userAppsDir, err := CreateUserAppsDir(home)
	if err != nil {
		log.Fatal(err)
	}
	apps := []app.App{}
	tomedo := &app.Tomedo{
		ServerURL: s.URL(),
		App:       macos.App{Path: filepath.Join(userAppsDir, "tomedo.app")},
	}
	if !tomedo.Exists() {
		apps = append(apps, tomedo)
	}
	if *installArzeko {
		arzeko := &app.Arzeko{
			ServerURL: s.URL(),
			Arch:      runtime.GOARCH,
			App:       macos.App{Path: filepath.Join(userAppsDir, "Arzeko.app")},
			Home:      home,
		}
		if !arzeko.Exists() {
			apps = append(apps, arzeko)
		}
	}
	if len(apps) == 0 {
		os.Exit(0)
	}
	for _, a := range apps {
		if err := app.Install(a, dock); err != nil {
			log.Fatalf("%s: %v", a.Name(), err)
		}
	}
	if dock.Changed() {
		dock.Restart()
	}
}

func Usage() {
	header := fmt.Sprintf(`tomi - the missing tomedo-installer (version %s)

Usage: tomi [options]

Options:`, version)
	fmt.Fprintln(flag.CommandLine.Output(), header)
	flag.PrintDefaults()
}

// CreateUserAppsDir creates a localized user application directory if missing
// and returns its path.
func CreateUserAppsDir(home string) (string, error) {
	dir := filepath.Join(home, "Applications")
	if _, err := os.Stat(dir); err == nil {
		return dir, nil
	}
	if err := os.Mkdir(dir, 0700); err != nil {
		return "", fmt.Errorf("create user apps dir: %v", err)
	}
	if f, err := os.Create(filepath.Join(dir, ".localized")); err != nil {
		return "", fmt.Errorf("localize user apps dir: %v", err)
	} else {
		f.Close()
	}
	return dir, nil
}
