// tomi
package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

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
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Println("server address missing")
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
	tomedo := app.Tomedo{ServerURL: u}
	if err := tomedo.Install(); err != nil {
		log.Fatalf("install tomedo: %v", err)
	}
}

func Usage() {
	header := fmt.Sprintf(`tomi - the missing tomedo-installer (version %s)

Usage: tomi <tomedo_server_url>`, version)
	fmt.Fprintln(flag.CommandLine.Output(), header)
}
