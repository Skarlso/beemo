package main

import (
	"flag"
	"log"

	"github.com/Skarlso/beemo/internal"
	"github.com/Skarlso/beemo/pkg"
)

func init() {
	flag.BoolVar(&internal.Debug, "v", false, "-v for verbose output")
	flag.BoolVar(&pkg.Opts.AutoTLS, "auto-tls", false, "--auto-tls")
	flag.StringVar(&pkg.Opts.CacheDir, "cache-dir", "", "--cache-dir /home/user/.server/.cache")
	flag.StringVar(&pkg.Opts.ServerKeyPath, "server-key-path", "", "--server-key-file /home/user/.server/server.key")
	flag.StringVar(&pkg.Opts.ServerCrtPath, "server-crt-path", "", "--server-crt-file /home/user/.server/server.crt")
	flag.StringVar(&pkg.Opts.Port, "port", "9998", "--port 443")
	flag.StringVar(&pkg.Opts.Hostname, "hostname", "localhost", "--hostname beemo.org")
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
