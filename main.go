// tomi
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/bjoernalbers/tomi/macos"
	"github.com/bjoernalbers/tomi/pkg"
	"github.com/bjoernalbers/tomi/tomedo"
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
	server := tomedo.DefaultServer()
	flag.StringVar(&server.Addr, "a", server.Addr, "address of tomedo server")
	flag.StringVar(&server.Port, "p", server.Port, "port of tomedo server")
	flag.StringVar(&server.Path, "P", server.Path, "path of tomedo server")
	installArzeko := flag.Bool("A", false, "install Arzeko as well")
	buildPackage := flag.Bool("b", false, "build installation package to install tomedo from official demo server")
	flag.Parse()
	if os.Geteuid() == 0 {
		log.Fatal("please run as regular user, not as root or with sudo!")
	}
	if *buildPackage {
		if err := pkg.Build(version); err != nil {
			log.Fatalf("build package: %v", err)
		}
		os.Exit(0)
	}
	home := os.Getenv("HOME")
	if home == "" {
		log.Fatalf("$HOME is not set")
	}
	userAppsDir, err := macos.CreateUserAppsDir(home)
	if err != nil {
		log.Fatal(err)
	}
	apps := []tomedo.App{}
	apps = append(apps, &tomedo.Tomedo{
		Dir:       userAppsDir,
		ServerURL: server.URL(),
	})
	if *installArzeko {
		apps = append(apps, &tomedo.Arzeko{
			Dir:       userAppsDir,
			ServerURL: server.URL(),
			Arch:      runtime.GOARCH,
			Home:      home,
		})
	}
	for _, a := range apps {
		if a.Exists() {
			continue
		}
		if err := a.Install(); err != nil {
			log.Fatalf("%s: %v", a.Name(), err)
		}
		if err := a.Configure(); err != nil {
			log.Fatalf("%s: %v", a.Name(), err)
		}
		if err := dock.Add(a.Path()); err != nil {
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
