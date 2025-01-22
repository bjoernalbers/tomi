// tomi
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bjoernalbers/tomi/macos"
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
		if err := buildPkg(version); err != nil {
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

func buildPkg(version string) error {
	payloadDir, err := os.MkdirTemp("", "payload-*")
	if err != nil {
		return fmt.Errorf("create payload dir: %v", err)
	}
	defer os.RemoveAll(payloadDir)
	scriptsDir, err := os.MkdirTemp("", "scripts-*")
	if err != nil {
		return fmt.Errorf("create scripts dir: %v", err)
	}
	defer os.RemoveAll(scriptsDir)
	preinstall, err := os.OpenFile(filepath.Join(scriptsDir, "preinstall"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0777)
	if err != nil {
		return fmt.Errorf("create preinstall script: %v", err)
	}
	script := `#!/bin/sh

set -eu

export PATH=/bin:/sbin:/usr/bin:/usr/sbin

tomi="$(dirname "$0")/tomi"

loggedInUser() {
	local user=$(stat -f %Su /dev/console)

	if [ -z "${user}" ]; then exit 1; fi
	if [ "${user}" = 'root' ]; then exit 1; fi
	echo "${user}"
}

if username=$(loggedInUser); then
	echo "Running tomi as user '${username}'..."
	sudo -u "${username}" --set-home "${tomi}"
else
	echo "ERROR: Could not determine active user."
	exit 1
fi
`
	if _, err := io.Copy(preinstall, strings.NewReader(script)); err != nil {
		return fmt.Errorf("copy preinstall script: %v", err)
	}
	preinstall.Close()
	tomi, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get tomi executable: %v", err)
	}
	if output, err := exec.Command("/bin/cp", "-p", tomi, filepath.Join(scriptsDir, "tomi")).CombinedOutput(); err != nil {
		return fmt.Errorf("copy tomi: %v", string(output))
	}
	cmd := exec.Command("/usr/bin/pkgbuild",
		"--identifier", "de.bjoernalbers.tomedo-installer",
		"--version", version,
		"--scripts", scriptsDir,
		"--root", payloadDir,
		"tomedo-installer.pkg")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}
