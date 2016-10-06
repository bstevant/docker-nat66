package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

var dev = flag.String("dev", "", "Egress network device")
var prefix = flag.String("prefix", "", "Internal IPv6 prefix")

func main() {

	flag.Parse()
	if flag.NArg() != 0 {
		log.Fatal("Bad number of argument: %d, expected 0", flag.NArg())
		os.Exit(2)
	}
	if *dev == "" && *prefix == "" {
		log.Fatal("Flag -dev and -prefix are mandatory")
		os.Exit(2)
	}

	// Set initial NAT rule
	set_nat66_postrouting(*dev, *prefix)

	// Clear when signal received
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			unset_nat66_postrouting(*dev, *prefix)
			clear_bindings()
		}
	}()

	// Start docker
	init_docker()

}
