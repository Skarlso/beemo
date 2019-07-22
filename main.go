package main

import (
	"flag"
	"log"

	"github.com/Skarlso/acquia-beemo/pkg"
)

func init() {
	flag.BoolVar(&pkg.Debug, "v", false, "-v")
	flag.Parse()
}

func main() {
	if pkg.Debug {
		pkg.LogDebug("[DEBUG] Logger has been turned on!")
	}

	if err := pkg.Serve(); err != nil {
		log.Fatal("Failure starting Beemo: ", err)
	}
}
