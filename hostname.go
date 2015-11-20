package main

/*
#include <unistd.h>
#include <string.h>
*/
import "C"

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
)

func main() {
	f_flag := flag.Bool("f", true, "Include domain information in the printed name. This is the default behavior.")
	s_flag := flag.Bool("s", false, "Trim off any domain information from the printed name.")
	flag.Usage = usage

	flag.Parse()

	*f_flag = true

	if flag.NArg() > 1 {
		usage()
	}

	if flag.NArg() == 1 {
		if err := setHostname(flag.Arg(0)); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		hostname, err := getHostname()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if *s_flag {
			parts := strings.Split(hostname, ".")
			if len(parts) > 1 {
				hostname = parts[0]
			}
		}

		fmt.Println(hostname)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: hostname [-fs] [name-of-host]")
	os.Exit(1)
}

func setHostname(hostname string) error {
	hname := C.CString(hostname)
	_, err := C.sethostname(hname, C.int(C.strlen(hname)))
	return err
}

func getHostname() (hname string, err error) {
	hname, err = syscall.Sysctl("kern.hostname")
	return
}
