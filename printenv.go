package main

import (
	"flag"
	"fmt"
	"os"
)

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
	fmt.Fprintf(os.Stderr, "usage: printenv [name]\n")
	os.Exit(0)
}

// /usr/bin/printenv
