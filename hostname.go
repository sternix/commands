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

func setHostname(hostname string) error {
	hname := C.CString(hostname)
	_, err := C.sethostname(hname, C.int(C.strlen(hname)))
	return err
}

func getHostname() (hname string, err error) {
	hname, err = syscall.Sysctl("kern.hostname")
	return
}

func main() {
	f_flag := flag.Bool("f", true, "")
	s_flag := flag.Bool("s", false, "")

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
