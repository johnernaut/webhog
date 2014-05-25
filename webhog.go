package main

import (
	"flag"
	"fmt"
	"github.com/johnernaut/webhog/webhog"
	"os"
	"runtime"
)

const VERSION = "v0.1.0"

func main() {
	version := flag.Bool("version", false, "current version")
	flag.Parse()
	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	// All the parallelism are belong to us!
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Load configuration file.
	webhog.LoadConfig()

	// Load DB instance.
	webhog.LoadDB()

	// Load route handlers
	webhog.LoadRoutes()
}
