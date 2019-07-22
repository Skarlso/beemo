package main

import (
	"log"

	"github.com/Skarlso/acquia-beemo/pkg"
)

func main() {
	if err := pkg.Serve(); err != nil {
		log.Fatal("Failure starting Beemo: ", err)
	}
}
