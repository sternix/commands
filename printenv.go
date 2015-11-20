package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: printenv [name]\n")
}

func main() {
	help := flag.Bool("?", false, "")
	flag.Parse()

	if *help {
		usage()
		os.Exit(0)
	}


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

// /usr/bin/printenv
