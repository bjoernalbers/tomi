// tomi
package main

import (
	"flag"
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("tomi: ")
}

func main() {
	if os.Geteuid() == 0 {
		log.Fatal("please run as regular user, not as root (with sudo)")
	}
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("server address missing")
	}
	log.Print("sorry, not implemened yet.")
}
