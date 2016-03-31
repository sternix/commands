package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		environ := os.Environ()
		for _, env := range environ {
			fmt.Println(env)
		}
	} else {
		if env := os.Getenv(flag.Arg(0)); env != "" {
			fmt.Println(env)
		}
	}
}

func usage() {
	log.Println("usage: printenv [name]")
	os.Exit(0)
}

// /usr/bin/printenv
