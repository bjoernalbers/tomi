// tomi
package main

import (
	"flag"
	"log"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("Server address missing")
	}
	log.Print("Sorry, not implemened yet.")
}
