package main

import (
	"flag"
	"log"

	"github.com/Skarlso/acquia-beemo/internal"
	"github.com/Skarlso/acquia-beemo/pkg"
)

func init() {
	flag.BoolVar(&internal.Debug, "v", false, "-v")
	flag.Parse()
}

func main() {
	if internal.Debug {
		internal.LogDebug("[DEBUG] Logger has been turned on!")
	}

	if err := pkg.Serve(); err != nil {
		log.Fatal("Failure starting Beemo: ", err)
	}
}
