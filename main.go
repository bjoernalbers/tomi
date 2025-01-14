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
	userAppsDir, err := macos.CreateUserAppsDir(home)
	if err != nil {
		log.Fatal(err)
	}
	apps := []app.App{}
	apps = append(apps, &app.Tomedo{
		ServerURL: s.URL(),
		App:       macos.App{Path: filepath.Join(userAppsDir, "tomedo.app")},
	})
	if *installArzeko {
		apps = append(apps, &app.Arzeko{
			ServerURL: s.URL(),
			Arch:      runtime.GOARCH,
			App:       macos.App{Path: filepath.Join(userAppsDir, "Arzeko.app")},
			Home:      home,
		})
	}
	for _, a := range apps {
		if a.Exists() {
			continue
		}
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
