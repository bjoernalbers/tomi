// tomi
package main

import (
	"flag"
	"log"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("tomi: ")
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("server address missing")
	}
	log.Print("sorry, not implemened yet.")
}
