package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

import (
	"github.com/sternix/commands/lib/sysctl"
)

var (
	fflag = flag.Bool("f", true, "Include domain information in the printed name. This is the default behavior.")
	sflag = flag.Bool("s", false, "Trim off any domain information from the printed name.")
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	*fflag = true

	if flag.NArg() > 1 {
		usage()
	}

	if flag.NArg() == 1 {
		if err := setHostname(flag.Arg(0)); err != nil {
			log.Fatalln(err)
		}
	} else {
		hostname, err := getHostname()
		if err != nil {
			log.Fatalln(err)
		}

		if *sflag {
			parts := strings.Split(hostname, ".")
			if len(parts) > 1 {
				hostname = parts[0]
			}
		}

		fmt.Println(hostname)
	}
}

func usage() {
	log.Fatalln("usage: hostname [-fs] [name-of-host]")
}

func setHostname(hostname string) error {
	return sysctl.SetString("kern.hostname", hostname)
}

func getHostname() (string, error) {
	return sysctl.ByName("kern.hostname")
}
