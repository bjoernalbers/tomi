// tomi
package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime"

	"github.com/bjoernalbers/tomi/app"
)

// version gets set via ldflags
var version = "unset"

func init() {
	log.SetFlags(0)
	log.SetPrefix("tomi: ")
	flag.Usage = Usage
}

func main() {
	installArzeko := flag.Bool("A", false, "install Arzeko as well")
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Println("tomedo server URL required")
		Usage()
		os.Exit(1)
	}
	serverURL := flag.Args()[0]
	if os.Geteuid() == 0 {
		log.Fatal("please run as regular user, not as root or with sudo!")
	}
	u, err := url.Parse(serverURL)
	if err != nil {
		log.Fatal(err)
	}
	home := os.Getenv("HOME")
	if home == "" {
		log.Fatalf("$HOME is not set")
	}
	userAppsDir, err := app.CreateUserAppsDir(home)
	if err != nil {
		log.Fatal(err)
	}
	apps := []app.App{}
	tomedo := &app.Tomedo{ServerURL: u, InstallDir: userAppsDir}
	if !app.Exists(tomedo) {
		apps = append(apps, tomedo)
	}
	if *installArzeko {
		arzeko := &app.Arzeko{ServerURL: u, Arch: runtime.GOARCH, InstallDir: userAppsDir, Home: home}
		if !app.Exists(arzeko) {
			apps = append(apps, arzeko)
		}
	}
	if len(apps) == 0 {
		os.Exit(0)
	}
	for _, a := range apps {
		if err := app.Install(home, a); err != nil {
			log.Fatalf("%s: %v", a.Name(), err)
		}
	}
	app.RestartDock()
}

func Usage() {
	header := fmt.Sprintf(`tomi - the missing tomedo-installer (version %s)

Usage: tomi [options] <tomedo_server_url>

Options:`, version)
	fmt.Fprintln(flag.CommandLine.Output(), header)
	flag.PrintDefaults()
}
